package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindVirtualHostBlocksByNameInConfigFile(t *testing.T) {
	configFile := getConfigFile(t, r2dtoolsConfigFileName)
	blocks := configFile.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, blocks, 2)
}

func TestDeleteDirectiveByNameInConfigFile(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		config, block := getFirstVirtualHostBlock(t, "r2dtools.work.gd")
		block.DeleteDirectiveByName(UseCanonicalName)
		err := config.Dump()
		require.Nil(t, err)

		_, block = getFirstVirtualHostBlock(t, "r2dtools.work.gd")
		directives := block.FindDirectives(UseCanonicalName)
		require.Empty(t, directives)
	})
}

func TestAddDirectiveInConfigFile(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		configFile := getConfigFile(t, r2dtoolsConfigFileName)
		directive := NewDirective("test", []string{"test_value"})
		configFile.AddDirective(directive, true, true)
		err := configFile.Dump()
		require.Nil(t, err)

		configFile = getConfigFile(t, r2dtoolsConfigFileName)
		directives := configFile.FindDirectives("test")
		require.Len(t, directives, 1)
		require.Equal(t, "test_value", directives[0].GetFirstValue())
	})
}

func TestDeleteBlock(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		configFile := getConfigFile(t, r2dtoolsConfigFileName)
		blocks := configFile.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, blocks, 2)
		block := blocks[0]

		configFile.DeleteVirtualHostBlock(block)
		err := configFile.Dump()
		assert.Nil(t, err)

		configFile = getConfigFile(t, r2dtoolsConfigFileName)
		blocks = configFile.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		assert.Len(t, blocks, 1)
	})
}

func getConfigFile(t *testing.T, name string) *ConfigFile {
	config := getConfig(t)

	configFile := config.GetConfigFile(name)
	require.NotNil(t, configFile)

	return configFile
}
