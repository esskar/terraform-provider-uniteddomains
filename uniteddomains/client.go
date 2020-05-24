package uniteddomains

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

// DefaultSchema is the value used for the URL in case
// no schema is explicitly defined
var DefaultSchema = "https"

// Client is a UnitedDomains client representation
type Client struct {
	BaseURL  string
	Csrf     string
	HTTP     *http.Client
	LoggedIn bool
}

// NewClient returns a new UnitedDomains client
func NewClient(email string, password string, configTLS *tls.Config) (*Client, error) {

	httpClient := cleanhttp.DefaultClient()
	httpClient.Transport.(*http.Transport).TLSClientConfig = configTLS

	client := Client{
		BaseURL:  "https://www.united-domains.de",
		Csrf:     "",
		HTTP:     httpClient,
		LoggedIn: false,
	}

	return &client, nil
}

// Creates a new request with necessary headers
func (client *Client) newRequest(method string, endpoint string, body []byte) (*http.Request, error) {

	var err error
	if !client.LoggedIn {
		client.LoggedIn, err = client.login()
	}

	if err != nil {
		return nil, err
	}

	var urlStr = client.BaseURL + endpoint
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("Error during parsing request URL: %s", err)
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("Error during creation of request: %s", err)
	}

	//req.Header.Add("X-API-Key", client.APIKey)
	//req.Header.Add("Accept", "application/json")

	if method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func (client *Client) login() (bool, error) {

	httpClient := client.HTTP

	url, err := url.Parse(client.BaseURL + "/login")
	if err != nil {
		return false, fmt.Errorf("Error while trying to login, request URL: %s", err)
	}

	if client.Csrf != "" {

	} else {

		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			return false, fmt.Errorf("Error during creation of request: %s", err)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			return false, err
		}
		data, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		if err != nil {
			return false, err
		}
		log.Printf("%s", data)
		if resp.StatusCode == 200 {
			return true, nil
		}
	}

	return false, nil
}
