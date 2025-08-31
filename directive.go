package goapacheconf

import (
	"github.com/r2dtools/goapacheconf/internal/container"
	"github.com/r2dtools/goapacheconf/internal/rawparser"
)

type Directive struct {
	rawDirective *rawparser.Directive
	container    container.EntryContainer
}

func (d *Directive) GetName() string {
	return d.rawDirective.Identifier
}

func (d *Directive) GetValues() []string {
	return d.rawDirective.GetExpressions()
}

func (d *Directive) GetFirstValue() string {
	values := d.GetValues()

	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (d *Directive) AddValue(expression string) {
	expressions := d.rawDirective.GetExpressions()
	expressions = append(expressions, expression)

	d.rawDirective.SetValues(expressions)
}

func (d *Directive) SetValues(expressions []string) {
	d.rawDirective.SetValues(expressions)
}

func (d *Directive) SetValue(expression string) {
	d.SetValues([]string{expression})
}

func NewDirective(name string, values []string) Directive {
	directiveValues := []*rawparser.Value{}

	for _, value := range values {
		directiveValues = append(directiveValues, &rawparser.Value{Expression: value})
	}

	rawDirective := &rawparser.Directive{
		Identifier: name,
		Values:     directiveValues,
	}

	return Directive{
		rawDirective: rawDirective,
	}
}
