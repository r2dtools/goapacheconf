package goapacheconf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	config := getConfig(t)
	require.Len(t, config.parsedFiles, 102)
}

func TestGetEnabledModules(t *testing.T) {
	config := getConfig(t)
	modules := config.GetEnabledModules()
	require.Len(t, modules, 37)
}

func TestFindDirectives(t *testing.T) {
	config := getConfig(t)
	directives := config.FindDirectives("SSLEngine")
	require.Len(t, directives, 4)
}

func TestFindblocks(t *testing.T) {
	config := getConfig(t)
	blocks := config.FindBlocks("VirtualHost")
	require.Len(t, blocks, 9)
}

func getConfig(t *testing.T) *Config {
	config, err := GetConfig("./test/apache2", "")
	require.Nil(t, err)

	return config
}

func testWithConfigFileRollback(t *testing.T, configFilePath string, testFunc func(t *testing.T)) {
	configFileContent, err := os.ReadFile(configFilePath)
	assert.Nil(t, err)

	testFunc(t)

	err = os.WriteFile(configFilePath, configFileContent, 0666)
	assert.Nil(t, err)
}
