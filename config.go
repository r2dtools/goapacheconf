package goapacheconf

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/r2dtools/goapacheconf/internal/rawdumper"
	"github.com/r2dtools/goapacheconf/internal/rawparser"
	"github.com/unknwon/com"
)

var ErrInvalidDirective = errors.New("entry is not a directive")
var ErrInvalidBlock = errors.New("entry is not a block")

type Config struct {
	rawParser   *rawparser.RawParser
	rawDumper   *rawdumper.RawDumper
	parsedFiles map[string]*rawparser.Config
	serverRoot  string
	configRoot  string
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
	return identifier == "include" || identifier == "includeoptional"
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
