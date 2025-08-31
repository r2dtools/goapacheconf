package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	config := getConfig(t)
	require.Len(t, config.parsedFiles, 102)
}

func getConfig(t *testing.T) *Config {
	config, err := GetConfig("./test/apache2", "")
	require.Nil(t, err)

	return config
}
