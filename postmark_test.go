package postmark

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
	"goji.io"
)

type PostmarkTestSuite struct {
	suite.Suite

	mux    *goji.Mux
	server *httptest.Server
	client *Client
}

func (s *PostmarkTestSuite) SetupSuite() {
	s.mux = goji.NewMux()
	s.server = httptest.NewServer(s.mux)

	transport := &http.Transport{
		Proxy: func(_ *http.Request) (*url.URL, error) {
			return url.Parse(s.server.URL)
		},
	}

	s.client = NewClient("", "")
	s.client.HTTPClient = &http.Client{Transport: transport}
	s.client.BaseURL = s.server.URL
}

func (s *PostmarkTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
}

func TestPostmarkSuite(t *testing.T) {
	suite.Run(t, new(PostmarkTestSuite))
}
