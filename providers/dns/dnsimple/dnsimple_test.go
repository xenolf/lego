package dnsimple

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	dnsimpleLiveTest   bool
	dnsimpleOauthToken string
	dnsimpleDomain     string
	dnsimpleBaseURL    string
)

func init() {
	dnsimpleOauthToken = os.Getenv("DNSIMPLE_OAUTH_TOKEN")
	dnsimpleDomain = os.Getenv("DNSIMPLE_DOMAIN")
	dnsimpleBaseURL = "https://api.sandbox.dnsimple.com"

	if len(dnsimpleOauthToken) > 0 && len(dnsimpleDomain) > 0 {
		baseURL := os.Getenv("DNSIMPLE_BASE_URL")

		if baseURL != "" {
			dnsimpleBaseURL = baseURL
		}

		dnsimpleLiveTest = true
	}
}

func restoreEnv() {
	os.Setenv("DNSIMPLE_OAUTH_TOKEN", dnsimpleOauthToken)
	os.Setenv("DNSIMPLE_BASE_URL", dnsimpleBaseURL)
}

//
// NewDNSProvider
//

func TestNewDNSProviderValid(t *testing.T) {
	defer restoreEnv()
	os.Setenv("DNSIMPLE_OAUTH_TOKEN", "123")

	provider, err := NewDNSProvider()

	assert.NotNil(t, provider)
	assert.Equal(t, "lego", provider.client.UserAgent)
	assert.NoError(t, err)
}

func TestNewDNSProviderValidWithBaseUrl(t *testing.T) {
	defer restoreEnv()
	os.Setenv("DNSIMPLE_OAUTH_TOKEN", "123")
	os.Setenv("DNSIMPLE_BASE_URL", "https://api.dnsimple.test")

	provider, err := NewDNSProvider()

	assert.NotNil(t, provider)
	assert.NoError(t, err)

	assert.Equal(t, provider.client.BaseURL, "https://api.dnsimple.test")
}

func TestNewDNSProviderInvalidWithMissingOauthToken(t *testing.T) {
	defer restoreEnv()
	os.Setenv("DNSIMPLE_OAUTH_TOKEN", "")

	provider, err := NewDNSProvider()

	assert.Nil(t, provider)
	assert.EqualError(t, err, "DNSimple OAuth token is missing")
}

//
// NewDNSProviderCredentials
//

func TestNewDNSProviderCredentialsValid(t *testing.T) {
	provider, err := NewDNSProviderCredentials("123", "")

	assert.NotNil(t, provider)
	assert.Equal(t, "lego", provider.client.UserAgent)
	assert.NoError(t, err)
}

func TestNewDNSProviderCredentialsValidWithBaseUrl(t *testing.T) {
	provider, err := NewDNSProviderCredentials("123", "https://api.dnsimple.test")

	assert.NotNil(t, provider)
	assert.NoError(t, err)

	assert.Equal(t, provider.client.BaseURL, "https://api.dnsimple.test")
}

func TestNewDNSProviderCredentialsInvalidWithMissingOauthToken(t *testing.T) {
	provider, err := NewDNSProviderCredentials("", "")

	assert.Nil(t, provider)
	assert.EqualError(t, err, "DNSimple OAuth token is missing")
}

//
// Present
//

func TestLiveDNSimplePresent(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	provider, err := NewDNSProviderCredentials(dnsimpleOauthToken, dnsimpleBaseURL)
	assert.NoError(t, err)

	err = provider.Present(dnsimpleDomain, "", "123d==")
	assert.NoError(t, err)
}

//
// Cleanup
//

func TestLiveDNSimpleCleanUp(t *testing.T) {
	if !dnsimpleLiveTest {
		t.Skip("skipping live test")
	}

	time.Sleep(time.Second * 1)

	provider, err := NewDNSProviderCredentials(dnsimpleOauthToken, dnsimpleBaseURL)
	assert.NoError(t, err)

	err = provider.CleanUp(dnsimpleDomain, "", "123d==")
	assert.NoError(t, err)
}
