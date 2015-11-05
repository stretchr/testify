package require

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

//go:generate go run ../_codegen/main.go -output-package=require -template=require.go.tmpl
