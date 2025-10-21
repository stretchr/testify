// This program reads all assertion functions from the assert package and
// automatically generates the corresponding requires and forwarded assertions

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
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
)

func main() {
	flag.Parse()

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

	// Generate header
	if err := tmplHead.Execute(buff, struct {
		Name    string
		Imports map[string]string
	}{
		*outputPkg,
		importer.Imports(),
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
	tmpl, err := template.New("function").Funcs(template.FuncMap{
		"replace": strings.ReplaceAll,
	}).Parse(funcTemplate)
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

		funcs = append(funcs, testFunc{*outputPkg, fdocs, fn})
		importer.AddImportsFrom(sig.Params())
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
	files := make(map[string]*ast.File)
	fileList := make([]*ast.File, len(pd.GoFiles))
	for i, fname := range pd.GoFiles {
		src, err := os.ReadFile(path.Join(pd.Dir, fname))
		if err != nil {
			return nil, nil, err
		}
		f, err := parser.ParseFile(fset, fname, src, parser.ParseComments|parser.AllErrors)
		if err != nil {
			return nil, nil, err
		}
		files[fname] = f
		fileList[i] = f
	}

	cfg := types.Config{
		Importer: importer.For("source", nil),
	}
	info := types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	tp, err := cfg.Check(pkg, fset, fileList, &info)
	if err != nil {
		return nil, nil, err
	}

	scope := tp.Scope()

	ap, _ := ast.NewPackage(fset, files, nil, nil)
	docs := doc.New(ap, pkg, 0)

	return scope, docs, nil
}

type testFunc struct {
	CurrentPkg string
	DocInfo    *doc.Func
	TypeInfo   *types.Func
}

func (f *testFunc) Qualifier(p *types.Package) string {
	if p == nil || p.Name() == f.CurrentPkg {
		return ""
	}
	return p.Name()
}

func (f *testFunc) Params() string {
	sig := f.TypeInfo.Type().(*types.Signature)
	params := sig.Params()
	p := ""
	comma := ""
	to := params.Len()
	var i int

	if sig.Variadic() {
		to--
	}
	for i = 1; i < to; i++ {
		param := params.At(i)
		p += fmt.Sprintf("%s%s %s", comma, param.Name(), types.TypeString(param.Type(), f.Qualifier))
		comma = ", "
	}
	if sig.Variadic() {
		param := params.At(params.Len() - 1)
		p += fmt.Sprintf("%s%s ...%s", comma, param.Name(), types.TypeString(param.Type().(*types.Slice).Elem(), f.Qualifier))
	}
	return p
}

func (f *testFunc) ForwardedParams() string {
	sig := f.TypeInfo.Type().(*types.Signature)
	params := sig.Params()
	p := ""
	comma := ""
	to := params.Len()
	var i int

	if sig.Variadic() {
		to--
	}
	for i = 1; i < to; i++ {
		param := params.At(i)
		p += fmt.Sprintf("%s%s", comma, param.Name())
		comma = ", "
	}
	if sig.Variadic() {
		param := params.At(params.Len() - 1)
		p += fmt.Sprintf("%s%s...", comma, param.Name())
	}
	return p
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
	name := f.DocInfo.Name
	nameF := name + "f"
	comment := f.Comment()

	// 1. Best effort replacer for mentions, calls, etc.
	//    that can preserve references to the original function.
	bestEffortReplacer := strings.NewReplacer(
		"["+name+"]", "["+name+"]", // ref to origin func, keep as is
		"["+nameF+"]", "["+nameF+"]", // ref to format func code, keep as is
		name+" ", nameF+" ", // mention in text -> replace
		name+"(", nameF+"(", // function call -> replace
		name+",", nameF+",", // mention in enumeration -> replace
		name+".", nameF+".", // closure of sentence -> replace
		name+"\n", nameF+"\n", // end of line -> replace
	)
	comment = bestEffortReplacer.Replace(comment)

	// 2 Find single line assertion calls of any kind, exluding multi-line ones.
	//   example: // assert.Equal(t, expected, actual) <-- the call must be closed on the same line
	assertFormatFuncExp := regexp.MustCompile(`assert\.` + nameF + `\(.*\)`)
	// 2.1 Extract params and existing message if any.
	//     Note: Unless we start parsing the parameters properly, this is a best-effort solution.
	//     If an assertion call ends with a string parameter, we consider that the message.
	//     Please adjust the assertion examples accordingly if needed.
	const minErrorMessageLength = 10
	paramsExp := regexp.MustCompile(`([^()]*)\((.*)\)`)
	strParamExp := regexp.MustCompile(`"[^"]*"$`)
	comment = assertFormatFuncExp.ReplaceAllStringFunc(comment, func(s string) string {
		oBraces := strings.Count(s, "(")
		cBraces := strings.Count(s, ")")
		if oBraces != cBraces {
			// Skip multi-line examples, where assert call is not closed on the same line.
			return s
		}

		m := paramsExp.FindStringSubmatch(s)
		prefix, params, msg := m[1], strings.Split(m[2], ", "), "error message"

		last := strings.TrimSpace(params[len(params)-1])
		// If last param is a string, consider it the message.
		// It is is too short, it is an assertion value, not a message.
		if strParamExp.MatchString(last) && len(last) > minErrorMessageLength+2 {
			msg = strings.Trim(msg, `"`) + ":"
			params = params[:len(params)-1]
		}

		// Rebuild the call with formatted message, reuse existing message if any.
		params = append(params, `"`+msg+` %s", "formatted"`)
		return prefix + "(" + strings.Join(params, ", ") + ")"
	})

	// 3. Replace calls to multi-line assertions end. Examles like:
	//    search:  //	}, time.Second, 10*time.Millisecond, "condition must never be true")
	//    replace: //	}, time.Second, 10*time.Millisecond, "condition must never be true, more: %s", "formatted")
	endFuncWithStringExp := regexp.MustCompile(`(//[\s]*\},.* )"([^"]+)"\)(\n|$)`)
	comment = endFuncWithStringExp.ReplaceAllString(comment, `$1 "$2, more: %s", "formatted")$3`)

	return comment
}

func (f *testFunc) CommentWithoutT(receiver string) string {
	search := fmt.Sprintf("assert.%s(t, ", f.DocInfo.Name)
	replace := fmt.Sprintf("%s.%s(", receiver, f.DocInfo.Name)
	return strings.Replace(f.Comment(), search, replace, -1)
}

// Standard header https://go.dev/s/generatedcode.
var headerTemplate = `// Code generated with github.com/stretchr/testify/_codegen; DO NOT EDIT.

package {{.Name}}

import (
{{range $path, $name := .Imports}}
	{{$name}} "{{$path}}"{{end}}
)
`

var funcTemplate = `{{.Comment}}
func (fwd *AssertionsForwarder) {{.DocInfo.Name}}({{.Params}}) bool {
	return assert.{{.DocInfo.Name}}({{.ForwardedParams}})
}`
