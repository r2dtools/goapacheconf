package goapacheconf

import (
	"os"

	"github.com/r2dtools/goapacheconf/internal/rawparser"
)

type ConfigFile struct {
	FilePath   string
	configFile *rawparser.Config
	config     *Config
}

func (c *ConfigFile) FindDirectives(directiveName string) []Directive {
	var directives []Directive

	for _, entry := range c.configFile.GetEntries() {
		directives = append(directives, c.config.findDirectivesRecursively(directiveName, c.configFile, entry, true)...)
	}

	return directives
}

func (c *ConfigFile) FindBlocks(blockName string) []Block {
	var blocks []Block

	for _, entry := range c.configFile.GetEntries() {
		blocks = append(blocks, c.config.findBlocksRecursively(blockName, c.FilePath, c.configFile, entry, true)...)
	}

	return blocks
}

func (c *ConfigFile) FindVirtualHostBlocks() []VirtualHostBlock {
	return findVirtualhostBlocks(c)
}

func (c *ConfigFile) FindVirtualHostBlocksByServerName(serverName string) []VirtualHostBlock {
	return findVirtualHostBlocksByServerName(c, serverName)
}

func (c *ConfigFile) DeleteDirective(directive Directive) {
	deleteDirective(c.configFile, directive)
}

func (c *ConfigFile) DeleteDirectiveByName(directiveName string) {
	deleteDirectiveByName(c.configFile, directiveName)
}

func (c *ConfigFile) AddDirective(directive Directive, begining bool, endWithNewLine bool) {
	addDirective(c.configFile, directive, begining, endWithNewLine)
}

func (c *ConfigFile) DeleteVirtualHostBlock(virtualHostBlock VirtualHostBlock) {
	deleteBlock(c.configFile, virtualHostBlock.Block)
}

func (c *ConfigFile) AddBlock(name string, parameters []string) Block {
	return newBlock(c.configFile, c.config, name, parameters, false)
}

func (c *ConfigFile) Dump() error {
	content, err := c.config.rawDumper.Dump(c.configFile)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(c.FilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(content)

	if err != nil {
		return err
	}

	return nil
}
