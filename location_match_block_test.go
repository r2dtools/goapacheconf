package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindLocationMatchBlocks(t *testing.T) {
	config := getConfig(t)

	vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, vBlocks, 2)

	vBlock := vBlocks[0]
	blocks := vBlock.FindLocationMatchBlocks()
	require.Len(t, blocks, 1)

	lBlock := blocks[0]
	require.Equal(t, "\"^/.well-known/acme-challenge/(.*/|)\\.\"", lBlock.GetLocationMatch())
}

func TestAddLocationMatchBlock(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		configFile := getConfigFile(t, r2dtoolsConfigFileName)

		vBlocks := configFile.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)

		vBlock := vBlocks[0]
		blocks := vBlock.FindLocationMatchBlocks()
		require.Len(t, blocks, 1)

		vBlock.AddLocationMatchBlock("~/location")
		_, err := configFile.Dump()
		require.Nil(t, err)

		configFile = getConfigFile(t, r2dtoolsConfigFileName)
		vBlocks = configFile.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)

		vBlock = vBlocks[0]
		blocks = vBlock.FindLocationMatchBlocks()
		require.Len(t, blocks, 2)

		lBlock := blocks[1]
		require.Equal(t, "~/location", lBlock.GetLocationMatch())
	})
}

func TestDeleteLocationMatchBlock(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		config := getConfig(t)

		vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)

		vBlock := vBlocks[0]
		blocks := vBlock.FindLocationMatchBlocks()
		require.Len(t, blocks, 1)

		lBlock := blocks[0]

		vBlock.DeleteLocationMatchBlock(lBlock)
		err := config.Dump()
		require.Nil(t, err)

		config = getConfig(t)
		vBlocks = config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)

		vBlock = vBlocks[0]
		blocks = vBlock.FindLocationMatchBlocks()
		require.Empty(t, blocks)
	})
}
