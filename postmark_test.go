package postmark

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"goji.io"
)

var (
	tMux    = goji.NewMux()  //nolint:gochecknoglobals // test infrastructure
	tServer *httptest.Server //nolint:gochecknoglobals // test infrastructure
	client  *Client          //nolint:gochecknoglobals // test infrastructure
)

func init() { //nolint:gochecknoinits // test infrastructure
	tServer = httptest.NewServer(tMux)

	transport := &http.Transport{
		Proxy: func(_ *http.Request) (*url.URL, error) {
			// Reroute...
			return url.Parse(tServer.URL)
		},
	}

	client = NewClient("", "")
	client.HTTPClient = &http.Client{Transport: transport}
	client.BaseURL = tServer.URL
}
