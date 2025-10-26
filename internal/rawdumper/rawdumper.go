package rawdumper

import (
	"errors"
	"strings"

	"github.com/r2dtools/goapacheconf/internal/rawparser"
)

const (
	tab   = "	"
	space = " "
)

type RawDumper struct {
	nestingLevel int
}

func (d *RawDumper) Dump(config *rawparser.Config) (string, error) {
	if config == nil {
		return "", errors.New("config is empty")
	}

	result := d.DumpEntries(config.Entries, false)

	return result, nil
}

func (d *RawDumper) DumpEntries(entries []*rawparser.Entry, isInContainer bool) string {
	var result string

	for i, entry := range entries {
		if entry != nil {
			result += strings.Join(entry.StartNewLines, "")
			result += d.DumpEntry(entry, isInContainer && i == 0)
			result += strings.Join(entry.EndNewLines, "")
		}
	}

	return result
}

func (d *RawDumper) DumpEntry(entry *rawparser.Entry, isFirstInContainer bool) string {
	result := ""

	if entry.BlockDirective != nil {
		result += d.dumpBlockDirective(entry)
	} else if entry.Directive != nil {
		result += d.dumpDirective(entry)
	} else if entry.Comment != nil {
		result += d.dumpComment(entry, isFirstInContainer)
	}

	return result
}

func (d *RawDumper) dumpBlockDirective(entry *rawparser.Entry) string {
	blockDirective := entry.BlockDirective
	result := d.getCurrentIdent() + "<" + blockDirective.Identifier
	parameters := strings.Join(blockDirective.GetParametersExpressions(), space)

	if parameters != "" {
		result += space + parameters
	}

	result += space + ">"

	if blockDirective.Content != nil {
		d.increaseNestingLevel()
		result += d.DumpEntries(blockDirective.GetEntries(), true)
		d.decreaseNestingLevel()
	}

	result += d.getCurrentIdent() + "</" + blockDirective.Identifier + ">"

	return result
}

func (d *RawDumper) dumpDirective(entry *rawparser.Entry) string {
	expression := strings.Join(entry.Directive.GetExpressions(), space)
	identfier := d.getCurrentIdent() + entry.GetIdentifier()

	if expression == "" {
		return identfier
	}

	return identfier + space + expression
}

func (d *RawDumper) dumpComment(entry *rawparser.Entry, isFirstInContainer bool) string {
	if isFirstInContainer && len(entry.StartNewLines) == 0 {
		return space + entry.Comment.Value
	}

	return d.getCurrentIdent() + entry.Comment.Value
}

func (d *RawDumper) getCurrentIdent() string {
	return strings.Repeat(tab, d.nestingLevel)
}

func (d *RawDumper) increaseNestingLevel() {
	d.nestingLevel++
}

func (d *RawDumper) decreaseNestingLevel() {
	if d.nestingLevel > 0 {
		d.nestingLevel--
	}
}
