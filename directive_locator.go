package goapacheconf

type directiveLocator interface {
	FindDirectives(directiveName DirectiveName) []Directive
}

func findRewriteRuleDirectives(locator directiveLocator) []RewriteRuleDirective {
	var directives []RewriteRuleDirective

	for _, directive := range locator.FindDirectives(RewriteRule) {
		directives = append(directives, RewriteRuleDirective{
			Directive: directive,
		})
	}

	return directives
}
