// This program reads all assertion functions from the assert package and
// automatically generates the corresponding requires and forwarded assertions

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/build/constraint"
	"go/doc"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/stretchr/testify/_codegen/internal/imports"
)

var (
	pkg       = flag.String("assert-path", "github.com/stretchr/testify/assert", "Path to the assert package")
	includeF  = flag.Bool("include-format-funcs", false, "include format functions such as Errorf and Equalf")
	outputPkg = flag.String("output-package", "", "package for the resulting code")
	tmplFile  = flag.String("template", "", "What file to load the function template from")
	out       = flag.String("out", "", "What file to write the source code to")
	goVersion string
)

func main() {
	flag.Parse()
	var err error
	goVersion, err = inferGoVersion()
	if err != nil {
		log.Fatal(err)
	}

	scope, docs, err := parsePackageSource(*pkg)
	if err != nil {
		log.Fatal(err)
	}

	importer, funcs, err := analyzeCode(scope, docs)
	if err != nil {
		log.Fatal(err)
	}

	if err := generateCode(importer, funcs); err != nil {
		log.Fatal(err)
	}
}

func generateCode(importer imports.Importer, funcs []testFunc) error {
	buff := bytes.NewBuffer(nil)

	tmplHead, tmplFunc, err := parseTemplates()
	if err != nil {
		return err
	}
	if strings.Contains(funcTemplate, "assert.") {
		importer.AddImport(*pkg, "assert")
	}

	// Generate header
	if err := tmplHead.Execute(buff, struct {
		Name      string
		Imports   map[string]string
		GoVersion string
	}{
		*outputPkg,
		importer.Imports(),
		goVersion,
	}); err != nil {
		return err
	}

	// Generate funcs
	for _, fn := range funcs {
		buff.Write([]byte("\n\n"))
		if err := tmplFunc.Execute(buff, &fn); err != nil {
			return err
		}
	}

	code, err := format.Source(buff.Bytes())
	if err != nil {
		return err
	}

	// Write file
	output, err := outputFile()
	if err != nil {
		return err
	}
	defer output.Close()
	_, err = io.Copy(output, bytes.NewReader(code))
	return err
}

func inferGoVersion() (string, error) {
	filename := os.Getenv("GOFILE")
	if filename == "" {
		return "", nil
	}

	source, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("read GOFILE %q: %w", filename, err)
	}
	for _, line := range strings.Split(string(source), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package ") {
			break
		}
		if !strings.HasPrefix(line, "//go:build ") {
			continue
		}

		expr, err := constraint.Parse(line)
		if err != nil {
			return "", fmt.Errorf("parse build constraint in %q: %w", filename, err)
		}
		tag, ok := expr.(*constraint.TagExpr)
		if !ok || !strings.HasPrefix(tag.Tag, "go1.") {
			return "", nil
		}
		return strings.TrimPrefix(tag.Tag, "go"), nil
	}
	return "", nil
}

func parseTemplates() (*template.Template, *template.Template, error) {
	tmplHead, err := template.New("header").Parse(headerTemplate)
	if err != nil {
		return nil, nil, err
	}
	if *tmplFile != "" {
		f, err := os.ReadFile(*tmplFile)
		if err != nil {
			return nil, nil, err
		}
		funcTemplate = string(f)
	}
	tmpl, err := template.New("function").Parse(funcTemplate)
	if err != nil {
		return nil, nil, err
	}
	return tmplHead, tmpl, nil
}

func outputFile() (*os.File, error) {
	filename := *out
	if filename == "-" || (filename == "" && *tmplFile == "") {
		return os.Stdout, nil
	}
	if filename == "" {
		filename = strings.TrimSuffix(strings.TrimSuffix(*tmplFile, ".tmpl"), ".go") + ".go"
	}
	return os.Create(filename)
}

// analyzeCode takes the types scope and the docs and returns the import
// information and information about all the assertion functions.
func analyzeCode(scope *types.Scope, docs *doc.Package) (imports.Importer, []testFunc, error) {
	testingT := scope.Lookup("TestingT").Type().Underlying().(*types.Interface)

	importer := imports.New(*outputPkg)
	var funcs []testFunc
	// Go through all the top level functions
	for _, fdocs := range docs.Funcs {
		// Find the function
		obj := scope.Lookup(fdocs.Name)

		fn, ok := obj.(*types.Func)
		if !ok {
			continue
		}
		// Check function signature has at least two arguments
		sig := fn.Type().(*types.Signature)
		if sig.Params().Len() < 2 {
			continue
		}
		// Check first argument is of type testingT
		first, ok := sig.Params().At(0).Type().(*types.Named)
		if !ok {
			continue
		}
		firstType, ok := first.Underlying().(*types.Interface)
		if !ok {
			continue
		}
		if !types.Implements(firstType, testingT) {
			continue
		}

		// Skip functions ending with f
		if strings.HasSuffix(fdocs.Name, "f") && !*includeF {
			continue
		}
		if (sig.TypeParams().Len() > 0) != (goVersion != "") {
			continue
		}
		results := sig.Results()
		if results.Len() == 0 || !types.Identical(results.At(results.Len()-1).Type(), types.Typ[types.Bool]) {
			return nil, nil, fmt.Errorf("assertion function %s must return bool as its final result", fdocs.Name)
		}

		funcs = append(funcs, testFunc{*outputPkg, fdocs, fn})
		for i := 1; i < sig.Params().Len(); i++ {
			importer.AddImportsFrom(sig.Params().At(i).Type())
		}
		importer.AddImportsFrom(sig.Results())
		for i := 0; i < sig.TypeParams().Len(); i++ {
			importer.AddImportsFrom(sig.TypeParams().At(i).Constraint())
		}
	}
	return importer, funcs, nil
}

// parsePackageSource returns the types scope and the package documentation from the package
func parsePackageSource(pkg string) (*types.Scope, *doc.Package, error) {
	pd, err := build.Import(pkg, ".", 0)
	if err != nil {
		return nil, nil, err
	}

	fset := token.NewFileSet()
	fileList := make([]*ast.File, len(pd.GoFiles))
	for i, fname := range pd.GoFiles {
		src, err := os.ReadFile(path.Join(pd.Dir, fname))
		if err != nil {
			return nil, nil, err
		}
		// SkipObjectResolution for less memory usage in the go/types era
		f, err := parser.ParseFile(fset, fname, src, parser.ParseComments|parser.AllErrors|parser.SkipObjectResolution)
		if err != nil {
			return nil, nil, err
		}
		fileList[i] = f
	}

	cfg := types.Config{
		Importer: importer.ForCompiler(fset, "source", nil),
	}
	info := types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	tp, err := cfg.Check(pkg, fset, fileList, &info)
	if err != nil {
		return nil, nil, err
	}

	scope := tp.Scope()

	docs, err := doc.NewFromFiles(fset, fileList, pkg)
	if err != nil {
		return nil, nil, err
	}

	return scope, docs, nil
}

type testFunc struct {
	CurrentPkg string
	DocInfo    *doc.Func
	TypeInfo   *types.Func
}

func (f *testFunc) signature() *types.Signature {
	return f.TypeInfo.Type().(*types.Signature)
}

func (f *testFunc) Qualifier(p *types.Package) string {
	if p == nil || p.Name() == f.CurrentPkg {
		return ""
	}
	return p.Name()
}

func (f *testFunc) Params() string {
	sig := f.signature()
	params := sig.Params()
	var p strings.Builder
	comma := ""
	to := params.Len()
	var i int

	if sig.Variadic() {
		to--
	}
	for i = 1; i < to; i++ {
		param := params.At(i)
		p.WriteString(comma)
		p.WriteString(param.Name())
		p.WriteString(" ")
		p.WriteString(types.TypeString(param.Type(), f.Qualifier))
		comma = ", "
	}
	if sig.Variadic() {
		param := params.At(params.Len() - 1)
		fmt.Fprintf(&p, "%s%s ...%s", comma, param.Name(), types.TypeString(param.Type().(*types.Slice).Elem(), f.Qualifier))
	}
	return p.String()
}

func (f *testFunc) TypeParams() string {
	typeParams := f.signature().TypeParams()
	if typeParams.Len() == 0 {
		return ""
	}

	var p strings.Builder
	p.WriteByte('[')
	for i := 0; i < typeParams.Len(); i++ {
		if i > 0 {
			p.WriteString(", ")
		}
		typeParam := typeParams.At(i)
		p.WriteString(typeParam.Obj().Name())
		p.WriteByte(' ')
		p.WriteString(types.TypeString(typeParam.Constraint(), f.Qualifier))
	}
	p.WriteByte(']')
	return p.String()
}

func (f *testFunc) TypeArgs() string {
	typeParams := f.signature().TypeParams()
	if typeParams.Len() == 0 {
		return ""
	}

	var p strings.Builder
	p.WriteByte('[')
	for i := 0; i < typeParams.Len(); i++ {
		if i > 0 {
			p.WriteString(", ")
		}
		p.WriteString(typeParams.At(i).Obj().Name())
	}
	p.WriteByte(']')
	return p.String()
}

func (f *testFunc) Results() string {
	return f.formatResults(f.signature().Results())
}

func (f *testFunc) RequireResults() string {
	results := f.signature().Results()
	return f.formatResults(resultsWithoutSuccess(results))
}

func (f *testFunc) HasRequireResults() bool {
	return f.signature().Results().Len() > 1
}

func (f *testFunc) RequireResultNames() string {
	var names strings.Builder
	for i := 0; i < f.signature().Results().Len()-1; i++ {
		if i > 0 {
			names.WriteString(", ")
		}
		fmt.Fprintf(&names, "result%d", i)
	}
	return names.String()
}

func (f *testFunc) formatResults(results *types.Tuple) string {
	switch results.Len() {
	case 0:
		return ""
	case 1:
		return " " + types.TypeString(results.At(0).Type(), f.Qualifier)
	default:
		return " " + types.TypeString(results, f.Qualifier)
	}
}

func resultsWithoutSuccess(results *types.Tuple) *types.Tuple {
	vars := make([]*types.Var, results.Len()-1)
	for i := range vars {
		vars[i] = results.At(i)
	}
	return types.NewTuple(vars...)
}

func (f *testFunc) ForwardedParams() string {
	sig := f.TypeInfo.Type().(*types.Signature)
	params := sig.Params()
	var p strings.Builder
	comma := ""
	to := params.Len()
	var i int

	if sig.Variadic() {
		to--
	}
	for i = 1; i < to; i++ {
		param := params.At(i)
		p.WriteString(comma)
		p.WriteString(param.Name())
		comma = ", "
	}
	if sig.Variadic() {
		param := params.At(params.Len() - 1)
		fmt.Fprintf(&p, "%s%s...", comma, param.Name())
	}
	return p.String()
}

func (f *testFunc) ParamsFormat() string {
	return strings.Replace(f.Params(), "msgAndArgs", "msg string, args", 1)
}

func (f *testFunc) ForwardedParamsFormat() string {
	return strings.Replace(f.ForwardedParams(), "msgAndArgs", "append([]interface{}{msg}, args...)", 1)
}

func (f *testFunc) Comment() string {
	return "// " + strings.Replace(strings.TrimSpace(f.DocInfo.Doc), "\n", "\n// ", -1)
}

func (f *testFunc) CommentFormat() string {
	search := f.DocInfo.Name
	replace := fmt.Sprintf("%sf", f.DocInfo.Name)
	comment := strings.ReplaceAll(f.Comment(), search, replace)

	// NOTE: Some functions have msgAndArgs at the end. It needs to be omitted. (currently only EventuallyWithT)
	// Change here if the original comment changed.
	comment = strings.Replace(comment, `, "external state has not changed to 'true'; still false"`, "", 1)

	exp := regexp.MustCompile(replace + `(\[[^\n]*\])?\((([^()]*|\([^()]*\))*)\)`)
	return exp.ReplaceAllString(comment, replace+`$1($2, "error message %s", "formatted")`)
}

func (f *testFunc) CommentWithoutT(receiver string) string {
	search := regexp.MustCompile(fmt.Sprintf(`assert\.%s(\[[^\n]*\])?\(t, `, regexp.QuoteMeta(f.DocInfo.Name)))
	replace := fmt.Sprintf("%s.%s$1(", receiver, f.DocInfo.Name)
	return search.ReplaceAllString(f.Comment(), replace)
}

func requireComment(comment string) string {
	const returnClause = "Returns whether the assertion was successful (true) or not (false)"
	const emptyString = ""

	comment = strings.ReplaceAll(comment, returnClause+".", emptyString)
	comment = strings.ReplaceAll(comment, returnClause, emptyString)

	// convert from: "if cond {\n  body\n}" to "cond\nbody"
	ifBlockReg := regexp.MustCompile(`(?m)^([\t ]*)if (.+) \{\n((?:.*\n)*?)\s*\}`)

	comment = ifBlockReg.ReplaceAllStringFunc(comment, func(match string) string {
		nestedMatch := ifBlockReg.FindStringSubmatch(match)
		indent, cond, body := nestedMatch[1], nestedMatch[2], nestedMatch[3]
		out := []string{indent + cond}

		for _, line := range strings.Split(
			strings.TrimRight(body, "\n"),
			"\n",
		) {
			if t := strings.TrimSpace(line); t != emptyString {
				out = append(out, indent+t)
			}
		}

		return strings.Join(out, "\n")
	})

	final := strings.TrimSpace(comment)
	return "// " + strings.ReplaceAll(final, "\n", "\n// ")
}

func (f *testFunc) CommentRequire() string {
	comment := strings.ReplaceAll(f.DocInfo.Doc, "assert.", "require.")
	// Preserve assert.CollectT, even in package 'require'
	comment = strings.ReplaceAll(comment, "require.CollectT", "assert.CollectT")
	return requireComment(comment)
}

func (f *testFunc) CommentRequireWithoutT(receiver string) string {
	assertCallRe := regexp.MustCompile(`assert\.(\w+)(\[[^\n]*\])?\(t, `)
	comment := assertCallRe.ReplaceAllString(f.DocInfo.Doc, receiver+".$1$2(")
	return requireComment(comment)
}

// Standard header https://go.dev/s/generatedcode.
var headerTemplate = `{{if .GoVersion}}//go:build go{{.GoVersion}}

{{end}}// Code generated with github.com/stretchr/testify/_codegen; DO NOT EDIT.

package {{.Name}}

{{if .Imports}}
import (
{{range $path, $name := .Imports}}
	{{$name}} "{{$path}}"{{end}}
)
{{end}}
`

var funcTemplate = `{{.Comment}}
func (fwd *AssertionsForwarder) {{.DocInfo.Name}}{{.TypeParams}}({{.Params}}){{.Results}} {
	return assert.{{.DocInfo.Name}}{{.TypeArgs}}({{.ForwardedParams}})
}`
