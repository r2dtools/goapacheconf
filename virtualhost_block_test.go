package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getVirtualHostBlockDirective(t *testing.T, serverName, directiveName string) (*Config, Directive) {
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
