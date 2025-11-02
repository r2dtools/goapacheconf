package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindLocationBlocks(t *testing.T) {
	config := getConfig(t)

	vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, vBlocks, 2)

	vBlock := vBlocks[0]
	blocks := vBlock.FindLocationBlocks()
	require.Len(t, blocks, 3)

	lBlock := blocks[0]
	require.Equal(t, "/plesk-stat/", lBlock.GetLocation())
}

func TestAddLocationBlock(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		configFile := getConfigFile(t, r2dtoolsConfigFileName)

		vBlocks := configFile.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)

		vBlock := vBlocks[0]
		blocks := vBlock.FindLocationBlocks()
		require.Len(t, blocks, 3)

		vBlock.AddLocationBlock("/location", false)
		_, err := configFile.Dump()
		require.Nil(t, err)

		configFile = getConfigFile(t, r2dtoolsConfigFileName)
		vBlocks = configFile.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)

		vBlock = vBlocks[0]
		blocks = vBlock.FindLocationBlocks()
		require.Len(t, blocks, 4)

		lBlock := blocks[3]
		require.Equal(t, "/location", lBlock.GetLocation())
	})
}

func TestDeleteLocationBlock(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		config := getConfig(t)

		vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)

		vBlock := vBlocks[0]
		blocks := vBlock.FindLocationBlocks()
		require.Len(t, blocks, 3)

		lBlock := blocks[0]
		require.Equal(t, "/plesk-stat/", lBlock.GetLocation())

		vBlock.DeleteLocationBlock(lBlock)
		err := config.Dump()
		require.Nil(t, err)

		config = getConfig(t)
		vBlocks = config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)

		vBlock = vBlocks[0]
		blocks = vBlock.FindLocationBlocks()
		require.Len(t, blocks, 2)
	})
}
