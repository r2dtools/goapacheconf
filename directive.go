package goapacheconf

import (
	"slices"

	"github.com/r2dtools/goapacheconf/internal/rawparser"
)

type Directive struct {
	rawDirective *rawparser.Directive
	container    entryContainer
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

func (d *Directive) setContainer(container entryContainer) {
	d.container = container
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

func deleteDirectiveByName(c entryContainer, directiveName string) {
	deleteDirectiveInEntityContainer(c, func(rawDirective *rawparser.Directive) bool {
		return rawDirective.Identifier == directiveName
	})
}

func deleteDirective(c entryContainer, directive Directive) {
	deleteDirectiveInEntityContainer(c, func(rawDirective *rawparser.Directive) bool {
		return directive.rawDirective == rawDirective
	})
}

func deleteDirectiveInEntityContainer(c entryContainer, callback func(directive *rawparser.Directive) bool) {
	entries := c.GetEntries()
	dEntries := []*rawparser.Entry{}
	indexesToDelete := []int{}

	for index, entry := range entries {
		if entry.Directive == nil {
			continue
		}

		if callback(entry.Directive) {
			indexesToDelete = append(indexesToDelete, index)
		}
	}

	for index, entry := range entries {
		if !slices.Contains(indexesToDelete, index) {
			dEntries = append(dEntries, entry)
		}
	}

	setEntries(c, dEntries)
}

func addDirective(c entryContainer, directive Directive, toBegining bool, endWithNewLine bool) {
	entries := c.GetEntries()
	directive.setContainer(c)
	entry := &rawparser.Entry{
		Directive: directive.rawDirective,
	}

	if endWithNewLine {
		entry.EndNewLines = []string{"\n"}
	}

	var prevEntry *rawparser.Entry

	if len(entries) > 0 && !toBegining {
		prevEntry = entries[len(entries)-1]
	}

	if prevEntry == nil || len(prevEntry.EndNewLines) == 0 {
		entry.StartNewLines = []string{"\n"}
	}

	if toBegining {
		entries = append([]*rawparser.Entry{entry}, entries...)
	} else {
		entries = append(entries, entry)
	}

	setEntries(c, entries)
}
