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
