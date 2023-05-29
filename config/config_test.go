package config

import (
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
