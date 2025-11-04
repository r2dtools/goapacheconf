package goapacheconf

import (
	"slices"

	"github.com/r2dtools/goapacheconf/internal/rawdumper"
	"github.com/r2dtools/goapacheconf/internal/rawparser"
)

type BlockName string

const (
	VirtualHost   BlockName = "VirtualHost"
	Directory     BlockName = "Directory"
	IfModule      BlockName = "IfModule"
	Proxy         BlockName = "Proxy"
	Location      BlockName = "Location"
	LocationMatch BlockName = "LocationMatch"
)

type Block struct {
	FilePath  string
	IfModules []string
	config    *Config
	container entryContainer
	rawBlock  *rawparser.BlockDirective
	rawDumper *rawdumper.RawDumper
}

func (b *Block) GetName() string {
	return b.rawBlock.Identifier
}

func (b *Block) GetParameters() []string {
	return b.rawBlock.GetParametersExpressions()
}

func (b *Block) GetFirstParameter() string {
	parameters := b.GetParameters()

	if len(parameters) > 0 {
		return parameters[0]
	}

	return ""
}

func (b *Block) SetParameters(parameters []string) {
	b.rawBlock.SetParameters(parameters)
}

func (b *Block) FindDirectives(directiveName DirectiveName) []Directive {
	var directives []Directive

	for _, entry := range b.rawBlock.GetEntries() {
		directives = append(directives, b.config.findDirectivesRecursively(directiveName, b.FilePath, b.rawBlock, entry, b.IfModules, true)...)
	}

	return directives
}

func (b *Block) FindRewriteRuleDirectives() []RewriteRuleDirective {
	return findDirectives[RewriteRuleDirective](b, RewriteRule)
}

func (b *Block) FindBlocks(blockName BlockName) []Block {
	var blocks []Block

	for _, entry := range b.rawBlock.GetEntries() {
		blocks = append(blocks, b.config.findBlocksRecursively(blockName, b.FilePath, b.rawBlock, entry, b.IfModules, true)...)
	}

	return blocks
}

func (b *Block) FindIfModuleBlocks() []IfModuleBlock {
	return findBlocks[IfModuleBlock](b, IfModule)
}

func (b *Block) FindIfModuleBlocksByModuleName(moduleName string) []IfModuleBlock {
	return findIfModuleBlocksByModuleName(b, moduleName)
}

func (b *Block) AppendDirective(directive Directive) Directive {
	directive.IfModules = b.IfModules
	entries := b.rawBlock.GetEntries()
	directive.setContainer(b.rawBlock)

	var prevEntry *rawparser.Entry

	if len(entries) > 0 {
		prevEntry = entries[len(entries)-1]
	}

	if prevEntry == nil || len(prevEntry.EndNewLines) == 0 {
		directive.PrependNewLine()
	}

	entries = append(entries, directive.entry)

	setEntries(b.rawBlock, entries)

	return directive
}

func (b *Block) PrependDirective(directive Directive) Directive {
	directive.IfModules = b.IfModules
	entries := b.rawBlock.GetEntries()
	directive.setContainer(b.rawBlock)

	directive.PrependNewLine()
	entries = append([]*rawparser.Entry{directive.entry}, entries...)

	setEntries(b.rawBlock, entries)

	return directive
}

func (b *Block) DeleteDirective(directive Directive) {
	deleteDirective(b.rawBlock, directive)
}

func (b *Block) DeleteDirectiveByName(directiveName string) {
	deleteDirectiveByName(b.rawBlock, directiveName)
}

func (b *Block) Dump() string {
	entry := rawparser.Entry{
		BlockDirective: b.rawBlock,
	}

	return b.rawDumper.DumpEntry(&entry, false)
}

func (b *Block) AddBlock(name string, parameters []string, begining bool) Block {
	return newBlock(
		b.rawBlock,
		b.config,
		name,
		parameters,
		b.IfModules,
		begining,
	)
}

func (b *Block) GetDirectiveOrder(directive Directive) int {
	order := -1

	for index, entry := range b.rawBlock.GetEntries() {
		if entry.Directive == directive.entry.Directive {
			return index
		}
	}

	return order
}

func (b *Block) GetBlockOrder(block Block) int {
	order := -1

	for index, entry := range b.rawBlock.GetEntries() {
		if entry.BlockDirective == block.rawBlock {
			return index
		}
	}

	return order
}

func (b *Block) ChangeDirectiveOrder(directive Directive, order int) {
	currentOrder := b.GetDirectiveOrder(directive)

	if currentOrder == -1 || currentOrder == order {
		return
	}

	entries := b.rawBlock.GetEntries()

	if order > len(entries)-1 {
		order = len(entries) - 1
	}

	if order < 0 {
		order = 0
	}

	entries = slices.Insert(entries, order, directive.entry)

	if order > currentOrder {
		entries = slices.Delete(entries, currentOrder, min(currentOrder+1, len(entries)-1))
	} else {
		entries = slices.Delete(entries, currentOrder+1, min(currentOrder+2, len(entries)-1))
	}

	b.rawBlock.SetEntries(entries)
}

func newBlock(c entryContainer, config *Config, name string, parameters []string, ifModules []string, begining bool) Block {
	rawBlock := &rawparser.BlockDirective{
		Identifier: name,
		Content:    &rawparser.BlockContent{},
	}
	rawBlock.SetParameters(parameters)

	block := Block{
		IfModules: ifModules,
		config:    config,
		container: c,
		rawBlock:  rawBlock,
		rawDumper: &rawdumper.RawDumper{},
	}

	entries := c.GetEntries()

	indexToInsert := -1
	similarBlocksIndexes := []int{}

	for index, entry := range entries {
		if entry.BlockDirective != nil && entry.BlockDirective.Identifier == name {
			similarBlocksIndexes = append(similarBlocksIndexes, index)
		}
	}

	if len(similarBlocksIndexes) != 0 {
		if begining {
			indexToInsert = similarBlocksIndexes[0]

			// skip block comments befor insert
			for i := similarBlocksIndexes[0] - 1; i >= 0; i-- {
				if entries[i].Comment == nil {
					break
				}

				indexToInsert = i
			}
		} else {
			indexToInsert = similarBlocksIndexes[len(similarBlocksIndexes)-1]

			if indexToInsert == len(entries)-1 {
				indexToInsert = -1
			} else {
				indexToInsert += 1
			}
		}
	}

	entry := &rawparser.Entry{
		BlockDirective: rawBlock,
		EndNewLines:    []string{"\n\n"},
	}

	if indexToInsert == -1 {
		entries = append(entries, entry)
	} else {
		if indexToInsert == 0 {
			entry.StartNewLines = []string{"\n"}
		}
		entries = slices.Insert(entries, indexToInsert, entry)
	}

	setEntries(c, entries)

	return block
}

func deleteBlock(c entryContainer, block Block) {
	deleteBlockEntityContainer(c, func(rawBlock *rawparser.BlockDirective) bool {
		return block.rawBlock == rawBlock
	})
}

func deleteBlockEntityContainer(c entryContainer, callback func(block *rawparser.BlockDirective) bool) {
	entries := c.GetEntries()
	dEntries := []*rawparser.Entry{}
	indexesToDelete := []int{}

	for index, entry := range entries {
		if entry.BlockDirective == nil {
			continue
		}

		if callback(entry.BlockDirective) {
			indexesToDelete = append(indexesToDelete, index)

			continue
		}

		deleteBlockEntityContainer(entry.BlockDirective, callback)
	}

	for index, entry := range entries {
		if !slices.Contains(indexesToDelete, index) {
			dEntries = append(dEntries, entry)
		}
	}

	setEntries(c, dEntries)
}

func addDirectoryBlock(b *Block, isRegex bool, match string, begining bool) DirectoryBlock {
	parameters := []string{}

	if isRegex {
		parameters = append(parameters, "~")
	}

	if match != "" {
		parameters = append(parameters, match)
	}

	block := b.AddBlock(string(Directory), parameters, begining)

	return DirectoryBlock{
		Block: block,
	}
}
