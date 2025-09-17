package goapacheconf

import "strings"

type IfModuleBlock struct {
	Block
}

func (b *IfModuleBlock) GetModuleName() string {
	parameters := b.GetParameters()

	if len(parameters) == 0 {
		return ""
	}

	name := parameters[0]
	parts := strings.Split(name, ".")

	if len(parts) == 2 {
		return parts[0]
	}

	return strings.ToLower(name)
}
