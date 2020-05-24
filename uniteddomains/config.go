package uniteddomains

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
)

// Config describes de configuration interface of this provider
type Config struct {
	Email         string
	Password      string
	InsecureHTTPS bool
	CACertificate string
}

// Client returns a new client for accessing UnitedDomains
func (c *Config) Client() (*Client, error) {
	tlsConfig := &tls.Config{}

	if c.CACertificate != "" {

		caCert, _, err := pathorcontents.Read(c.CACertificate)
		if err != nil {
			return nil, fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		tlsConfig.RootCAs = caCertPool
	}

	tlsConfig.InsecureSkipVerify = c.InsecureHTTPS

	client, err := NewClient(c.Email, c.Password, tlsConfig)

	if err != nil {
		return nil, fmt.Errorf("Error setting up UnitedDomains client: %s", err)
	}

	log.Print("UnitedDomains Client configured")

	return client, nil
}
