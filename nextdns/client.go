package nextdns

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

type Client struct {
	*http.Client
	BaseURL string
}
type transport struct {
	apiKey              string
	underlyingTransport http.RoundTripper
}

func (c *Client) MustGet(endpoint string) []byte {
	url := c.BaseURL + endpoint
	req, err := c.Get(url)
	if err != nil {
		log.Printf("Error getting %s: %v", url, err)
		return nil
	}

	if req.StatusCode != 200 {
		log.Printf("HTTP error getting %s: %v", url, req.Status)
		return nil
	}
	defer req.Body.Close()

	var content bytes.Buffer
	_, err = content.ReadFrom(req.Body)
	if err != nil {
		log.Printf("Error getting %s: %v", url, err)
		return nil
	}
	return content.Bytes()
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("x-api-key", t.apiKey)
	return t.underlyingTransport.RoundTrip(req)
}

func NewClient(apiKey string) *Client {
	return &Client{
		Client: &http.Client{
			Transport: &transport{
				apiKey:              apiKey,
				underlyingTransport: http.DefaultTransport,
			},
			Timeout: 10 * time.Second,
		},
		BaseURL: "https://api.nextdns.io/",
	}
}
