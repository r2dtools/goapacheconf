package goapacheconf

type DirectiveUnion interface {
	RewriteRuleDirective | AliasDirective | LoadModuleDirective
}

type directiveLocator interface {
	FindDirectives(directiveName DirectiveName) []Directive
}

func findDirectives[T DirectiveUnion](locator directiveLocator, name DirectiveName) []T {
	var directives []T

	for _, directive := range locator.FindDirectives(name) {
		directives = append(directives, T{Directive: directive})
	}

	return directives
}
