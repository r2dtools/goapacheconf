package goapacheconf

import "slices"

type blockLocator interface {
	FindBlocks(blockName string) []Block
}

type virtualHostBlockLocator interface {
	FindVirtualHostBlocks() []VirtualHostBlock
}

func findVirtualhostBlocks(locator blockLocator) []VirtualHostBlock {
	var blocks []VirtualHostBlock

	for _, block := range locator.FindBlocks("VirtualHost") {
		blocks = append(blocks, VirtualHostBlock{
			Block: block,
		})
	}

	return blocks
}

func findVirtualHostBlocksByServerName(locator virtualHostBlockLocator, serverName string) []VirtualHostBlock {
	var blocks []VirtualHostBlock

	for _, block := range locator.FindVirtualHostBlocks() {
		serverNames := block.GetServerNames()

		if slices.Contains(serverNames, serverName) {
			blocks = append(blocks, block)
		}
	}

	return blocks
}
