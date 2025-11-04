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
	var (
		directives []Directive
		ifModules  []string
	)

	for _, entry := range c.configFile.GetEntries() {
		directives = append(directives, c.config.findDirectivesRecursively(directiveName, c.FilePath, c.configFile, entry, ifModules, true)...)
	}

	return directives
}

func (c *ConfigFile) FindRewriteRuleDirectives() []RewriteRuleDirective {
	return findDirectives[RewriteRuleDirective](c, RewriteRule)
}

func (c *ConfigFile) FindAlliasDirectives() []AliasDirective {
	return findDirectives[AliasDirective](c, Alias)
}

func (c *ConfigFile) FindBlocks(blockName string) []Block {
	var (
		blocks    []Block
		ifModules []string
	)

	for _, entry := range c.configFile.GetEntries() {
		blocks = append(blocks, c.config.findBlocksRecursively(blockName, c.FilePath, c.configFile, entry, ifModules, true)...)
	}

	return blocks
}

func (c *ConfigFile) FindVirtualHostBlocks() []VirtualHostBlock {
	return findBlocks[VirtualHostBlock](c, VirtualHost)
}

func (c *ConfigFile) FindVirtualHostBlocksByServerName(serverName string) []VirtualHostBlock {
	return findVirtualHostBlocksByServerName(c, serverName)
}

func (c *ConfigFile) FindIfModuleBlocks() []IfModuleBlock {
	return findBlocks[IfModuleBlock](c, IfModule)
}

func (c *ConfigFile) FindIfModuleBlocksByModuleName(moduleName string) []IfModuleBlock {
	return findIfModuleBlocksByModuleName(c, moduleName)
}

func (c *ConfigFile) DeleteDirective(directive Directive) {
	deleteDirective(c.configFile, directive)
}

func (c *ConfigFile) DeleteDirectiveByName(directiveName string) {
	deleteDirectiveByName(c.configFile, directiveName)
}

func (c *ConfigFile) AppendDirective(directive Directive) Directive {
	entries := c.configFile.GetEntries()
	directive.setContainer(c.configFile)

	var prevEntry *rawparser.Entry

	if len(entries) > 0 {
		prevEntry = entries[len(entries)-1]
	}

	if prevEntry == nil || len(prevEntry.EndNewLines) == 0 {
		directive.PrependNewLine()
	}

	entries = append(entries, directive.entry)

	setEntries(c.configFile, entries)

	return directive
}

func (c *ConfigFile) PrependDirective(directive Directive) Directive {
	entries := c.configFile.GetEntries()
	directive.setContainer(c.configFile)

	directive.PrependNewLine()
	entries = append([]*rawparser.Entry{directive.entry}, entries...)

	setEntries(c.configFile, entries)

	return directive
}

func (c *ConfigFile) DeleteVirtualHostBlock(virtualHostBlock VirtualHostBlock) {
	deleteBlock(c.configFile, virtualHostBlock.Block)
}

func (c *ConfigFile) AddBlock(name string, parameters []string, begining bool) Block {
	return newBlock(c.configFile, c.config, name, parameters, nil, begining)
}

func (c *ConfigFile) AppendAliasDirective(aliasDirective AliasDirective) AliasDirective {
	directive := c.AppendDirective(aliasDirective.Directive)

	return AliasDirective{Directive: directive}
}

func (c *ConfigFile) PrependAliasDirective(aliasDirective AliasDirective) AliasDirective {
	directive := c.PrependDirective(aliasDirective.Directive)

	return AliasDirective{Directive: directive}
}

func (c *ConfigFile) DeleteAliasDirective(aliasDirective AliasDirective) {
	deleteDirective(c.configFile, aliasDirective.Directive)
}

func (c *ConfigFile) Dump() (string, error) {
	content, err := c.config.rawDumper.Dump(c.configFile)

	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(c.FilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)

	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = file.WriteString(content)

	if err != nil {
		return "", err
	}

	return content, nil
}
