// This program reads all assertion functions from the assert package and
// automatically generates the corersponding requires and forwarded assertions

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/ernesto-jimenez/gogen/imports"
)

var (
	pkg       = flag.String("assert-path", "github.com/stretchr/testify/assert", "Path to the assert package")
	outputPkg = flag.String("output-package", "", "package for the resulting code")
	tmplFile  = flag.String("template", "", "What file to load the function template from")
	out       = flag.String("out", "", "What file to write the source code to")
)

func main() {
	flag.Parse()

	tmplHead, err := template.New("header").Parse(defaultTemplate)
	if err != nil {
		log.Fatal(err)
	}
	if *tmplFile != "" {
		f, err := ioutil.ReadFile(*tmplFile)
		if err != nil {
			log.Fatal(err)
		}
		funcTemplate = string(f)
	}
	tmpl, err := template.New("function").Parse(funcTemplate)
	if err != nil {
		log.Fatal(err)
	}

	pd, err := build.Import(*pkg, ".", 0)
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	files := make(map[string]*ast.File)
	fileList := make([]*ast.File, len(pd.GoFiles))
	for i, fname := range pd.GoFiles {
		src, err := ioutil.ReadFile(path.Join(pd.SrcRoot, pd.ImportPath, fname))
		if err != nil {
			log.Fatal(err)
		}
		f, err := parser.ParseFile(fset, fname, src, parser.ParseComments|parser.AllErrors)
		if err != nil {
			log.Fatal(err)
		}
		files[fname] = f
		fileList[i] = f
	}

	cfg := types.Config{
		Importer: importer.Default(),
	}
	info := types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	tp, err := cfg.Check(*pkg, fset, fileList, &info)
	if err != nil {
		log.Fatal(err)
	}

	scope := tp.Scope()
	testingT := scope.Lookup("TestingT").Type().Underlying().(*types.Interface)

	ap, _ := ast.NewPackage(fset, files, nil, nil)
	docs := doc.New(ap, *pkg, 0)

	imports := imports.New(*outputPkg)
	funcs := make([]Func, 0)
	// Go through all the top level functions
	for _, fdocs := range docs.Funcs {
		// Find the function
		obj := scope.Lookup(fdocs.Name)

		fn, ok := obj.(*types.Func)
		if !ok {
			continue
		}
		// Check function signatuer has at least two arguments
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

		funcs = append(funcs, Func{*outputPkg, fdocs, fn})
		imports.AddImportsFrom(sig.Params())
	}

	var output *os.File
	if *out == "-" || (*out == "" && *tmplFile == "") {
		*out = "-"
		output = os.Stdout
	} else if *out == "" {
		*out = strings.TrimSuffix(strings.TrimSuffix(*tmplFile, ".tmpl"), ".go") + ".go"
	}
	if *out != "-" {
		output, err = os.Create(*out)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := tmplHead.Execute(output, struct {
		Name    string
		Imports map[string]string
	}{
		*outputPkg,
		imports.Imports(),
	}); err != nil {
		log.Fatal(err)
	}
	for _, fn := range funcs {
		output.Write([]byte("\n\n"))
		if err := tmpl.Execute(output, &fn); err != nil {
			log.Fatal(err)
		}
	}
}

var defaultTemplate = `/*
* CODE GENERATED AUTOMATICALLY WITH github.com/stretchr/testify/_codegen
* THIS FILE MUST NOT BE EDITED BY HAND
*/

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

type Func struct {
	CurrentPkg string
	DocInfo    *doc.Func
	TypeInfo   *types.Func
}

func (f *Func) Qualifier(p *types.Package) string {
	if p == nil || p.Name() == f.CurrentPkg {
		return ""
	}
	return p.Name()
}

func (f *Func) Params() string {
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

func (f *Func) ForwardedParams() string {
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

func (f *Func) Comment() string {
	return "// " + strings.Replace(strings.TrimSpace(f.DocInfo.Doc), "\n", "\n// ", -1)
}
