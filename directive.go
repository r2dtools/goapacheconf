package goapacheconf

import (
	"slices"
	"strings"

	"github.com/r2dtools/goapacheconf/internal/rawparser"
)

const (
	LoadModule              = "LoadModule"
	ServerName              = "ServerName"
	ServerAlias             = "ServerAlias"
	DocumentRoot            = "DocumentRoot"
	SSLEngine               = "SSLEngine"
	SSLCertificateFile      = "SSLCertificateFile"
	SSLCertificateKeyFile   = "SSLCertificateKeyFile"
	SSLCertificateChainFile = "SSLCertificateChainFile"
	UseCanonicalName        = "UseCanonicalName"
	Include                 = "Include"
	IncludeOptional         = "IncludeOptional"
	RewriteEngine           = "RewriteEngine"
	SetSysEnv               = "SetSysEnv"
	RewriteCond             = "RewriteCond"
	RewriteRule             = "RewriteRule"
	ListenPort              = "Listen"
	Alias                   = "Alias"
	Order                   = "Order"
	Allow                   = "Allow"
	Satisfy                 = "Satisfy"
	CustomLog               = "CustomLog"
	ErrorLog                = "ErrorLog"
	ProxyPass               = "ProxyPass"
)

type Directive struct {
	IfModules []string
	entry     *rawparser.Entry
	container entryContainer
}

func (d *Directive) GetName() string {
	return d.entry.Directive.Identifier
}

func (d *Directive) GetValues() []string {
	return d.entry.Directive.GetExpressions()
}

func (d *Directive) GetFirstValue() string {
	values := d.GetValues()

	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (d *Directive) GetValuesAsString() string {
	return strings.Join(d.GetValues(), " ")
}

func (d *Directive) AddValue(expression string) {
	expressions := d.entry.Directive.GetExpressions()
	expressions = append(expressions, expression)

	d.entry.Directive.SetValues(expressions)
}

func (d *Directive) SetValues(expressions []string) {
	d.entry.Directive.SetValues(expressions)
}

func (d *Directive) SetValue(expression string) {
	d.SetValues([]string{expression})
}

func (d *Directive) AppendNewLine() {
	d.entry.EndNewLines = append(d.entry.EndNewLines, "\n")
}

func (d *Directive) PrependNewLine() {
	d.entry.StartNewLines = append(d.entry.StartNewLines, "\n")
}

func (d *Directive) HasPrependedNewLines() bool {
	return len(d.entry.StartNewLines) > 0
}

func (d *Directive) HasAppendedNewLines() bool {
	return len(d.entry.EndNewLines) > 0
}

func (d *Directive) setContainer(container entryContainer) {
	d.container = container
}

func deleteDirectiveByName(c entryContainer, directiveName string) {
	deleteDirectiveInEntityContainer(c, func(rawDirective *rawparser.Directive) bool {
		return rawDirective.Identifier == directiveName
	})
}

func deleteDirective(c entryContainer, directive Directive) {
	deleteDirectiveInEntityContainer(c, func(rawDirective *rawparser.Directive) bool {
		return directive.entry.Directive == rawDirective
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

func NewDirective(name string, values []string) Directive {
	directiveValues := []*rawparser.Value{}

	for _, value := range values {
		directiveValues = append(directiveValues, &rawparser.Value{Expression: value})
	}

	rawDirective := &rawparser.Directive{
		Identifier: name,
		Values:     directiveValues,
	}
	entry := &rawparser.Entry{Directive: rawDirective}
	directive := Directive{
		entry: entry,
	}

	return directive
}
