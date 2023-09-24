package acme

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"software.sslmate.com/src/go-pkcs12"
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

// expandACMEUser creates a new instance of an ACME user from set
// email_address and private_key_pem fields, and a registration
// if one exists.
func expandACMEUser(accountKeyPem string, emailAddress string) (*acmeUser, error) {
	key, err := privateKeyFromPEM([]byte(accountKeyPem))
	if err != nil {
		return nil, err
	}

	user := &acmeUser{
		key:   key,
		Email: emailAddress,
	}

	return user, nil
}

// expandACMEClient creates a connection to an ACME server from resource data,
// and also returns the user.
//
// If loadReg is supplied, the registration information is loaded in to the
// user's registration, if it exists - if the account cannot be resolved by the
// private key, then the appropriate error is returned.
func expandACMEClient(accountKeyPem string, emailAddress string, serverUrl string, keyType string, loadReg bool) (*lego.Client, *acmeUser, error) {
	user, err := expandACMEUser(accountKeyPem, emailAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting user data: %s", err.Error())
	}

	config := lego.NewConfig(user)
	config.CADirURL = serverUrl

	// Note this function is used by both the registration and certificate
	// resources, but key type is not necessary during registration, so
	// it's okay if it's empty for that.
	if keyType != "" {
		config.Certificate.KeyType = certcrypto.KeyType(keyType)
	}

	var client *lego.Client
	newClient := func() error {
		client, err = lego.NewClient(config)
		if err != nil {
			if isAbleToRetry(err.Error()) {
				return err
			} else {
				return backoff.Permanent(err)
			}
		}

		// Populate user's registration resource if needed
		if loadReg {
			user.Registration, err = client.Registration.ResolveAccountByKey()
			if err != nil {
				if isAbleToRetry(err.Error()) {
					return err
				} else {
					return backoff.Permanent(err)
				}
			}
		}
		return nil
	}
	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = 30 * time.Minute
	err = backoff.Retry(newClient, reconnectBackoff)
	if err != nil {
		return nil, nil, err
	}

	return client, user, nil
}

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
	for _, ic := range cb[1:] {
		issuer = append(issuer, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ic.Raw})...)
	}

	return
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
