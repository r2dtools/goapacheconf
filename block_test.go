package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlockIfModules(t *testing.T) {
	config := getConfig(t)

	vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, vBlocks, 2)

	vBlock := vBlocks[0]
	require.Equal(t, []string{"ssl"}, vBlock.IfModules)
	blocks := vBlock.FindBlocks(Proxy)
	require.Len(t, blocks, 1)

	block := blocks[0]
	require.Equal(t, []string{"ssl", "proxy_http"}, block.IfModules)
}

func TestFindIfModuleBlocks(t *testing.T) {
	config := getConfig(t)

	vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, vBlocks, 2)

	vBlock := vBlocks[1]
	blocks := vBlock.FindIfModuleBlocks()
	require.Len(t, blocks, 2)

	iBlock := blocks[0]
	blocks = iBlock.FindIfModuleBlocksByModuleName("rewrite")
	require.Len(t, blocks, 1)
	require.Equal(t, []string{"proxy_http"}, blocks[0].IfModules)

	blocks = vBlock.FindIfModuleBlocksByModuleName("rewrite")
	require.Len(t, blocks, 1)
}

func TestAddDirectiveToBlock(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		configFile := getConfigFile(t, r2dtoolsConfigFileName)
		vBlocks := configFile.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)

		vBlock := vBlocks[0]
		directive := NewDirective("Test", []string{"test"})
		directive.AppendNewLine()
		directive = vBlock.PrependDirective(directive)
		_, err := configFile.Dump()
		require.Nil(t, err)

		configFile = getConfigFile(t, r2dtoolsConfigFileName)
		vBlocks = configFile.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)
		vBlock = vBlocks[0]

		directives := vBlock.FindDirectives("Test")
		require.Len(t, directives, 1)
		require.Equal(t, []string{"test"}, directives[0].GetValues())
	})
}

func TestDeleteDirectiveFromBlock(t *testing.T) {
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

func TestGetOrder(t *testing.T) {
	config := getConfig(t)

	vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, vBlocks, 2)

	vBlock := vBlocks[0]
	directives := vBlock.FindAlliasDirectives()
	require.NotEmpty(t, directives)

	directive := directives[0]
	order := vBlock.GetDirectiveOrder(directive.Directive)
	require.Equal(t, 13, order)

	blocks := vBlock.FindDirectoryBlocks()
	require.NotEmpty(t, blocks)
	block := blocks[0]

	order = vBlock.GetBlockOrder(block.Block)
	require.Equal(t, 24, order)
}
