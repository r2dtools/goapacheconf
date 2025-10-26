package goapacheconf

type LocationMatchBlock struct {
	Block
}

func (b *LocationMatchBlock) GetLocationMatch() string {
	return b.GetFirstParameter()
}

func (b *LocationMatchBlock) SetLocationMatch(locationMatch string) {
	parameters := b.GetParameters()
	parameters[0] = locationMatch

	b.SetParameters(parameters)
}
