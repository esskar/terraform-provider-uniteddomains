package uniteddomains

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

// DefaultSchema is the value used for the URL in case
// no schema is explicitly defined
var DefaultSchema = "https"

// Client is a UnitedDomains client representation
type Client struct {
	BaseURL    string
	Csrf       string
	CsrfMeta   string
	CsrfScript string
	Email      string
	HTTP       *http.Client
	LoggedIn   bool
	Password   string
	SessionId  string
}

// NewClient returns a new UnitedDomains client
func NewClient(email string, password string, configTLS *tls.Config) (*Client, error) {

	uri, err := url.Parse("https://www.united-domains.de")
	if err != nil {
		return nil, fmt.Errorf("Error during parsing request URL: %s", err)
	}

	httpClient := cleanhttp.DefaultClient()
	httpClient.Transport.(*http.Transport).TLSClientConfig = configTLS

	/*
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			return nil, err
		}
		var cookies []*http.Cookie
		cookies = append(cookies, &http.Cookie{
			Name:  "CookieSettingsGroupId",
			Value: "2464190"})
		jar.SetCookies(uri, cookies)
		httpClient.Jar = jar*/

	client := Client{
		BaseURL:    uri.String(),
		Csrf:       "",
		CsrfMeta:   "",
		CsrfScript: "",
		Email:      email,
		HTTP:       httpClient,
		LoggedIn:   false,
		Password:   password,
		SessionId:  "",
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

func (client *Client) prepareRequest(req *http.Request, contentType string, acceptType string) *http.Request {
	if acceptType != "" {
		req.Header.Set("Accept", acceptType)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Origin", client.BaseURL)
	req.AddCookie(&http.Cookie{
		Name:  "CookieSettingsGroupId",
		Value: "2464190",
	})
	if client.SessionId != "" {
		req.AddCookie(&http.Cookie{
			Name:  "SESSID",
			Value: client.SessionId,
		})
	}
	return req
}

func (client *Client) prepareXhr(req *http.Request, contentType string, acceptType string) *http.Request {
	req = client.prepareRequest(req, contentType, acceptType)

	req.Header.Set("http-x-csrf-token", client.CsrfMeta)
	req.Header.Set("x-csrf-token", client.CsrfScript)
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	return req
}

func (client *Client) login() (bool, error) {

	loginUrl, err := url.Parse(client.BaseURL + "/login")
	if err != nil {
		return false, fmt.Errorf("Error while trying to login, request URL: %s", err)
	}

	httpClient := client.HTTP

	if client.Csrf != "" {
		values := url.Values{
			"csrf":     {client.Csrf},
			"selector": {"login"},
			"email":    {client.Email},
			"pwd":      {client.Password},
			"submit":   {"Login"},
		}

		req, err := http.NewRequest("POST", loginUrl.String(), strings.NewReader(values.Encode()))
		if err != nil {
			return false, fmt.Errorf("Error during creation of request: %s", err)
		}
		req = client.prepareRequest(req, "application/x-www-form-urlencoded", "text/html")

		dump, _ := httputil.DumpRequestOut(req, false)
		log.Print(string(dump))

		res, err := httpClient.Do(req)
		if err != nil {
			return false, err
		}
		defer res.Body.Close()

		dump, _ = httputil.DumpResponse(res, true)
		log.Print(string(dump))

		log.Printf("Failed to login, POST: %d %s %s", res.StatusCode, res.Status, loginUrl.String())

		if res.StatusCode != 302 {
			return false, fmt.Errorf("Failed to login, POST: %d %s %s", res.StatusCode, res.Status, loginUrl.String())
		}

		log.Print("Successfully logged in.")

		return true, nil
	} else {

		setUserLanguageUrl, err := url.Parse(client.BaseURL + "/set-user-language")
		if err != nil {
			return false, fmt.Errorf("Error while trying to login, failed to parse URL: %s", err)
		}

		//syncSessionUrl, err := url.Parse(client.BaseURL + "/sync-session/request")
		//if err != nil {
		//	return false, fmt.Errorf("Error while trying to login, failed to parse URL: %s", err)
		//}

		err = client.initLogin(loginUrl, httpClient)
		if err != nil {
			return false, err
		}

		err = client.setUserLanguage(setUserLanguageUrl, httpClient)
		if err != nil {
			return false, err
		}

		//err = client.syncSession(syncSessionUrl, httpClient)
		//if err != nil {
		//	return false, err
		//}

		return client.login()
	}
}

func (client *Client) initLogin(loginUrl *url.URL, httpClient *http.Client) error {
	err := client.fetchCsrfAndSession(loginUrl, httpClient)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) fetchCsrfAndSession(loginUrl *url.URL, httpClient *http.Client) error {
	req, err := http.NewRequest("GET", loginUrl.String(), nil)
	if err != nil {
		return fmt.Errorf("Error during creation of request: %s", err)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to login, GET: %d %s %s", res.StatusCode, res.Status, loginUrl.String())
	}

	for _, cookie := range res.Cookies() {
		if cookie.Name == "SESSID" {
			client.SessionId = strings.TrimSpace(cookie.Value)
			break
		}
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	s := doc.Find("#login-form-1 input[name=csrf]").First()
	csrf, _ := s.Attr("value")
	if csrf == "" {
		return fmt.Errorf("Failed to get CSRF token.")
	}

	s = doc.Find("meta[name=csrf]").First()
	csrfMeta, _ := s.Attr("content")

	csrfScript, err := goquery.OuterHtml(doc.First())
	if err != nil {
		return err
	}

	pattern := "\"CSRF_TOKEN\":\""
	ndx := strings.Index(csrfScript, pattern)
	if ndx > -1 {
		csrfScript = csrfScript[ndx+len(pattern):]
		ndx = strings.Index(csrfScript, "\"")
		csrfScript = csrfScript[:ndx]
	} else {
		csrfScript = ""
	}

	log.Printf("Successfully fetched csrf tokens: %s, %s, %s", csrf, csrfMeta, csrfScript)

	client.Csrf = csrf
	client.CsrfMeta = csrfMeta
	client.CsrfScript = csrfScript

	return nil
}

func (client *Client) setUserLanguage(setUserLanguageUrl *url.URL, httpClient *http.Client) error {
	values := url.Values{
		"language": {"en-US"},
	}
	req, err := http.NewRequest("POST", setUserLanguageUrl.String(), strings.NewReader(values.Encode()))
	if err != nil {
		return fmt.Errorf("Error during creation of request: %s", err)
	}
	req = client.prepareXhr(req, "application/x-www-form-urlencoded; charset=UTF-8", "*/*")

	dump, _ := httputil.DumpRequestOut(req, false)
	log.Print(string(dump))

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	dump, _ = httputil.DumpResponse(res, false)
	log.Print(string(dump))

	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to set user language, POST: %d %s %s", res.StatusCode, res.Status, setUserLanguageUrl.String())
	}

	log.Print("Successfully set user language")

	return nil
}

func (client *Client) syncSession(syncSessionUrl *url.URL, httpClient *http.Client) error {
	values := url.Values{
		"method":   {"GET"},
		"uri":      {"/login/"},
		"referrer": {""},
	}
	req, err := http.NewRequest("POST", syncSessionUrl.String(), strings.NewReader(values.Encode()))
	if err != nil {
		return fmt.Errorf("Error during creation of request: %s", err)
	}
	req = client.prepareXhr(req, "application/x-www-form-urlencoded; charset=UTF-8", "*/*")

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to sync session, POST: %d %s %s", res.StatusCode, res.Status, syncSessionUrl.String())
	}

	log.Print("Successfully synced session")

	return nil
}
