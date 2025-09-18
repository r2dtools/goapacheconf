package goapacheconf

import "slices"

type blockLocator interface {
	FindBlocks(blockName BlockName) []Block
}

type virtualHostBlockLocator interface {
	FindVirtualHostBlocks() []VirtualHostBlock
}

type ifModuleBlockLocator interface {
	FindIfModuleBlocks() []IfModuleBlock
}

func findVirtualHostBlocks(locator blockLocator) []VirtualHostBlock {
	var blocks []VirtualHostBlock

	for _, block := range locator.FindBlocks(VirtualHost) {
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

func findIfModuleBlocks(locator blockLocator) []IfModuleBlock {
	var blocks []IfModuleBlock

	for _, block := range locator.FindBlocks(IfModule) {
		blocks = append(blocks, IfModuleBlock{
			Block: block,
		})
	}

	return blocks
}

func findIfModuleBlocksByModuleName(locator ifModuleBlockLocator, moduleName string) []IfModuleBlock {
	var blocks []IfModuleBlock

	for _, block := range locator.FindIfModuleBlocks() {
		if moduleName == block.GetModuleName() {
			blocks = append(blocks, block)
		}
	}

	return blocks
}
