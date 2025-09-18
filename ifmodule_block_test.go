package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddDirectiveToIfModuleBlock(t *testing.T) {
	config := getConfig(t)

	vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, vBlocks, 2)

	vBlock := vBlocks[1]
	blocks := vBlock.FindIfModuleBlocksByModuleName("proxy_http")
	require.Len(t, blocks, 1)

	block := blocks[0]
	directive := block.AddDirective("Test", []string{"test"}, true, true)
	require.Equal(t, []string{"proxy_http"}, directive.IfModules)

	nBlock := block.AddBlock("Test", []string{"test"}, true)
	require.Equal(t, []string{"proxy_http"}, nBlock.IfModules)
}
