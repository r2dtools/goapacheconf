package goapacheconf

import "strings"

type LoadModuleDirective struct {
	Directive
}

func (d *LoadModuleDirective) GetModuleName() string {
	values := d.GetValues()

	if len(values) == 0 {
		return ""
	}

	name := strings.TrimSuffix(values[0], "_module")

	return strings.ToLower(name)
}
