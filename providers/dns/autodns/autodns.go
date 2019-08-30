package autodns

import (
	"fmt"
	"github.com/go-acme/lego/v3/challenge/dns01"
	"github.com/go-acme/lego/v3/platform/config/env"
	"net/http"
	"net/url"
	"path"
)

const (
	defaultEndpoint string = `https://api.autodns.com/v1/`
	demoEndpoint    string = `https://api.demo.autodns.com/v1/`

	defaultEndpointContext int = 4
	demoEndpointContext    int = 1
)

type Config struct {
	Endpoint   *url.URL
	Username   string `json:"username"`
	Password   string `json:"password"`
	Context    int    `json:"-"`
	HTTPClient *http.Client
}

func NewDefaultConfig() *Config {
	endpoint, _ := url.Parse(defaultEndpoint)

	return &Config{
		Endpoint:   endpoint,
		Context:    defaultEndpointContext,
		HTTPClient: &http.Client{},
	}
}

type DNSProvider struct {
	config          *Config
	zoneNameservers map[string]string
	currentRecords  []*ResourceRecord
}

func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get("AUTODNS_API_USER", "AUTODNS_API_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("autodns: %v", err)
	}

	rawEndpoint := env.GetOrDefaultString("AUTODNS_ENDPOINT", defaultEndpoint)
	endpoint, err := url.Parse(rawEndpoint)
	if err != nil {
		return nil, fmt.Errorf("autodns: %v", err)
	}

	config := NewDefaultConfig()
	config.Username = values["AUTODNS_API_USER"]
	config.Password = values["AUTODNS_API_PASSWORD"]
	config.Endpoint = endpoint

	provider := &DNSProvider{config: config}

	// Because autodns needs the nameservers for each request, we query them all here and put them
	// in our state to avoid making a lot of requests later.
	req, err := provider.makeRequest(http.MethodPost, path.Join("zone", "_search"), nil)
	if err != nil {
		return nil, fmt.Errorf("autodns: %v", err)
	}

	var resp *DataZoneResponse
	if err := provider.sendRequest(req, &resp); err != nil {
		return nil, fmt.Errorf("autodns: %v", err)
	}

	provider.zoneNameservers = make(map[string]string, len(resp.Data))

	for _, zone := range resp.Data {
		provider.zoneNameservers[zone.Name] = zone.VirtualNameServer
	}

	return provider, nil
}

// Present creates a TXT record to fulfill the dns-01 challenge
func (d *DNSProvider) Present(domain, token, keyAuth string) error {

	// Get all current records for this domain to be able to restore them later and not clear everything
	// since the api does not support adding/removing single txt records.
	resp, err := d.getRecords(domain)
	if err != nil {
		return fmt.Errorf("autodns: getRecords: %v", err)
	}

	if len(resp.Data) > 0 {
		d.currentRecords = resp.Data[0].ResourceRecords
	}

	// Add the actual record
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	_, err = d.addTxtRecord(domain, fqdn, value)
	if err != nil {
		return fmt.Errorf("autodns: %v", err)
	}
	return nil
}

// CleanUp removes the TXT record previously created
func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	if err := d.restoreRecords(domain, "_acme-challenge"); err != nil {
		return fmt.Errorf("autodns: restoreRecords: %v", err)
	}

	return nil
}
