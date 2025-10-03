package postmark

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// FuzzDoRequestJSONUnmarshal tests the JSON unmarshaling logic in doRequest
// with various malformed and edge case JSON inputs
func FuzzDoRequestJSONUnmarshal(f *testing.F) {
	// Seed corpus with valid and edge case JSON responses
	f.Add(`{"ErrorCode":0,"Message":"Success"}`)
	f.Add(`{"ErrorCode":422,"Message":"Invalid request"}`)
	f.Add(`{}`)
	f.Add(`{"ErrorCode":null,"Message":""}`)
	f.Add(`{"ErrorCode":"invalid","Message":123}`)
	f.Add(`{"Extra":"field","ErrorCode":0}`)
	f.Add(`null`)
	f.Add(`""`)
	f.Add(`[]`)
	f.Add(`{"nested":{"deep":{"value":true}}}`)

	f.Fuzz(func(t *testing.T, jsonResponse string) {
		// Create a test server that returns the fuzzed JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, jsonResponse)
		}))
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		var result map[string]interface{}

		// Test that the function doesn't panic with malformed JSON
		err := client.doRequest(context.Background(), http.MethodGet, "test", nil, &result, serverToken)
		// We expect either success or an error, but never a panic
		// The function should handle malformed JSON gracefully
		if err != nil {
			// Ensure error messages don't contain the raw input to prevent injection
			if strings.Contains(err.Error(), jsonResponse) && len(jsonResponse) > 100 {
				t.Errorf("Error message should not contain raw input for large payloads")
			}
		}
	})
}

// FuzzAPIErrorHandling tests error response parsing with various error structures
func FuzzAPIErrorHandling(f *testing.F) {
	// Seed corpus with different error response formats
	f.Add(int64(422), `{"ErrorCode":422,"Message":"The request was invalid"}`)
	f.Add(int64(401), `{"ErrorCode":401,"Message":"Unauthorized"}`)
	f.Add(int64(500), `{"ErrorCode":500,"Message":"Internal server error"}`)
	f.Add(int64(400), `{"ErrorCode":"invalid","Message":"Bad request"}`)
	f.Add(int64(404), `{"Message":"Not found"}`)
	f.Add(int64(503), `{}`)
	f.Add(int64(429), `{"ErrorCode":null,"Message":null}`)
	f.Add(int64(400), `"plain text error"`)
	f.Add(int64(500), `<html><body>Error</body></html>`)

	f.Fuzz(func(t *testing.T, statusCode int64, errorResponse string) {
		// Ensure status code is in valid HTTP error range
		if statusCode < 400 || statusCode > 599 {
			statusCode = 400 + (statusCode % 200) // Normalize to 400-599 range
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(int(statusCode))
			_, _ = io.WriteString(w, errorResponse)
		}))
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		var result map[string]interface{}

		err := client.doRequest(context.Background(), http.MethodGet, "test", nil, &result, serverToken)

		// For error status codes, we should always get an error
		if err == nil {
			t.Errorf("Expected error for HTTP status %d, but got nil", statusCode)
		}

		// If it's an APIError, validate its structure
		var apiErr APIError
		if errors.As(err, &apiErr) {
			// Error code should be reasonable (not negative, not extremely large)
			if apiErr.ErrorCode < 0 || apiErr.ErrorCode > 999999 {
				t.Logf("APIError has unusual ErrorCode: %d", apiErr.ErrorCode)
			}

			// Message should not be excessively long to prevent DoS
			if len(apiErr.Message) > 10000 {
				t.Errorf("Error message is excessively long: %d characters", len(apiErr.Message))
			}
		}
	})
}

// FuzzJSONPayloadMarshaling tests the JSON marshaling of request payloads
func FuzzJSONPayloadMarshaling(f *testing.F) {
	// Seed corpus with different payload types
	f.Add(`{"string":"value","number":42,"boolean":true}`)
	f.Add(`{"nested":{"deep":{"array":[1,2,3]}}}`)
	f.Add(`{"empty":"","null":null,"zero":0}`)
	f.Add(`{"unicode":"test-fire","special":"tab\there"}`)
	f.Add(`{"large_number":9223372036854775807}`)

	f.Fuzz(func(t *testing.T, payloadJSON string) {
		// Try to unmarshal the fuzzed input into a generic map
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
			// Skip invalid JSON as we're testing the marshaling, not parsing
			return
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify the request has proper headers
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
			}

			// Read and validate the request body can be parsed
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Failed to read request body: %v", err)
				return
			}

			var requestPayload map[string]interface{}
			if err := json.Unmarshal(body, &requestPayload); err != nil {
				t.Errorf("Request body is not valid JSON: %v", err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{"success":true}`)
		}))
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		var result map[string]interface{}

		// Test that marshaling doesn't panic with complex payloads
		err := client.doRequest(context.Background(), http.MethodPost, "test", payload, &result, serverToken)
		if err != nil {
			t.Logf("Request failed (may be expected): %v", err)
		}
	})
}
