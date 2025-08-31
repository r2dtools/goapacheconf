package goapacheconf

type DirectoryBlock struct {
	Block
}

func (d *DirectoryBlock) IsRegex() bool {
	parameters := d.GetParameters()

	return len(parameters) > 1
}

func (d *DirectoryBlock) GetLocationMatch() string {
	parameters := d.GetParameters()

	if len(parameters) > 1 {
		return parameters[1]
	}

	if len(parameters) == 1 {
		return parameters[0]
	}

	return ""
}

func (d *DirectoryBlock) SetLocationMatch(match string) {
	parameters := d.GetParameters()

	if len(parameters) > 1 {
		parameters[1] = match
	} else {
		parameters[0] = match
	}

	d.SetParameters(parameters)
}
