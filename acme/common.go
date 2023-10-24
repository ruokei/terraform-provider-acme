package acme

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/registration"
	"github.com/myklst/terraform-provider-acme/v2/acme/dnsplugin"
	"software.sslmate.com/src/go-pkcs12"
)

const (
	DefaultMaxElapsedTime = 120 * time.Minute
)

// acmeUser implements acme.User.
type acmeUser struct {

	// The email address for the account.
	Email string

	// The registration resource object.
	Registration *registration.Resource

	// The private key for the account.
	key crypto.PrivateKey
}

func (u acmeUser) GetEmail() string {
	return u.Email
}
func (u acmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u acmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func trimStringQuotes(input string) string {
	return strings.TrimPrefix(strings.TrimSuffix(input, "\""), "\"")
}

func expandDNSChallenge(ctx context.Context, dns *dnsChallenge, nameServers []string) (challenge.ProviderTimeout, func(), error) {
	var providerName string

	if !dns.Provider.IsUnknown() && !dns.Provider.IsNull() {
		providerName = dns.Provider.ValueString()
	} else {
		return nil, nil, fmt.Errorf("DNS challenge provider not defined")
	}

	// Config only needs to be set if it's defined, otherwise existing env/SDK
	// defaults are fine.

	config := make(map[string]string)
	if !dns.Config.IsUnknown() && !dns.Config.IsNull() {
		dns.Config.ElementsAs(ctx, &config, false)
	}

	return dnsplugin.NewClient(providerName, config, nameServers)
}

// certSecondsRemaining takes an certificate.Resource, parses the
// certificate, and computes the seconds that it has remaining.
func certSecondsRemaining(cert *certificate.Resource) (int64, error) {
	x509Certs, err := parsePEMBundle(cert.Certificate)
	if err != nil {
		return 0, err
	}
	c := x509Certs[0]

	if c.IsCA {
		return 0, fmt.Errorf("first certificate is a CA certificate")
	}

	expiry := c.NotAfter.Unix()
	now := time.Now().Unix()

	return (expiry - now), nil
}

// certDaysRemaining takes an certificate.Resource, parses the
// certificate, and computes the days that it has remaining.
func certDaysRemaining(cert *certificate.Resource) (int64, error) {
	remaining, err := certSecondsRemaining(cert)
	if err != nil {
		return 0, fmt.Errorf("unable to calculate time to certificate expiry: %s", err)
	}

	return remaining / 86400, nil
}

// splitPEMBundle gets a slice of x509 certificates from
// parsePEMBundle.
//
// The first certificate split is returned as the issued certificate,
// with the rest returned as the issuer (intermediate) chain.
//
// Technically, it will be possible for issuer to be empty, if there
// are zero certificates in the intermediate chain. This is highly
// unlikely, however.
func splitPEMBundle(bundle []byte) (cert []byte, certNotAfter string, issuer []byte, err error) {
	cb, err := parsePEMBundle(bundle)
	if err != nil {
		return
	}

	// lego always returns the issued cert first, if the CA is first there is a problem
	if cb[0].IsCA {
		err = fmt.Errorf("first certificate is a CA certificate")
		return
	}

	cert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cb[0].Raw})
	certNotAfter = cb[0].NotAfter.Format(time.RFC3339)
	issuer = make([]byte, 0)
	for i, ic := range cb[1:] {
		issuer = append(issuer, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ic.Raw})...)
		if i < len(cb)-2 {
			issuer = append(issuer, '\n')
		}
	}

	return
}

// bundleToPKCS12 packs an issued certificate (and any supplied
// intermediates) into a PFX file.  The private key is included in
// the archive if it is a non-zero value.
//
// The returned archive is base64-encoded.
func bundleToPKCS12(bundle, key []byte, password string) ([]byte, error) {
	cb, err := parsePEMBundle(bundle)
	if err != nil {
		return nil, err
	}

	// lego always returns the issued cert first, if the CA is first there is a problem
	if cb[0].IsCA {
		return nil, fmt.Errorf("first certificate is a CA certificate")
	}

	pk, err := privateKeyFromPEM(key)
	if err != nil {
		return nil, err
	}

	pfxData, err := pkcs12.Encode(rand.Reader, pk, cb[0], cb[1:], password)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(pfxData)))
	base64.StdEncoding.Encode(buf, pfxData)
	return buf, nil
}

// parsePEMBundle parses a certificate bundle from top to bottom and returns
// a slice of x509 certificates. This function will error if no certificates are found.
//
// TODO: This was taken from lego directly, consider exporting it there, or
// consolidating with other TF crypto functions.
func parsePEMBundle(bundle []byte) ([]*x509.Certificate, error) {
	var certificates []*x509.Certificate
	var certDERBlock *pem.Block

	for {
		certDERBlock, bundle = pem.Decode(bundle)
		if certDERBlock == nil {
			break
		}

		if certDERBlock.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(certDERBlock.Bytes)
			if err != nil {
				return nil, err
			}
			certificates = append(certificates, cert)
		}
	}

	if len(certificates) == 0 {
		return nil, errors.New("no certificates were found while parsing the bundle")
	}

	return certificates, nil
}

// privateKeyFromPEM converts a PEM block into a crypto.PrivateKey.
func privateKeyFromPEM(pemData []byte) (crypto.PrivateKey, error) {
	var result *pem.Block
	rest := pemData
	for {
		result, rest = pem.Decode(rest)
		if result == nil {
			return nil, fmt.Errorf("cannot decode supplied PEM data")
		}
		switch result.Type {
		case "RSA PRIVATE KEY":
			return x509.ParsePKCS1PrivateKey(result.Bytes)
		case "EC PRIVATE KEY":
			return x509.ParseECPrivateKey(result.Bytes)
		}
	}
}

// csrFromPEM converts a PEM block into an *x509.CertificateRequest.
func csrFromPEM(pemData []byte) (*x509.CertificateRequest, error) {
	var result *pem.Block
	rest := pemData
	for {
		result, rest = pem.Decode(rest)
		if result == nil {
			return nil, fmt.Errorf("cannot decode supplied PEM data")
		}
		if result.Type == "CERTIFICATE REQUEST" {
			return x509.ParseCertificateRequest(result.Bytes)
		}
	}
}
