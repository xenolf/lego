// Package sonic implements a DNS provider for solving the DNS-01 challenge using Sonic.net based on DNSMadeEasy
package sonic

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
)

// Environment variables names.
const (
	envNamespace = "SONIC_"

	EnvAPIUserId    = envNamespace + "USERID"
	EnvAPIApiKey = envNamespace + "APIKEY"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

// Config is used to configure the creation of the DNSProvider.
type Config struct {
	UserId          string
	ApiKey             string
	HTTPClient         *http.Client
	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
}

// NewDefaultConfig returns a default configuration for the DNSProvider.
func NewDefaultConfig() *Config {
	return &Config{
		TTL:                env.GetOrDefaultInt(EnvTTL, dns01.DefaultTTL),
		PropagationTimeout: env.GetOrDefaultSecond(EnvPropagationTimeout, dns01.DefaultPropagationTimeout),
		PollingInterval:    env.GetOrDefaultSecond(EnvPollingInterval, dns01.DefaultPollingInterval),
		HTTPClient: &http.Client{
			Timeout: env.GetOrDefaultSecond(EnvHTTPTimeout, 10*time.Second),
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

// DNSProvider implements the challenge.Provider interface.
type DNSProvider struct {
	config *Config
	client *Client
}

// NewDNSProvider returns a DNSProvider instance configured for Sonic.
// Credentials must be passed in the environment variables:
// SONIC_USERID and SONIC_APIKEY.
// Credentials are created by calling the API with a username/password pair
// https://public-api.sonic.net/dyndns#requesting_an_api_key for the specific hostname
// NOTE: SONIC does not support `_` in DNS entries created via the API.
// To get around this issue, a manual CNAME needs to be created
// re-pointing _acme-challenge.hostname to acme-challenge.hostname
func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get(EnvAPIUserId, EnvAPIApiKey)
	if err != nil {
		return nil, fmt.Errorf("sonic: %w", err)
	}

	config := NewDefaultConfig()
	config.UserId = values[EnvAPIUserId]
	config.ApiKey = values[EnvAPIApiKey]

	return NewDNSProviderConfig(config)
}

// NewDNSProviderConfig return a DNSProvider instance configured for Sonic.
func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("sonic: the configuration of the DNS provider is nil")
	}

	client, err := NewClient(config.UserId, config.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("sonic: %w", err)
	}

	client.HTTPClient = config.HTTPClient

	return &DNSProvider{
		client: client,
		config: config,
	}, nil
}

// Present creates a TXT record using the specified parameters.
// NOTE: SONIC does not support `_` in DNS entries.
func (d *DNSProvider) Present(domainName, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domainName, keyAuth)

    fqdn = strings.Replace(fqdn, "_acme-challenge", "acme-challenge", 1)

    if strings.Index(fqdn, "_") > 0 {
		return fmt.Errorf("sonic: unable to create FQDNs with an _ like %s", fqdn)
    }

    // Sonic does not support trining . in hostname
    fqdn = dns01.UnFqdn(fqdn)

	err := d.client.CreateOrUpdateRecord(fqdn, value, d.config.TTL)
	if err != nil {
		return fmt.Errorf("sonic: unable to create record for %s: %w", fqdn, err)
	}
	return nil
}

// CleanUp removes the TXT records matching the specified parameters.
func (d *DNSProvider) CleanUp(domainName, token, keyAuth string) error {
	return nil
}

// Timeout returns the timeout and interval to use when checking for DNS propagation.
// Adjusting here to cope with spikes in propagation times.
func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}
