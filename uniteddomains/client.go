package uniteddomains

import (
	"crypto/tls"
	"net/http"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

// DefaultSchema is the value used for the URL in case
// no schema is explicitly defined
var DefaultSchema = "https"

// Client is a UnitedDomains client representation
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// NewClient returns a new UnitedDomains client
func NewClient(email string, password string, configTLS *tls.Config) (*Client, error) {

	httpClient := cleanhttp.DefaultClient()
	httpClient.Transport.(*http.Transport).TLSClientConfig = configTLS

	client := Client{
		BaseURL: "https://www.united-domains.de/",
		HTTP:    httpClient,
	}

	return &client, nil
}
