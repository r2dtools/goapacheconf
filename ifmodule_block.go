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
		return strings.TrimPrefix(parts[0], "mod_")
	}

	return strings.ToLower(name)
}

func (b *IfModuleBlock) FindBlocks(blockName BlockName) []Block {
	var blocks []Block
	ifModules := append(b.IfModules, b.GetModuleName())

	for _, entry := range b.rawBlock.GetEntries() {
		blocks = append(blocks, b.config.findBlocksRecursively(blockName, b.FilePath, b.rawBlock, entry, ifModules, true)...)
	}

	return blocks
}

func (b *IfModuleBlock) FindIfModuleBlocks() []IfModuleBlock {
	return findBlocks[IfModuleBlock](b, IfModule)
}

func (b *IfModuleBlock) FindIfModuleBlocksByModuleName(moduleName string) []IfModuleBlock {
	return findIfModuleBlocksByModuleName(b, moduleName)
}

func (b *IfModuleBlock) AddBlock(name string, parameters []string, begining bool) Block {
	ifModules := append(b.IfModules, b.GetModuleName())

	return newBlock(
		b.rawBlock,
		b.config,
		name,
		parameters,
		ifModules,
		begining,
	)
}

func (b *IfModuleBlock) AddDirective(name string, values []string, begining bool, endWithNewLine bool) Directive {
	ifModules := append(b.IfModules, b.GetModuleName())

	return newDirective(
		b.rawBlock,
		name,
		values,
		ifModules,
		begining,
		endWithNewLine,
	)
}
