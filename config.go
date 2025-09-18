package goapacheconf

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/r2dtools/goapacheconf/internal/rawdumper"
	"github.com/r2dtools/goapacheconf/internal/rawparser"
	"github.com/unknwon/com"
)

var ErrInvalidDirective = errors.New("entry is not a directive")
var ErrInvalidBlock = errors.New("entry is not a block")

type Config struct {
	rawParser      *rawparser.RawParser
	rawDumper      *rawdumper.RawDumper
	parsedFiles    map[string]*rawparser.Config
	enabledModules map[string]struct{}
	serverRoot     string
	configRoot     string
}

func (c *Config) GetConfigFile(configFileName string) *ConfigFile {
	for configFilePath, config := range c.parsedFiles {
		pConfigFileName := filepath.Base(configFilePath)

		if configFileName == pConfigFileName {
			return &ConfigFile{
				FilePath:   configFilePath,
				configFile: config,
				config:     c,
			}
		}
	}

	return nil
}

func (c *Config) Dump() error {
	for filePath, config := range c.parsedFiles {
		content, err := c.rawDumper.Dump(config)

		if err != nil {
			return err
		}

		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)

		if err != nil {
			return err
		}

		defer file.Close()

		_, err = file.WriteString(content)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) ParseFile(configPath string) error {
	return c.parseRecursively(configPath)
}

func (c *Config) GetEnabledModules() map[string]struct{} {
	if c.enabledModules == nil {
		moduleDirectives := c.FindDirectives(LoadModule)
		modules := make(map[string]struct{}, len(moduleDirectives))

		for _, moduleDirective := range moduleDirectives {
			loadModuleDirective := LoadModuleDirective{
				Directive: moduleDirective,
			}

			name := loadModuleDirective.GetModuleName()

			if name != "" {
				modules[name] = struct{}{}
			}

		}

		c.enabledModules = modules
	}

	return c.enabledModules
}

func (c *Config) IsModuleEnabled(name string) bool {
	enabledModules := c.GetEnabledModules()
	name = strings.ToLower(name)
	_, ok := enabledModules[name]

	return ok
}

func (c *Config) FindVirtualHostBlocks() []VirtualHostBlock {
	return findVirtualhostBlocks(c)
}

func (c *Config) FindVirtualHostBlocksByServerName(serverName string) []VirtualHostBlock {
	return findVirtualHostBlocksByServerName(c, serverName)
}

func (c *Config) FindBlocks(blockName BlockName) []Block {
	var blocks []Block

	keys := slices.Collect(maps.Keys(c.parsedFiles))
	sort.Strings(keys)

	for _, key := range keys {
		tree, ok := c.parsedFiles[key]

		if !ok {
			continue
		}

		for _, entry := range tree.Entries {
			blocks = append(blocks, c.findBlocksRecursively(blockName, key, tree, entry, false)...)
		}
	}

	return blocks
}

func (c *Config) FindDirectives(directiveName DirectiveName) []Directive {
	var directives []Directive

	keys := slices.Collect(maps.Keys(c.parsedFiles))
	sort.Strings(keys)

	for _, key := range keys {
		tree, ok := c.parsedFiles[key]

		if !ok {
			continue
		}

		for _, entry := range tree.GetEntries() {
			directives = append(
				directives,
				c.findDirectivesRecursively(directiveName, tree, entry, false)...,
			)
		}
	}

	return directives
}

func (c *Config) AddConfigFile(filePath string) (*ConfigFile, error) {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		configFile := ConfigFile{
			FilePath:   filePath,
			configFile: &rawparser.Config{},
			config:     c,
		}

		return &configFile, nil
	}

	return nil, fmt.Errorf("file %s already exists", filePath)
}

func (c *Config) findDirectivesRecursively(
	directiveName DirectiveName,
	container entryContainer,
	entry *rawparser.Entry,
	withInclude bool,
) []Directive {
	var directives []Directive
	directive := entry.Directive
	blockDirective := entry.BlockDirective

	if directive != nil {
		identifier := directive.Identifier

		if withInclude && c.isInclude(identifier) {
			include := c.getAbsPath(directive.GetFirstValueStr())
			includeFiles, err := filepath.Glob(include)

			if err != nil {
				return directives
			}

			for _, includePath := range includeFiles {
				includeConfig, ok := c.parsedFiles[includePath]

				if !ok {
					continue
				}

				for _, entry := range includeConfig.GetEntries() {
					directives = append(
						directives,
						c.findDirectivesRecursively(directiveName, includeConfig, entry, withInclude)...,
					)
				}
			}
		}

		if identifier == string(directiveName) {
			directives = append(directives, Directive{
				rawDirective: directive,
				container:    container,
			})

			return directives
		}
	}

	if blockDirective == nil {
		return directives
	}

	for _, bEntry := range blockDirective.GetEntries() {
		directives = append(
			directives,
			c.findDirectivesRecursively(directiveName, blockDirective, bEntry, withInclude)...,
		)
	}

	return directives
}

func (c *Config) findBlocksRecursively(
	blockName BlockName,
	path string,
	container entryContainer,
	entry *rawparser.Entry,
	withInclude bool,
) []Block {
	var blocks []Block
	directive := entry.Directive
	blockDirective := entry.BlockDirective

	if withInclude && directive != nil && c.isInclude(directive.Identifier) {
		include := c.getAbsPath(directive.GetFirstValueStr())
		includeFiles, err := filepath.Glob(include)

		if err != nil {
			return blocks
		}

		for _, includePath := range includeFiles {
			includeConfig, ok := c.parsedFiles[includePath]

			if !ok {
				continue
			}

			for _, entry := range includeConfig.Entries {
				blocks = append(
					blocks,
					c.findBlocksRecursively(blockName, includePath, includeConfig, entry, withInclude)...,
				)
			}
		}

		return blocks
	}

	if blockDirective == nil {
		return blocks

	}

	block := Block{
		FilePath:  path,
		config:    c,
		container: container,
		rawBlock:  blockDirective,
		rawDumper: &rawdumper.RawDumper{},
	}

	if blockDirective.Identifier == string(blockName) {
		blocks = append(blocks, block)
	} else {
		// blocks can be nested
		for _, blockEntry := range blockDirective.GetEntries() {
			blocks = append(
				blocks,
				c.findBlocksRecursively(blockName, path, blockDirective, blockEntry, withInclude)...,
			)
		}
	}

	return blocks
}

func (c *Config) parse() error {
	c.parsedFiles = make(map[string]*rawparser.Config)

	return c.parseRecursively(c.configRoot)
}

func (c *Config) parseRecursively(configPath string) error {
	configFilePathAbs := c.getAbsPath(configPath)
	trees, err := c.parseFilesByPath(configFilePathAbs, false)

	if err != nil {
		return err
	}

	for _, tree := range trees {
		for _, entry := range tree.Entries {
			identifier := strings.ToLower(entry.GetIdentifier())
			// Parse the top-level included file
			if c.isInclude(identifier) {
				if entry.Directive == nil {
					return ErrInvalidDirective
				}

				includePath := entry.Directive.GetFirstValueStr()

				if includePath == "" {
					continue
				}

				if err := c.parseRecursively(includePath); err != nil {
					return err
				}

				continue
			}

			// Parse included files in blocks
			if entry.BlockDirective != nil {
				includePaths, err := c.findBlockIcludePathsRecursively(entry.BlockDirective)

				if err != nil {
					return err
				}

				for _, includePath := range includePaths {
					if err := c.parseRecursively(includePath); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (c *Config) findBlockIcludePathsRecursively(block *rawparser.BlockDirective) ([]string, error) {
	var includePaths []string

	for _, entry := range block.GetEntries() {
		identifier := strings.ToLower(entry.GetIdentifier())

		if !c.isInclude(identifier) {
			continue
		}

		if entry.Directive == nil {
			return includePaths, ErrInvalidDirective
		}

		includePath := entry.Directive.GetFirstValueStr()

		if includePath != "" {
			includePaths = append(includePaths, includePath)

			continue
		}

		if entry.BlockDirective != nil {
			blockIncludePaths, err := c.findBlockIcludePathsRecursively(entry.BlockDirective)

			if err != nil {
				return includePaths, err
			}

			includePaths = append(includePaths, blockIncludePaths...)
		}
	}

	return includePaths, nil
}

func (c *Config) parseFilesByPath(path string, override bool) ([]*rawparser.Config, error) {
	var (
		filePaths []string
		err       error
	)

	if com.IsFile(path) {
		filePaths = []string{path}
	} else {
		filePaths, err = filepath.Glob(path)

		if err != nil {
			return nil, err
		}
	}

	var trees []*rawparser.Config

	for _, filePath := range filePaths {
		if _, ok := c.parsedFiles[filePath]; ok && !override {
			continue
		}

		content, err := os.ReadFile(filePath)

		if err != nil {
			return nil, err
		}

		config, err := c.rawParser.Parse(string(content))

		if err != nil {
			return nil, err
		}

		c.parsedFiles[filePath] = config
		trees = append(trees, config)
	}

	return trees, nil
}

func (c *Config) getAbsPath(path string) string {
	path = strings.Trim(path, "'\"")

	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	return filepath.Clean(filepath.Join(c.serverRoot, path))
}

func (c *Config) isInclude(identifier string) bool {
	return identifier == strings.ToLower(Include) || identifier == strings.ToLower(IncludeOptional)
}

func GetConfig(serverRootPath, configFilePath string) (*Config, error) {
	var err error

	serverRootPath, err = filepath.Abs(serverRootPath)

	if err != nil {
		return nil, err
	}

	if configFilePath == "" {
		configFilePath = path.Join(serverRootPath, "apache2.conf")
	}

	if !filepath.IsAbs(configFilePath) {
		configFilePath = filepath.Clean(filepath.Join(serverRootPath, configFilePath))
	}

	if !com.IsFile(configFilePath) {
		return nil, fmt.Errorf("could not find '%s' config file", configFilePath)
	}

	rawParser, err := rawparser.GetRawParser()

	if err != nil {
		return nil, err
	}

	parser := Config{
		rawParser:  rawParser,
		rawDumper:  &rawdumper.RawDumper{},
		serverRoot: serverRootPath,
		configRoot: configFilePath,
	}

	if err := parser.parse(); err != nil {
		return nil, err
	}

	return &parser, nil
}
