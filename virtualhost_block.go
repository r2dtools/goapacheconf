package goapacheconf

import "strings"

type VirtualHostBlock struct {
	Block
}

func (v *VirtualHostBlock) GetServerNames() []string {
	serverNames := []string{}
	directives := v.FindDirectives("ServerName")

	if len(directives) == 0 {
		return serverNames
	}

	for _, directive := range directives {
		serverNames = append(serverNames, directive.GetValues()...)
	}

	for index, serverName := range serverNames {
		serverNames[index] = strings.Trim(serverName, " \"")
	}

	return serverNames
}

func (v *VirtualHostBlock) GetServerAliases() []string {
	serverAliases := []string{}
	directives := v.FindDirectives("ServerAlias")

	if len(directives) == 0 {
		return serverAliases
	}

	for _, directive := range directives {
		serverAliases = append(serverAliases, directive.GetValues()...)
	}

	for index, serverAlias := range serverAliases {
		serverAliases[index] = strings.Trim(serverAlias, " \"")
	}

	return serverAliases
}

func (v *VirtualHostBlock) GetDocumentRoot() string {
	directives := v.FindDirectives("DocumentRoot")

	if len(directives) == 0 {
		return ""
	}

	return directives[0].GetFirstValue()
}

func (s *VirtualHostBlock) GetAddresses() []Address {
	parameters := s.GetParameters()
	addresses := []Address{}

	for _, parameter := range parameters {
		address := CreateAddressFromString(parameter)
		addresses = append(addresses, address)
	}

	return addresses
}

func (s *VirtualHostBlock) HasSSL() bool {
	addresses := s.GetAddresses()

	for _, address := range addresses {
		if address.Port == "443" {
			return true
		}
	}

	sslDirectives := s.FindDirectives("SSLEngine")

	for _, sslDirective := range sslDirectives {
		if sslDirective.GetFirstValue() == "on" {
			return true
		}
	}

	return false
}

func (v *VirtualHostBlock) IsIpv6Enabled() bool {
	addresses := v.GetAddresses()

	for _, address := range addresses {
		if address.IsIpv6 {
			return true
		}
	}

	return false
}

func (v *VirtualHostBlock) IsIpv4Enabled() bool {
	addresses := v.GetAddresses()

	if len(addresses) == 0 {
		return true
	}

	for _, address := range addresses {
		if !address.IsIpv6 {
			return true
		}
	}

	return false
}

func (v *VirtualHostBlock) FindDirectoryBlocks() []DirectoryBlock {
	return findDirectoryBlocks(&v.Block)
}

func (v *VirtualHostBlock) AddDirectoryBlock(isRegex bool, match string, begining bool) DirectoryBlock {
	return addDirectiveBlock(&v.Block, isRegex, match, begining)
}

func (v *DirectoryBlock) DeleteDirectiveBlock(directiveBlock DirectoryBlock) {
	deleteBlock(v.rawBlock, directiveBlock.Block)
}
