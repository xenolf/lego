package alicloud

import (
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/dns"
	"github.com/xenolf/lego/acme"
)

// DNSProvider is an implementation of the acme.ChallengeProvider interface.
type DNSProvider struct {
	client *dns.Client
}

// NewDNSProvider returns a DNSProvider instance configured for AliCloud client.
// Credentials must be passed in the environment variables: ALICLOUD_API_KEY and ALICLOUD_API_SECRET.
func NewDNSProvider() (*DNSProvider, error) {
	apiKey := os.Getenv("ALICLOUD_API_KEY")
	apiSecret := os.Getenv("ALICLOUD_API_SECRET")
	if apiKey == "" {
		return nil, fmt.Errorf("AliCloud credentials missing")
	}
	if apiSecret == "" {
		return nil, fmt.Errorf("AliCloud credentials Secret missing")
	}
	return NewDNSProviderCredentials(apiKey, apiSecret)
}

// NewDNSProviderCredentials uses the supplied credentials to return a
// DNSProvider instance configured for AliCloud DNSProvider.
func NewDNSProviderCredentials(apiKey, apiSecret string) (*DNSProvider, error) {
	return &DNSProvider{
		client: dns.NewClient(apiKey, apiSecret),
	}, nil
}

// Present creates a record with a secret
func (provider *DNSProvider) Present(domain, token, keyAuth string) error {
	fqdn, value, _ := acme.DNS01Record(domain, keyAuth)
	zoneDomain, domainErr := provider.getHostedZone(domain)

	if domainErr != nil {
		return domainErr
	}

	name := provider.extractRecordName(fqdn, zoneDomain)

	_, err := provider.client.AddDomainRecord(&dns.AddDomainRecordArgs{
		DomainName: zoneDomain,
		RR:         name,
		Type:       "TXT",
		Value:      value,
		TTL:        600,
	})

	if err != nil {
		return fmt.Errorf("AliCloud call failed: %v", err)
	}
	return nil
}

// CleanUp removes a given record that was generated by Present
func (provider *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	fqdn, _, _ := acme.DNS01Record(domain, keyAuth)
	recordId, recordIdErr := provider.findTxtRecords(domain, fqdn)
	if recordIdErr != nil {
		return recordIdErr
	}

	_, err := provider.client.DeleteDomainRecord(&dns.DeleteDomainRecordArgs{
		RecordId: recordId,
	})

	return err
}

// Timeout returns the values (20*time.Minute, 20*time.Second) which
// are used by the acme package as timeout and check interval values
// when checking for DNS record propagation with AliCloud.
func (provider *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return 20 * time.Minute, 20 * time.Second
}

func (provider *DNSProvider) getHostedZone(domain string) (string, error) {
	var allDomains []dns.DomainType
	pagination := common.Pagination{
		PageNumber: 1,
		PageSize:   100,
	}

	args := &dns.DescribeDomainsArgs{}

	for {
		args.Pagination = pagination
		domains, err := provider.client.DescribeDomains(args)
		if err != nil {
			return "", err
		}
		allDomains = append(allDomains, domains...)
		if len(domains) < pagination.PageSize {
			break
		}
		pagination.PageNumber += 1
	}

	var hostedDomain dns.DomainType
	for _, d := range allDomains {
		if strings.HasSuffix(domain, d.DomainName) {
			if len(d.DomainName) > len(hostedDomain.DomainName) {
				hostedDomain = d
			}
		}
	}

	if hostedDomain.DomainName == "" {
		return "", fmt.Errorf("No matching AliCloud domain found for domain %s", domain)
	}
	return hostedDomain.DomainName, nil
}

func (provider *DNSProvider) findTxtRecords(domain, fqdn string) (string, error) {
	var recordId string
	var allRecords []dns.RecordType

	zoneDomain, err := provider.getHostedZone(domain)
	if err != nil {
		return "", err
	}
	recordName := provider.extractRecordName(fqdn, zoneDomain)

	pagination := common.Pagination{
		PageNumber: 1,
		PageSize:   100,
	}

	args := &dns.DescribeDomainRecordsArgs{
		DomainName: zoneDomain,
	}

	for {
		args.Pagination = pagination
		response, err := provider.client.DescribeDomainRecords(args)
		if err != nil {
			return "", err
		}
		allRecords = append(allRecords, response.DomainRecords.Record...)
		if len(response.DomainRecords.Record) < pagination.PageSize {
			break
		}
		pagination.PageNumber += 1
	}

	for _, record := range allRecords {
		if record.Type == "TXT" && record.RR == recordName && record.DomainName == zoneDomain {
			recordId = record.RecordId
		}
	}

	if recordId == "" {
		return "", fmt.Errorf("No matching AliCloud domain record found for domain %s  and record %s ", domain, recordName)
	}

	return recordId, nil
}

func (provider *DNSProvider) extractRecordName(fqdn, domain string) string {
	name := acme.UnFqdn(fqdn)
	if idx := strings.Index(name, "."+domain); idx != -1 {
		return name[:idx]
	}
	return name
}
