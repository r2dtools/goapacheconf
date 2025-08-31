package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirective(t *testing.T) {
	config := getConfig(t)

	directives := config.FindDirectives("SSLCertificateFile")
	assert.Len(t, directives, 3)

	directive := directives[2]
	assert.Equal(t, "SSLCertificateFile", directive.GetName())
	assert.Equal(t, "/opt/psa/var/certificates/cert9Jcn6w4", directive.GetFirstValue())
}
