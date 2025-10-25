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
		configFile.AddAliasDirective("/from", "/to")

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
