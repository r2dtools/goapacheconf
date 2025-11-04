package goapacheconf

type DirectiveUnion interface {
	RewriteRuleDirective | AliasDirective | LoadModuleDirective
}

type directiveLocator interface {
	FindDirectives(directiveName string) []Directive
}

func findDirectives[T DirectiveUnion](locator directiveLocator, name string) []T {
	var directives []T

	for _, directive := range locator.FindDirectives(name) {
		directives = append(directives, T{Directive: directive})
	}

	return directives
}
