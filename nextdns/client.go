package nextdns

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

// https://nextdns.github.io/api/
type Client struct {
	*http.Client
	BaseURL string
}
type transport struct {
	apiKey              string
	sem                 chan struct{} // semaphore
	underlyingTransport http.RoundTripper
}

func (c *Client) MustGet(endpoint string) []byte {
	// log.Println("starting MustGet for endpoint", endpoint)
	// defer log.Println("finished MustGet for endpoint", endpoint)
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
	t.sem <- struct{}{}
	req.Header.Add("x-api-key", t.apiKey)
	defer func() { <-t.sem }()
	return t.underlyingTransport.RoundTrip(req)
}

type Option func(*Client)

func WithMaxConcurrentRequests(n int) Option {
	return func(c *Client) {
		close(c.Client.Transport.(*transport).sem)
		c.Client.Transport.(*transport).sem = make(chan struct{}, n)
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.Client.Timeout = d
	}
}

func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		Client: &http.Client{
			Transport: &transport{
				apiKey:              apiKey,
				sem:                 make(chan struct{}, 3),
				underlyingTransport: http.DefaultTransport,
			},
			Timeout: 10 * time.Second,
		},
		BaseURL: "https://api.nextdns.io/",
	}
	for _, o := range opts {
		o(c)
	}
	return c
}
