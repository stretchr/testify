package require

type Assertions struct {
	t TestingT
}

func New(t TestingT) *Assertions {
	return &Assertions{
		t: t,
	}
}

//go:generate go run ../_codegen/main.go -output-package=require -template=require_forward.go.tmpl
