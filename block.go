package goapacheconf

import (
	"slices"

	"github.com/r2dtools/goapacheconf/internal/rawdumper"
	"github.com/r2dtools/goapacheconf/internal/rawparser"
)

type Block struct {
	FilePath  string
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

func (b *Block) SetParameters(parameters []string) {
	b.rawBlock.SetParameters(parameters)
}

func (b *Block) FindDirectives(directiveName string) []Directive {
	var directives []Directive

	for _, entry := range b.rawBlock.GetEntries() {
		directives = append(directives, b.config.findDirectivesRecursively(directiveName, b.rawBlock, entry, true)...)
	}

	return directives
}

func (b *Block) FindBlocks(blockName string) []Block {
	var blocks []Block

	for _, entry := range b.rawBlock.GetEntries() {
		blocks = append(blocks, b.config.findBlocksRecursively(blockName, b.FilePath, b.rawBlock, entry, true)...)
	}

	return blocks
}

func (b *Block) AddDirective(directive Directive, begining bool, endWithNewLine bool) {
	addDirective(b.rawBlock, directive, begining, endWithNewLine)
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
	return newBlock(b.rawBlock, b.config, name, parameters, begining)
}

func newBlock(c entryContainer, config *Config, name string, parameters []string, begining bool) Block {
	rawBlock := &rawparser.BlockDirective{
		Identifier: name,
		Content:    &rawparser.BlockContent{},
	}
	rawBlock.SetParameters(parameters)

	block := Block{
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

func findDirectoryBlocks(locator blockLocator) []DirectoryBlock {
	var directoryBlocks []DirectoryBlock

	for _, block := range locator.FindBlocks("Directory") {
		directoryBlocks = append(directoryBlocks, DirectoryBlock{
			Block: block,
		})
	}

	return directoryBlocks
}

func addDirectiveBlock(b *Block, isRegex bool, match string, begining bool) DirectoryBlock {
	parameters := []string{}

	if isRegex {
		parameters = append(parameters, "~")
	}

	if match != "" {
		parameters = append(parameters, match)
	}

	block := b.AddBlock("Directive", parameters, begining)

	return DirectoryBlock{
		Block: block,
	}
}
