package rawparser

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootConfParse(t *testing.T) {
	items := []struct {
		configName   string
		expectedName string
	}{
		{"apache2.conf", "apache2.conf.json"},
		{"r2dtools.work.gd.conf", "r2dtools.work.gd.conf.json"},
		{"webmail.r2dtools.work.gd.conf", "webmail.r2dtools.work.gd.conf.json"},
	}

	parser, err := GetRawParser()
	assert.Nilf(t, err, "could not create parser: %v", err)

	for _, item := range items {
		content, err := os.ReadFile("../../test/" + item.configName)
		assert.Nilf(t, err, "could not read config file")

		parsedConfig, err := parser.Parse(string(content))
		assert.Nilf(t, err, "could not parse config: %v", err)

		expectedData := &Config{}
		data, err := os.ReadFile("../../test/" + item.expectedName)
		assert.Nilf(t, err, "could not read file with expected data: %v", err)

		err = json.Unmarshal(data, expectedData)
		assert.Nilf(t, err, "could not decode expected data: %v", err)

		assert.Equal(t, expectedData, parsedConfig, "parsed data is invalid")
	}

}
