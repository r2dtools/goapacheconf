package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindDirective(t *testing.T) {
	config := getConfig(t)

	directives := config.FindDirectives("SSLCertificateFile")
	assert.Len(t, directives, 3)

	directive := directives[2]
	assert.Equal(t, "SSLCertificateFile", directive.GetName())
	assert.Equal(t, "/opt/psa/var/certificates/cert9Jcn6w4", directive.GetFirstValue())
}

func TestDirectiveChangeValue(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		certPath := "/path/to/certificate"

		config, directive := getVirtualHostBlockFirstDirective(t, "r2dtools.work.gd", "SSLCertificateFile")

		directive.SetValue(certPath)
		err := config.Dump()
		require.Nil(t, err)

		config, directive = getVirtualHostBlockFirstDirective(t, "r2dtools.work.gd", "SSLCertificateFile")
		require.Equal(t, certPath, directive.GetFirstValue())
	})
}
