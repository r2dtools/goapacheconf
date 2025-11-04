package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindDirectivesInVirtualHostBlock(t *testing.T) {
	config := getConfig(t)

	blocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, blocks, 2)

	block := blocks[0]
	directives := block.FindDirectives("ErrorLog")
	require.Len(t, directives, 1)

	directive := directives[0]
	require.Equal(t, "ErrorLog", directive.GetName())
	require.ElementsMatch(t, directive.GetValues(), []string{`"/var/www/vhosts/system/r2dtools.work.gd/logs/error_log"`})

	directives = block.FindDirectives("DirectoryIndex")
	require.Len(t, directives, 1)
	require.Len(t, directives[0].GetValues(), 7)

	directoryBlocks := block.FindDirectoryBlocks()
	require.Len(t, directoryBlocks, 3)
}

func TestVirtualHostBlock(t *testing.T) {
	config := getConfig(t)
	blocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")

	require.Len(t, blocks, 2)

	block := blocks[0]
	require.Equal(t, "VirtualHost", block.GetName())
	require.ElementsMatch(t, block.GetServerNames(), []string{"r2dtools.work.gd"})
	require.ElementsMatch(t, block.GetServerAliases(), []string{"www.r2dtools.work.gd", "ipv4.r2dtools.work.gd"})
	require.Equal(t, true, block.HasSSL())
	require.Equal(t, `"/var/www/vhosts/r2dtools.work.gd/httpdocs"`, block.GetDocumentRoot())
}

func TestAddDirectiveInVirtualHostBlock(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		config, block := getFirstVirtualHostBlock(t, "r2dtools.work.gd")
		directive := NewDirective("test", []string{"test_value"})
		directive.AppendNewLine()
		directive = block.PrependDirective(directive)
		err := config.Dump()
		require.Nil(t, err)

		_, directive = getVirtualHostBlockFirstDirective(t, "r2dtools.work.gd", "test")
		require.Equal(t, "test_value", directive.GetFirstValue())
	})
}

func TestDeleteDirectiveByNameInVirtualHostBlock(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		config, block := getFirstVirtualHostBlock(t, "r2dtools.work.gd")
		block.DeleteDirectiveByName("CustomLog")
		err := config.Dump()
		require.Nil(t, err)

		_, block = getFirstVirtualHostBlock(t, "r2dtools.work.gd")
		directives := block.FindDirectives("CustomLog")
		require.Empty(t, directives)
	})
}

func TestDeleteDirectiveInVirtualHostBlock(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		config, block := getFirstVirtualHostBlock(t, "r2dtools.work.gd")
		directives := block.FindDirectives("Alias")
		require.Len(t, directives, 8)

		block.DeleteDirective(directives[2])
		err := config.Dump()
		require.Nil(t, err)

		_, block = getFirstVirtualHostBlock(t, "r2dtools.work.gd")
		directives = block.FindDirectives("Alias")
		require.Len(t, directives, 7)
		require.Equal(t, []string{"/ftpstat", "/var/www/vhosts/system/r2dtools.work.gd/statistics/ftpstat"}, directives[2].GetValues())
	})
}

func getVirtualHostBlockFirstDirective(t *testing.T, serverName string, directiveName string) (*Config, Directive) {
	config, block := getFirstVirtualHostBlock(t, serverName)
	directives := block.FindDirectives(directiveName)
	require.Len(t, directives, 1)

	return config, directives[0]
}

func getFirstVirtualHostBlock(t *testing.T, serverName string) (*Config, VirtualHostBlock) {
	config := getConfig(t)

	blocks := config.FindVirtualHostBlocksByServerName(serverName)
	require.NotEmpty(t, blocks)

	return config, blocks[0]
}
