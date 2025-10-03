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

func (s *PostmarkTestSuite) TestDoRequestInvalidURL() {
	// Create a client with an invalid base URL to trigger request creation error
	client := NewClient("server-token", "account-token")
	client.BaseURL = "ht!tp://invalid url with spaces"

	var result map[string]string
	err := client.doRequest(context.Background(), http.MethodGet, "test", nil, &result, serverToken)

	s.Require().Error(err, "Should fail with invalid URL")
}

func (s *PostmarkTestSuite) TestDoRequestMarshalError() {
	// Use a channel as payload which cannot be marshaled to JSON
	invalidPayload := make(chan int)

	var result map[string]string
	err := s.client.doRequest(context.Background(), http.MethodPost, "test", invalidPayload, &result, serverToken)

	s.Require().Error(err, "Should fail when payload cannot be marshaled")
	s.Contains(err.Error(), "json")
}

func (s *PostmarkTestSuite) TestDoRequestNilDestination() {
	responseJSON := `{"message": "success"}`

	s.mux.Get("/nil-dest", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	// Pass nil as destination - should not error, just skip unmarshaling
	err := s.client.doRequest(context.Background(), http.MethodGet, "nil-dest", nil, nil, serverToken)

	s.Require().NoError(err, "Should handle nil destination gracefully")
}

func (s *PostmarkTestSuite) TestDoRequestHTTPErrorInvalidJSON() {
	// Test error response that has invalid JSON
	s.mux.Get("/error-invalid-json", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`not valid json`))
	})

	var result map[string]string
	err := s.client.doRequest(context.Background(), http.MethodGet, "error-invalid-json", nil, &result, serverToken)

	s.Require().Error(err)
	s.Contains(err.Error(), "request failed with status")
}

func (s *PostmarkTestSuite) TestAPIError_Error() {
	apiErr := APIError{
		ErrorCode: 401,
		Message:   "Unauthorized: Missing or incorrect API token",
	}

	errMsg := apiErr.Error()
	s.Equal("Unauthorized: Missing or incorrect API token", errMsg)
}

func (s *PostmarkTestSuite) TestNewClient() {
	client := NewClient("test-server-token", "test-account-token")

	s.Require().NotNil(client)
	s.Equal("test-server-token", client.ServerToken)
	s.Equal("test-account-token", client.AccountToken)
	s.Equal(postmarkURL, client.BaseURL)
	s.NotNil(client.HTTPClient)
}
