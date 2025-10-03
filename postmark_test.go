package postmark

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PostmarkTestSuite struct {
	suite.Suite

	mux    *TestRouter
	server *httptest.Server
	client *Client
}

func (s *PostmarkTestSuite) SetupSuite() {
	s.mux = NewTestRouter()
	s.server = httptest.NewServer(s.mux)

	transport := &http.Transport{
		Proxy: func(_ *http.Request) (*url.URL, error) {
			return url.Parse(s.server.URL)
		},
	}

	s.client = NewClient("server-token", "account-token")
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

func (s *PostmarkTestSuite) TestDoRequestWithPayload() {
	responseJSON := `{"message": "success"}`
	payload := map[string]string{"test": "data"}

	s.mux.Post("/test", func(w http.ResponseWriter, req *http.Request) {
		s.Equal("application/json", req.Header.Get("Content-Type"))
		s.Equal("application/json", req.Header.Get("Accept"))
		s.Equal("server-token", req.Header.Get("X-Postmark-Server-Token"))
		_, _ = w.Write([]byte(responseJSON))
	})

	var result map[string]string
	err := s.client.doRequest(context.Background(), http.MethodPost, "test", payload, &result, serverToken)

	s.Require().NoError(err)
	s.Equal("success", result["message"])
}

func (s *PostmarkTestSuite) TestDoRequestWithAccountToken() {
	responseJSON := `{"message": "success"}`

	s.mux.Get("/account-test", func(w http.ResponseWriter, req *http.Request) {
		s.Equal("account-token", req.Header.Get("X-Postmark-Account-Token"))
		s.Empty(req.Header.Get("X-Postmark-Server-Token"))
		_, _ = w.Write([]byte(responseJSON))
	})

	var result map[string]string
	err := s.client.doRequest(context.Background(), http.MethodGet, "account-test", nil, &result, accountToken)

	s.Require().NoError(err)
	s.Equal("success", result["message"])
}

func (s *PostmarkTestSuite) TestDoRequestHTTPError() {
	s.mux.Get("/error-test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"ErrorCode": 500, "Message": "Internal Server Error"}`))
	})

	var result map[string]string
	err := s.client.doRequest(context.Background(), http.MethodGet, "error-test", nil, &result, serverToken)

	s.Require().Error(err)
}

func (s *PostmarkTestSuite) TestDoRequestContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var result map[string]string
	err := s.client.doRequest(ctx, http.MethodGet, "test", nil, &result, serverToken)

	s.Require().Error(err)
	s.Contains(err.Error(), "context canceled")
}

func (s *PostmarkTestSuite) TestDoRequestInvalidJSON() {
	s.mux.Get("/invalid-json", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`invalid json`))
	})

	var result map[string]string
	err := s.client.doRequest(context.Background(), http.MethodGet, "invalid-json", nil, &result, serverToken)

	s.Require().Error(err)
}
