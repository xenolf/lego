package checkdomain

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-acme/lego/v3/challenge/dns01"
	"github.com/go-acme/lego/v3/platform/config/env"
)

const (
	envEndpoint           = "CHECKDOMAIN_ENDPOINT"
	envToken              = "CHECKDOMAIN_TOKEN"
	envTTL                = "CHECKDOMAIN_TTL"
	envHTTPTimeout        = "CHECKDOMAIN_HTTP_TIMEOUT"
	envPropagationTimeout = "CHECKDOMAIN_PROPAGATION_TIMEOUT"
	envPollingInterval    = "CHECKDOMAIN_POLLING_INTERVAL"
)

const (
	defaultEndpoint string = "https://api.checkdomain.de"
	defaultTTL      int    = 300
)

// Config is used to configure the creation of the DNSProvider
type Config struct {
	Endpoint           *url.URL
	Token              string
	TTL                int
	Timeout            time.Duration
	PropagationTimeout time.Duration
	PollingInterval    time.Duration
}

// NewDefaultConfig returns a default configuration for the DNSProvider
func NewDefaultConfig() *Config {
	endpoint, _ := url.Parse(env.GetOrDefaultString(envEndpoint, defaultEndpoint))

	return &Config{
		Endpoint:           endpoint,
		Token:              env.GetOrDefaultString(envToken, ""),
		TTL:                env.GetOrDefaultInt(envTTL, defaultTTL),
		PropagationTimeout: env.GetOrDefaultSecond(envPropagationTimeout, 3*time.Minute),
		PollingInterval:    env.GetOrDefaultSecond(envPollingInterval, 5*time.Second),
		Timeout:            env.GetOrDefaultSecond(envHTTPTimeout, 30*time.Second),
	}
}

// DNSProvider implements challenge.Provider for the checkdomain API
// specified at https://developer.checkdomain.de/reference/.
type DNSProvider struct {
	config          *Config
	httpClient      *http.Client
	domainIDMapping map[string]int
}

func NewDNSProvider() (*DNSProvider, error) {
	config := NewDefaultConfig()

	if config.Endpoint == nil {
		return nil, fmt.Errorf("checkdomain: some information are invalid: CHECKDOMAIN_ENDPOINT")
	}

	if config.Token == "" {
		return nil, fmt.Errorf("checkdomain: some information are missing: CHECKDOMAIN_TOKEN")
	}

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config.Endpoint == nil {
		return nil, fmt.Errorf("checkdomain: invalid endpoint")
	}

	if config.Token == "" {
		return nil, fmt.Errorf("checkdomain: missing token")
	}

	return &DNSProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		domainIDMapping: make(map[string]int),
	}, nil
}

func (p *DNSProvider) present(present bool, domain, keyAuth string) error {
	domainID, err := p.getDomainIDByName(domain)
	if err != nil {
		return fmt.Errorf("checkdomain: %v", err)
	}

	if domainID == -1 {
		return fmt.Errorf("checkdomain: domain not found")
	}

	using, err := p.isUsingCheckdomainNameservers(domainID)
	if err != nil {
		return fmt.Errorf("checkdomain: %v", err)
	}

	if !using {
		return fmt.Errorf("checkdomain: not using checkdomain nameservers, can not update records")
	}

	name, value := dns01.GetRecord(domain, keyAuth)
	if present {
		err = p.createRecord(domainID, &Record{
			Name:  name,
			TTL:   p.config.TTL,
			Type:  "TXT",
			Value: value,
		})
	} else {
		// absent
		err = p.deleteRecord(domainID, "TXT", name, value)
	}

	if err != nil {
		return fmt.Errorf("checkdomain: %v", err)
	}

	return nil
}

// Present creates a TXT record to fulfill the dns-01 challenge
func (p *DNSProvider) Present(domain, token, keyAuth string) error {
	return p.present(true, domain, keyAuth)
}

// CleanUp removes the TXT record previously created
func (p *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	return p.present(false, domain, keyAuth)
}

func (p *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return p.config.PropagationTimeout, p.config.PollingInterval
}
