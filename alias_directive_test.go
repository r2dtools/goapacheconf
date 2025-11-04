package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindAliasDirective(t *testing.T) {
	config := getConfig(t)

	directives := config.FindAliasDirectives()
	require.Len(t, directives, 12)

	acmeAlias := directives[9]
	require.Equal(t, "/.well-known/acme-challenge", acmeAlias.GetFromLocation())
	require.Equal(t, "\"/var/www/vhosts/default/htdocs/.well-known/acme-challenge\"", acmeAlias.GetToLocation())
}

func TestAddAliasDirective(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		configFile := getConfigFile(t, r2dtoolsConfigFileName)
		aliasDirective := NewAliasDirective("/from", "/to")
		aliasDirective.AppendNewLine()
		configFile.AppendAliasDirective(aliasDirective)

		_, err := configFile.Dump()
		require.Nil(t, err)

		configFile = getConfigFile(t, r2dtoolsConfigFileName)
		aliasDirectives := configFile.FindAlliasDirectives()

		var newAliasDirective *AliasDirective

		for _, aliasDirective := range aliasDirectives {
			if aliasDirective.GetFromLocation() == "/from" && aliasDirective.GetToLocation() == "/to" {
				newAliasDirective = &aliasDirective

				break
			}
		}

		require.NotNil(t, newAliasDirective)
	})
}

func TestDeleteAliasDirective(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		configFile := getConfigFile(t, r2dtoolsConfigFileName)
		directives := configFile.FindAlliasDirectives()
		require.Len(t, directives, 8)

		virtualHostBlocks := configFile.FindVirtualHostBlocks()
		require.Len(t, virtualHostBlocks, 2)

		virtualHostBlock := virtualHostBlocks[0]
		virtualHostBlock.DeleteAliasDirective(directives[0])
		_, err := configFile.Dump()
		require.Nil(t, err)

		configFile = getConfigFile(t, r2dtoolsConfigFileName)
		directives = configFile.FindAlliasDirectives()
		require.Len(t, directives, 7)
	})
}
