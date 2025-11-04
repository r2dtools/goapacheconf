package goapacheconf

import "github.com/r2dtools/goapacheconf/internal/rawparser"

type RewriteRuleDirective struct {
	Directive
}

func (d *RewriteRuleDirective) GetRelatedRewiteCondDirectives() []Directive {
	entries := d.container.GetEntries()
	var dIndex int
	var rcDirectives []Directive

	for index, entry := range entries {
		if entry.Directive == d.entry.Directive {
			dIndex = index

			break
		}
	}

	for i := dIndex - 1; i >= 0; i-- {
		entry := entries[i]
		directive := entry.Directive
		block := entry.BlockDirective

		if directive != nil && directive.Identifier == RewriteCond {
			rcDirectives = append(rcDirectives, Directive{
				IfModules: d.IfModules,
				entry: &rawparser.Entry{
					StartNewLines: entry.StartNewLines,
					Directive:     directive,
					EndNewLines:   entry.EndNewLines,
				},
				container: d.container,
			})
		}

		if (directive != nil && directive.Identifier != RewriteCond) || block != nil {
			break
		}
	}

	return rcDirectives
}
