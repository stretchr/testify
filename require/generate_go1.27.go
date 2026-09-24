//go:build go1.27

package require

//go:generate sh -c "cd ../_codegen && go build && cd - && ../_codegen/_codegen -output-package=require -template=require.go.tmpl -out=require_go1.27.go -include-format-funcs"
//go:generate sh -c "cd ../_codegen && go build && cd - && ../_codegen/_codegen -output-package=require -template=require_forward.go.tmpl -out=require_forward_go1.27.go -include-format-funcs"
