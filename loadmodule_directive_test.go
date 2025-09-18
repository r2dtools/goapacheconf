package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetModuleName(t *testing.T) {
	config := getConfig(t)
	directives := config.FindLoadModuleDirectives()
	require.NotEmpty(t, directives)
	require.Equal(t, "access_compat", directives[0].GetModuleName())
}
