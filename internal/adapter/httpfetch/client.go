package httpfetch

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/janaka/web-analyzer/internal/service/analyzer"
)

type Client struct {
	http *http.Client
}

func New(dialTimeout, tlsTimeout, reqTimeout int) *Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout: time.Duration(dialTimeout) * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   time.Duration(tlsTimeout) * time.Second,
		ResponseHeaderTimeout: time.Duration(reqTimeout) * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
	}
	return &Client{
		http: &http.Client{
			Timeout:   time.Duration(reqTimeout) * time.Second,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "WebAnalyzerBot/1.0 (+https://example.com)")
	return c.http.Do(req)
}

func (c *Client) Head(ctx context.Context, url string) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	req.Header.Set("User-Agent", "WebAnalyzerBot/1.0 (+https://example.com)")
	return c.http.Do(req)
}

var _ analyzer.Fetcher = (*Client)(nil)
