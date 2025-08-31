package goapacheconf

import (
	"github.com/r2dtools/goapacheconf/internal/container"
	"github.com/r2dtools/goapacheconf/internal/rawdumper"
	"github.com/r2dtools/goapacheconf/internal/rawparser"
)

type Block struct {
	FilePath  string
	config    *Config
	container container.EntryContainer
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

func (b *Block) Dump() string {
	entry := rawparser.Entry{
		BlockDirective: b.rawBlock,
	}

	return b.rawDumper.DumpEntry(&entry, false)
}
