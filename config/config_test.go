package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHTTPEndpoint(t *testing.T) {
	// Let's try by default.
	want := "http://localhost:8080"
	a, _ := New()
	got := a.GetHTTPEndpoint()
	assert.Equal(t, want, got)
	tls := TLS{
		DomainNames: "domain.name.com,another.domain.name.com",
		Enable:      true,
	}
	a, _ = Get(WithTLS(tls))
	want = "https://domain.name.com:443"
	got = a.GetHTTPEndpoint()
	assert.Equal(t, want, got)
}

func TestTLSLetsEncryptVerify(t *testing.T) {
	// Happy path - all good.
	dir, err := os.MkdirTemp("", "example")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)
	tls := TLS{
		CacheDir:    dir,
		DomainNames: "darron.froese.org",
		Email:       "darron@froese.org",
		Enable:      true,
	}
	err = tls.LetsEncryptVerify()
	assert.NoError(t, err)
	// Create one - but it won't be there:
	dir, _ = os.MkdirTemp("", "example")
	// Because we'll delete it.
	os.RemoveAll(dir)
	// It will be recreated - so let's make sure it's deleted.
	defer os.RemoveAll(dir)
	tls.CacheDir = dir
	err = tls.LetsEncryptVerify()
	assert.NoError(t, err)
	// Let's remove the email address.
	tls.Email = ""
	err = tls.LetsEncryptVerify()
	assert.Error(t, err)
	// Let's remove the domain name.
	tls.DomainNames = ""
	err = tls.LetsEncryptVerify()
	assert.Error(t, err)
	// Let's remove the directory.
	tls.CacheDir = ""
	err = tls.LetsEncryptVerify()
	assert.Error(t, err)
}
