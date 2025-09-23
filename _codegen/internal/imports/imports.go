/*
The MIT License (MIT)

Copyright (c) 2015 Ernesto Jiménez

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package imports

import (
	"go/types"
	"os"
	"path/filepath"
	"strings"
)

type Importer interface {
	AddImportsFrom(t types.Type)
	Imports() map[string]string
}

// imports contains metadata about all the imports from a given package
type imports struct {
	currentpkg string
	imp        map[string]string
}

// AddImportsFrom adds imports used in the passed type
func (imp *imports) AddImportsFrom(t types.Type) {
	switch el := t.(type) {
	case *types.Basic:
	case *types.Slice:
		imp.AddImportsFrom(el.Elem())
	case *types.Pointer:
		imp.AddImportsFrom(el.Elem())
	case *types.Named:
		pkg := el.Obj().Pkg()
		if pkg == nil {
			return
		}
		if pkg.Name() == imp.currentpkg {
			return
		}
		imp.imp[cleanImportPath(pkg.Path())] = pkg.Name()
	case *types.Tuple:
		for i := 0; i < el.Len(); i++ {
			imp.AddImportsFrom(el.At(i).Type())
		}
	default:
	}
}

func cleanImportPath(ipath string) string {
	return gopathlessImportPath(
		vendorlessImportPath(ipath),
	)
}

func gopathlessImportPath(ipath string) string {
	paths := strings.Split(os.Getenv("GOPATH"), ":")
	for _, p := range paths {
		ipath = strings.TrimPrefix(ipath, filepath.Join(p, "src")+string(filepath.Separator))
	}
	return ipath
}

// vendorlessImportPath returns the devendorized version of the provided import path.
// e.g. "foo/bar/vendor/a/b" => "a/b"
func vendorlessImportPath(ipath string) string {
	// Devendorize for use in import statement.
	if i := strings.LastIndex(ipath, "/vendor/"); i >= 0 {
		return ipath[i+len("/vendor/"):]
	}
	if strings.HasPrefix(ipath, "vendor/") {
		return ipath[len("vendor/"):]
	}
	return ipath
}

// AddImportsFrom adds imports used in the passed type
func (imp *imports) Imports() map[string]string {
	return imp.imp
}

// New initializes a new structure to track packages imported by the currentpkg
func New(currentpkg string) Importer {
	return &imports{
		currentpkg: currentpkg,
		imp:        make(map[string]string),
	}
}
