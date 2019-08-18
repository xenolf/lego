package liquidweb

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/go-acme/lego/v3/challenge/dns01"
	"github.com/go-acme/lego/v3/platform/config/env"
	lwclient "github.com/liquidweb/liquidweb-go/client"
	"github.com/liquidweb/liquidweb-go/network"
)

// Config is used to configure the creation of the DNSProvider
type Config struct {
	URL                string
	Username           string
	Password           string
	Zone               string
	Timeout            time.Duration
	PropagationTimeout time.Duration
}

// NewDefaultConfig returns a default configuration for the DNSProvider
func NewDefaultConfig() *Config {
	config := &Config{
		Timeout:            env.GetOrDefaultSecond("LIQUID_WEB_TIMEOUT", 1*time.Minute),
		PropagationTimeout: env.GetOrDefaultSecond("LIQUID_WEB_PROPAGATION_TIMEOUT", 2*time.Minute),
	}

	return config
}

// DNSProvider is an implementation of the acme.ChallengeProvider interface
// that uses Liquid Web's REST API to manage TXT records for a domain.
type DNSProvider struct {
	config      *Config
	recordIDs   map[string]int
	recordIDsMu sync.Mutex
	client      *lwclient.API
}

// NewDNSProvider returns a DNSProvider instance configured for Liquid Web.
func NewDNSProvider() (*DNSProvider, error) {
	config := NewDefaultConfig()
	config.URL = env.GetOrDefaultString("LIQUID_WEB_URL", "")
	config.Username = env.GetOrDefaultString("LIQUID_WEB_USERNAME", "")
	config.Password = env.GetOrDefaultString("LIQUID_WEB_PASSWORD", "")
	config.Zone = env.GetOrDefaultString("LIQUID_WEB_ZONE", "")

	return NewDNSProviderConfig(config)
}

// NewDNSProviderConfig return a DNSProvider instance configured for Liquid Web.
func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the DNS provider is nil")
	}

	if config.URL == "" {
		return nil, fmt.Errorf("liquidweb: url is missing")
	}

	if config.Username == "" {
		return nil, fmt.Errorf("liquidweb: username is missing")
	}

	if config.Password == "" {
		return nil, fmt.Errorf("liquidweb: password is missing")
	}

	// Initial new LW go client.
	lwAPI, err := lwclient.NewAPI(config.Username, config.Password, config.URL, int(config.Timeout.Seconds()))

	if err != nil {
		log.Fatalf("Could not create Liquid Web API client: %v", err)
	}

	return &DNSProvider{
		config:    config,
		recordIDs: make(map[string]int),
		client:    lwAPI,
	}, nil
}

// Timeout returns the timeout and interval to use when checking for DNS propagation.
// Adjusting here to cope with spikes in propagation times.
func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.Timeout, d.config.PropagationTimeout
}

// Present creates a TXT record using the specified parameters
func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	params := &network.DNSRecordParams{
		Name:  fqdn[0 : len(fqdn)-1],
		RData: strconv.Quote(value),
		Type:  "TXT",
		Zone:  d.config.Zone,
	}

	dnsEntry, err := d.client.NetworkDNS.Create(params)
	if err != nil {
		return fmt.Errorf("could not create TXT record: %v", err)
	}

	d.recordIDsMu.Lock()
	d.recordIDs[token] = int(dnsEntry.ID)
	d.recordIDsMu.Unlock()

	return nil
}

// CleanUp removes the TXT record matching the specified parameters
func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	// get the record's unique ID from when we created it
	d.recordIDsMu.Lock()
	recordID, ok := d.recordIDs[token]
	d.recordIDsMu.Unlock()
	if !ok {
		return fmt.Errorf("unknown record ID for '%s'", domain)
	}
	fmt.Printf("%+v", d.recordIDs)

	params := &network.DNSRecordParams{ID: recordID}
	_, err := d.client.NetworkDNS.Delete(params)
	if err != nil {
		return fmt.Errorf("could not remove TXT record: %v", err)
	}
	// Delete record ID from map
	d.recordIDsMu.Lock()
	delete(d.recordIDs, domain)
	d.recordIDsMu.Unlock()

	return nil
}
