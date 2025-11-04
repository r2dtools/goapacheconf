package goapacheconf

import "slices"

type BlockUnion interface {
	VirtualHostBlock | IfModuleBlock | LocationBlock | DirectoryBlock | LocationMatchBlock
}

type blockLocator interface {
	FindBlocks(blockName string) []Block
}

type virtualHostBlockLocator interface {
	FindVirtualHostBlocks() []VirtualHostBlock
}

type ifModuleBlockLocator interface {
	FindIfModuleBlocks() []IfModuleBlock
}

func findBlocks[T BlockUnion](locator blockLocator, name string) []T {
	var blocks []T

	for _, block := range locator.FindBlocks(name) {
		blocks = append(blocks, T{Block: block})
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

func findIfModuleBlocksByModuleName(locator ifModuleBlockLocator, moduleName string) []IfModuleBlock {
	var blocks []IfModuleBlock

	for _, block := range locator.FindIfModuleBlocks() {
		if moduleName == block.GetModuleName() {
			blocks = append(blocks, block)
		}
	}

	return blocks
}
