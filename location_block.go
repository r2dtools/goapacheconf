package goapacheconf

type LocationBlock struct {
	Block
}

func (b *LocationBlock) GetLocation() string {
	return b.GetFirstParameter()
}

func (b *LocationBlock) SetLocation(location string) {
	parameters := b.GetParameters()
	parameters[0] = location

	b.SetParameters(parameters)
}
