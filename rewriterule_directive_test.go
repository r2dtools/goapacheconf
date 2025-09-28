package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRelatedRewriteCondDirectives(t *testing.T) {
	config := getConfig(t)
	vHostBlocks := config.FindVirtualHostBlocksByServerName("webmail.r2dtools.work.gd")
	require.Len(t, vHostBlocks, 2)

	vHostBlock := vHostBlocks[0]
	directives := vHostBlock.FindRewriteRuleDirectives()
	require.Len(t, directives, 1)

	directive := directives[0]
	rcDirectives := directive.GetRelatedRewiteCondDirectives()
	require.Len(t, rcDirectives, 2)
}
