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

func (c *ConfigFile) FindDirectives(directiveName DirectiveName) []Directive {
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

func (c *ConfigFile) FindBlocks(blockName BlockName) []Block {
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

func (c *ConfigFile) AddDirective(name string, values []string, begining bool, endWithNewLine bool) Directive {
	return newDirective(c.configFile, name, values, nil, begining, endWithNewLine)
}

func (c *ConfigFile) DeleteVirtualHostBlock(virtualHostBlock VirtualHostBlock) {
	deleteBlock(c.configFile, virtualHostBlock.Block)
}

func (c *ConfigFile) AddBlock(name string, parameters []string) Block {
	return newBlock(c.configFile, c.config, name, parameters, nil, false)
}

func (c *ConfigFile) AddAliasDirective(fromLocation, toLocation string) AliasDirective {
	directive := c.AddDirective(Alias, []string{fromLocation, toLocation}, false, true)

	return AliasDirective{Directive: directive}
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
