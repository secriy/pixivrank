package pixivrank

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

// NewClient return the http client instance.
func NewClient() *http.Client {
	cookiejar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client := &http.Client{Jar: cookiejar}
	return client
}

// NewClientWithProxy return the http client instance with a proxy server.
func NewClientWithPorxy(proxyUrl string) *http.Client {
	proxy, _ := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := NewClient()
	client.Transport = tr
	return client
}
