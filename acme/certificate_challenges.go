package acme

import (
	"fmt"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/hashicorp/go-multierror"
	"github.com/myklst/terraform-provider-acme/acme/dnsplugin"
)

type dnsBlock struct {
	Config   map[string]interface{} `json:"config"`
	Provider string                 `json:"provider"`
}

// setCertificateChallengeProviders sets all of the challenge providers in the
// client that are needed for obtaining the certificate.
//
// The returned func() is a closer for all of the configured DNS providers that
// should be called when they are no longer needed (i.e. in a defer after one of
// the CRUD functions are complete).
func setCertificateChallengeProviders(client *lego.Client, providers []dnsBlock) (func(), error) {
	// DNS
	dnsClosers := make([]func(), 0)
	dnsCloser := func() {
		for _, f := range dnsClosers {
			f()
		}
	}

	dnsProvider, err := NewDNSProviderWrapper()
	if err != nil {
		return dnsCloser, fmt.Errorf("%s, %w", "failed to create DNS provider wrapper", err)
	}

	for _, providerRaw := range providers {
		if p, closer, err := expandDNSChallenge(providerRaw); err == nil {
			dnsProvider.providers = append(dnsProvider.providers, p)
			dnsClosers = append(dnsClosers, closer)
		} else {
			return dnsCloser, fmt.Errorf("%s, %w", "failed to expand DNS challenge", err)
		}
	}

	if err := client.Challenge.SetDNS01Provider(dnsProvider, []dns01.ChallengeOption{}...); err != nil {
		return dnsCloser, fmt.Errorf("%s, %w", "failed to set DNS 01 provider", err)
	}

	return dnsCloser, nil
}

func expandDNSChallenge(m dnsBlock) (challenge.ProviderTimeout, func(), error) {
	var providerName string

	if m.Provider != "" {
		providerName = m.Provider
	} else {
		return nil, nil, fmt.Errorf("DNS challenge provider not defined")
	}

	// Config only needs to be set if it's defined, otherwise existing env/SDK
	// defaults are fine.
	config := make(map[string]string)
	for k, v := range m.Config {
		config[k] = v.(string)
	}

	return dnsplugin.NewClient(providerName, config, []string{""})
}

// DNSProviderWrapper is a multi-provider wrapper to support multiple
// DNS challenges.
type DNSProviderWrapper struct {
	providers []challenge.ProviderTimeout
}

// NewDNSProviderWrapper returns an freshly initialized
// DNSProviderWrapper.
func NewDNSProviderWrapper() (*DNSProviderWrapper, error) {
	return &DNSProviderWrapper{}, nil
}

// Present implements challenge.Provider for DNSProviderWrapper.
func (d *DNSProviderWrapper) Present(domain, token, keyAuth string) error {
	var err error
	for _, p := range d.providers {
		err = p.Present(domain, token, keyAuth)
		if err != nil {
			err = multierror.Append(err, fmt.Errorf("error encountered while presenting token for DNS challenge: %s", err.Error()))
		}
	}

	return err
}

// CleanUp implements challenge.Provider for DNSProviderWrapper.
func (d *DNSProviderWrapper) CleanUp(domain, token, keyAuth string) error {
	var err error
	for _, p := range d.providers {
		err = p.CleanUp(domain, token, keyAuth)
		if err != nil {
			err = multierror.Append(err, fmt.Errorf("error encountered while cleaning token for DNS challenge: %s", err.Error()))
		}
	}

	return err
}

// Timeout implements challenge.ProviderTimeout for
// DNSProviderWrapper.
//
// The highest polling interval and timeout values defined across all
// providers is used.
func (d *DNSProviderWrapper) Timeout() (time.Duration, time.Duration) {
	var timeout, interval time.Duration
	for _, p := range d.providers {
		t, i := p.Timeout()
		if t > timeout {
			timeout = t
		}

		if i > interval {
			interval = i
		}
	}

	if timeout < 1 {
		timeout = dns01.DefaultPropagationTimeout
	}

	if interval < 1 {
		interval = dns01.DefaultPollingInterval
	}

	return timeout, interval
}
