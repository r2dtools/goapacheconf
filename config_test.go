package goapacheconf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var r2dtoolsConfigFilePath = "test/apache2/plesk.conf.d/vhosts/r2dtools.work.gd.conf"
var r2dtoolsConfigFileName = "r2dtools.work.gd.conf"

func TestGetConfig(t *testing.T) {
	config := getConfig(t)
	require.Len(t, config.parsedFiles, 102)
}

func TestGetEnabledModules(t *testing.T) {
	config := getConfig(t)
	modules := config.GetEnabledModules()
	require.Len(t, modules, 37)
}

func TestIsModuleEnabled(t *testing.T) {
	config := getConfig(t)
	require.True(t, config.IsModuleEnabled("ssl"))
}

func TestFindDirectives(t *testing.T) {
	config := getConfig(t)
	directives := config.FindDirectives(SSLEngine)
	require.Len(t, directives, 4)
}

func TestFindblocks(t *testing.T) {
	config := getConfig(t)
	blocks := config.FindBlocks(VirtualHost)
	require.Len(t, blocks, 9)
}

func TestAddConfigFile(t *testing.T) {
	configFilePath := "./test/apache2/sites-enabled/example.com.conf"

	config := getConfig(t)
	configFile, err := config.AddConfigFile(configFilePath)
	require.Nil(t, err)

	directive := NewDirective("TestDirective", []string{"test"})
	configFile.AddDirective(directive, true, true)

	block := configFile.AddBlock("TestBlock", []string{"test"})
	directive = NewDirective("TestBlockDirective", []string{"test", "directive"})
	block.AddDirective(directive, false, true)

	err = configFile.Dump()
	require.Nil(t, err)
	defer os.Remove(configFilePath)

	config = getConfig(t)
	configFile = config.GetConfigFile("example.com.conf")
	require.NotNil(t, configFile)
	directives := configFile.FindDirectives("TestDirective")
	require.Len(t, directives, 1)
	require.Equal(t, "TestDirective", directives[0].GetName())
	require.Equal(t, []string{"test"}, directives[0].GetValues())

	blocks := configFile.FindBlocks("TestBlock")
	require.Len(t, blocks, 1)
}

func getConfig(t *testing.T) *Config {
	config, err := GetConfig("./test/apache2", "")
	require.Nil(t, err)

	return config
}

func testWithConfigFileRollback(t *testing.T, configFilePath string, testFunc func(t *testing.T)) {
	configFileContent, err := os.ReadFile(configFilePath)
	require.Nil(t, err)

	defer func() {
		err = os.WriteFile(configFilePath, configFileContent, 0666)

		if err != nil {
			t.Log(err)
		}
	}()

	testFunc(t)
}
