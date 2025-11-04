package goapacheconf

type AliasDirective struct {
	Directive
}

func (d *AliasDirective) GetFromLocation() string {
	return d.GetFirstValue()
}

func (d *AliasDirective) GetToLocation() string {
	parameters := d.GetValues()

	if len(parameters) > 1 {
		return parameters[1]
	}

	return ""
}

func NewAliasDirective(toLocation, fromLocation string) AliasDirective {
	directive := NewDirective(Alias, []string{toLocation, fromLocation})

	return AliasDirective{Directive: directive}
}
