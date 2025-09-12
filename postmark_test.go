package postmark

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
	"goji.io"
	"goji.io/pat"
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

	s.mux.HandleFunc(pat.Post("/test"), func(w http.ResponseWriter, req *http.Request) {
		s.Equal("application/json", req.Header.Get("Content-Type"))
		s.Equal("application/json", req.Header.Get("Accept"))
		s.Equal("server-token", req.Header.Get("X-Postmark-Server-Token"))
		_, _ = w.Write([]byte(responseJSON))
	})

	var result map[string]string
	err := s.client.doRequest(context.Background(), parameters{
		Method:    "POST",
		Path:      "test",
		TokenType: serverToken,
		Payload:   payload,
	}, &result)

	s.Require().NoError(err)
	s.Equal("success", result["message"])
}

func (s *PostmarkTestSuite) TestDoRequestWithAccountToken() {
	responseJSON := `{"message": "success"}`

	s.mux.HandleFunc(pat.Get("/account-test"), func(w http.ResponseWriter, req *http.Request) {
		s.Equal("account-token", req.Header.Get("X-Postmark-Account-Token"))
		s.Empty(req.Header.Get("X-Postmark-Server-Token"))
		_, _ = w.Write([]byte(responseJSON))
	})

	var result map[string]string
	err := s.client.doRequest(context.Background(), parameters{
		Method:    "GET",
		Path:      "account-test",
		TokenType: accountToken,
	}, &result)

	s.Require().NoError(err)
	s.Equal("success", result["message"])
}

func (s *PostmarkTestSuite) TestDoRequestHTTPError() {
	s.mux.HandleFunc(pat.Get("/error-test"), func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"ErrorCode": 500, "Message": "Internal Server Error"}`))
	})

	var result map[string]string
	err := s.client.doRequest(context.Background(), parameters{
		Method:    "GET",
		Path:      "error-test",
		TokenType: serverToken,
	}, &result)

	s.Require().Error(err)
}

func (s *PostmarkTestSuite) TestDoRequestContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var result map[string]string
	err := s.client.doRequest(ctx, parameters{
		Method:    "GET",
		Path:      "test",
		TokenType: serverToken,
	}, &result)

	s.Require().Error(err)
	s.Contains(err.Error(), "context canceled")
}

func (s *PostmarkTestSuite) TestDoRequestInvalidJSON() {
	s.mux.HandleFunc(pat.Get("/invalid-json"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`invalid json`))
	})

	var result map[string]string
	err := s.client.doRequest(context.Background(), parameters{
		Method:    "GET",
		Path:      "invalid-json",
		TokenType: serverToken,
	}, &result)

	s.Require().Error(err)
}
