package container

import "github.com/r2dtools/goapacheconf/internal/rawparser"

type EntryContainer interface {
	GetEntries() []*rawparser.Entry
	SetEntries(entries []*rawparser.Entry)
}

func setEntries(c EntryContainer, entries []*rawparser.Entry) {
	entriesCount := len(entries)

	if entriesCount > 0 {
		if len(entries[0].StartNewLines) == 0 {
			entries[0].StartNewLines = []string{"\n"}
		}

		if len(entries[entriesCount-1].EndNewLines) == 0 {
			entries[entriesCount-1].EndNewLines = []string{"\n"}
		}
	}

	c.SetEntries(entries)
}
