package goapacheconf

type RewriteRuleDirective struct {
	Directive
}

func (d *RewriteRuleDirective) GetRelatedRewiteCondDirectives() []Directive {
	entries := d.container.GetEntries()
	var dIndex int
	var rcDirectives []Directive

	for index, entry := range entries {
		if entry.Directive == d.rawDirective {
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
				IfModules:    d.IfModules,
				container:    d.container,
				rawDirective: directive,
			})
		}

		if (directive != nil && directive.Identifier != RewriteCond) || block != nil {
			break
		}
	}

	return rcDirectives
}
