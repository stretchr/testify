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
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/ernesto-jimenez/gogen/imports"
)

var (
	pkg       = flag.String("assert-path", "github.com/stretchr/testify/assert", "Path to the assert package")
	includeF  = flag.Bool("include-format-funcs", false, "include format functions such as Errorf and Equalf")
	outputPkg = flag.String("output-package", "", "package for the resulting code")
	tmplFile  = flag.String("template", "", "What file to load the function template from")
	out       = flag.String("out", "", "What file to write the source code to")
)

type Context struct {
	files map[string]*ast.File
	fset  *token.FileSet
	scope *types.Scope
	docs  *doc.Package

	tags map[string]string
}

func main() {
	flag.Parse()

	var ctx Context

	err := ctx.parsePackageSource(*pkg)
	if err != nil {
		log.Fatal(err)
	}

	funcs, err := ctx.analyzeCode()
	if err != nil {
		log.Fatal(err)
	}

	err = ctx.generateCode(funcs)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Context) generateCode(funcs []testFunc) error {
	tags := map[string]struct{}{}
	for _, f := range funcs {
		tags[f.Tags] = struct{}{}
	}

	for tag := range tags {
		err := c.generateCodeTag(funcs, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) generateCodeTag(funcs []testFunc, tags string) error {
	buff := bytes.NewBuffer(nil)

	tmplHead, tmplFunc, err := parseTemplates()
	if err != nil {
		return err
	}

	// Make imports for given functions set (filtered by build tags)
	importer := imports.New(*outputPkg)
	for _, fn := range funcs {
		if fn.Tags != tags {
			continue
		}

		sig := fn.TypeInfo.Type().(*types.Signature)

		importer.AddImportsFrom(sig.Params())
	}

	// Generate header
	if err := tmplHead.Execute(buff, struct {
		Name      string
		BuildTags string
		Imports   map[string]string
	}{
		*outputPkg,
		tags,
		importer.Imports(),
	}); err != nil {
		return err
	}

	// Generate funcs
	for _, fn := range funcs {
		if fn.Tags != tags {
			continue
		}

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
	output, err := outputFile(tags)
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
		f, err := ioutil.ReadFile(*tmplFile)
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

func outputFile(tags string) (*os.File, error) {
	filename := *out
	if filename == "-" || (filename == "" && *tmplFile == "") {
		return os.Stdout, nil
	}
	if filename == "" {
		if tags != "" {
			tags = "_" + tags
		}
		filename = strings.TrimSuffix(strings.TrimSuffix(*tmplFile, ".tmpl"), ".go") + tags + ".go"
	}
	return os.Create(filename)
}

// analyzeCode takes the types scope and the docs and returns the import
// information and information about all the assertion functions.
func (c *Context) analyzeCode() ([]testFunc, error) {
	testingT := c.scope.Lookup("TestingT").Type().Underlying().(*types.Interface)

	var funcs []testFunc
	// Go through all the top level functions
	for _, fdocs := range c.docs.Funcs {
		// Find the function
		obj := c.scope.Lookup(fdocs.Name)

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

		tags := c.buildTags(obj)

		funcs = append(funcs, testFunc{
			CurrentPkg: *outputPkg,
			Tags:       tags,
			DocInfo:    fdocs,
			TypeInfo:   fn,
		})
	}
	return funcs, nil
}

// parsePackageSource returns the types scope and the package documentation from the package
func (c *Context) parsePackageSource(pkg string) error {
	pd, err := build.Import(pkg, ".", 0)
	if err != nil {
		return err
	}

	c.fset = token.NewFileSet()
	c.files = make(map[string]*ast.File)
	c.tags = make(map[string]string)
	fileList := make([]*ast.File, len(pd.GoFiles))
	for i, fname := range pd.GoFiles {
		src, err := ioutil.ReadFile(path.Join(pd.Dir, fname))
		if err != nil {
			return err
		}
		f, err := parser.ParseFile(c.fset, fname, src, parser.ParseComments|parser.AllErrors)
		if err != nil {
			return err
		}

		c.files[fname] = f
		fileList[i] = f

		c.parseBuildTags(fname, f)
	}

	cfg := types.Config{
		Importer: importer.For("source", nil),
	}
	info := types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	tp, err := cfg.Check(pkg, c.fset, fileList, &info)
	if err != nil {
		return err
	}

	c.scope = tp.Scope()

	ap, _ := ast.NewPackage(c.fset, c.files, nil, nil)
	c.docs = doc.New(ap, pkg, 0)

	return nil
}

func (c *Context) buildTags(o types.Object) string {
	tf := c.fset.File(o.Pos())

	return c.tags[tf.Name()]
}

func (c *Context) parseBuildTags(fname string, f *ast.File) {
	const pref = "// +build "

	for _, g := range f.Comments {
		for _, comm := range g.List {
			t := comm.Text
			if !strings.HasPrefix(t, pref) {
				continue
			}

			t = strings.TrimPrefix(t, pref)
			t = strings.TrimSpace(t)
			t = strings.ReplaceAll(t, ",", "-")
			t = strings.ReplaceAll(t, " ", "_")
			t = strings.ReplaceAll(t, "!", "N")

			c.tags[fname] = t
			return
		}
	}

	c.tags[fname] = ""
}

type testFunc struct {
	CurrentPkg string
	Tags       string
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
	search := fmt.Sprintf("%s", f.DocInfo.Name)
	replace := fmt.Sprintf("%sf", f.DocInfo.Name)
	comment := strings.Replace(f.Comment(), search, replace, -1)
	exp := regexp.MustCompile(replace + `\(((\(\)|[^\n])+)\)`)
	return exp.ReplaceAllString(comment, replace+`($1, "error message %s", "formatted")`)
}

func (f *testFunc) CommentWithoutT(receiver string) string {
	search := fmt.Sprintf("assert.%s(t, ", f.DocInfo.Name)
	replace := fmt.Sprintf("%s.%s(", receiver, f.DocInfo.Name)
	return strings.Replace(f.Comment(), search, replace, -1)
}

var headerTemplate = `{{ with .BuildTags }}// +build {{ . }}
{{ end }}
/*
* CODE GENERATED AUTOMATICALLY WITH github.com/stretchr/testify/_codegen
* THIS FILE MUST NOT BE EDITED BY HAND
*/

package {{.Name}}

{{ with .Imports }}
import (
{{range $path, $name := .}}
	{{$name}} "{{$path}}"{{end}}
)
{{ end }}
{{ if ne .Name "assert" }}var _ assert.TestingT // in case no function required assert package{{ end }}
`

var funcTemplate = `{{.Comment}}
func (fwd *AssertionsForwarder) {{.DocInfo.Name}}({{.Params}}) bool {
	return assert.{{.DocInfo.Name}}({{.ForwardedParams}})
}`
