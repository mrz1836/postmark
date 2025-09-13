package postmark

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// validateJSONError checks if an error is a valid JSON parsing error
func validateJSONError(t *testing.T, err error) {
	if !strings.Contains(err.Error(), "invalid character") &&
		!strings.Contains(err.Error(), "unexpected end") &&
		!strings.Contains(err.Error(), "cannot unmarshal") &&
		!strings.Contains(err.Error(), "json") &&
		!strings.Contains(err.Error(), "unmarshal") {
		t.Logf("Unexpected error type: %v", err)
	}
}

// validateInjectionPatterns checks for common injection patterns in query strings
func validateInjectionPatterns(t *testing.T, query string) {
	injectionPatterns := []string{
		"<script", "javascript:", "data:", "vbscript:",
		"onload=", "onerror=", "onclick=",
		"../", "..\\", "/etc/", "\\windows\\",
		"' OR ", "\" OR ", "; DROP ", "UNION SELECT",
		"\r\n", "\n\r", "\x00", "\x1f",
	}

	for _, pattern := range injectionPatterns {
		if strings.Contains(strings.ToLower(query), strings.ToLower(pattern)) {
			t.Logf("Potential injection pattern detected in query: %s", pattern)
		}
	}
}

// validateURLEncoding checks if URL parameters are properly encoded
func validateURLEncoding(t *testing.T, query, paramName, paramValue string) {
	parsedURL, err := url.ParseQuery(query)
	if err != nil {
		t.Errorf("Query parameters not properly encoded: %v", err)
		return
	}

	if val := parsedURL.Get(paramName); val != paramValue {
		t.Logf("Parameter encoding: %s -> %s", paramValue, val)
	}
}

// createBounceTestServer creates a test server for bounce testing
func createBounceTestServer(t *testing.T, paramName, paramValue string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.RawQuery
		validateInjectionPatterns(t, query)
		validateURLEncoding(t, query, paramName, paramValue)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"TotalCount":0,"Bounces":[]}`)
	}))
}

// validateErrorResponse checks if error responses contain sensitive information
func validateErrorResponse(t *testing.T, err error, paramValue string) {
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), paramValue) && len(paramValue) > 50 {
		t.Errorf("Error message should not contain large parameter values")
	}
}

// validateBouncesData validates bounce data for integrity
func validateBouncesData(t *testing.T, bounces []Bounce) {
	for _, bounce := range bounces {
		if bounce.ID < 0 {
			t.Errorf("Bounce ID should not be negative: %d", bounce.ID)
		}
		if len(bounce.Email) > 320 {
			t.Errorf("Email address is too long: %d characters", len(bounce.Email))
		}
	}
}

// validateBounceResponse validates successful bounce response
func validateBounceResponse(t *testing.T, bounces []Bounce, total int64) {
	if bounces == nil {
		t.Logf("Bounces is nil - may indicate empty response")
	}
	if total < 0 {
		t.Errorf("Total count should not be negative: %d", total)
	}
	validateBouncesData(t, bounces)
}

// FuzzGetBouncedTagsJSON tests the custom JSON parsing logic in GetBouncedTags
// This function has special handling for Postmark's unusual array response format
func FuzzGetBouncedTagsJSON(f *testing.F) {
	// Seed corpus with valid and edge case JSON array responses
	f.Add(`["tag1","tag2","tag3"]`)
	f.Add(`[]`)
	f.Add(`["single-tag"]`)
	f.Add(`["tag with spaces","tag-with-dashes","tag_with_underscores"]`)
	f.Add(`["unicode-test","emoji-fire","special-\t\n\r"]`)
	f.Add(`["very-long-tag-name-that-exceeds-normal-expectations-and-might-cause-issues"]`)
	f.Add(`[null]`)
	f.Add(`[123,"string",true]`)
	f.Add(`["","",""]`)
	f.Add(`["tag1",null,"tag3"]`)

	f.Fuzz(func(t *testing.T, jsonArray string) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, jsonArray)
		}))
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		// Test the GetBouncedTags function which has custom JSON handling
		tags, err := client.GetBouncedTags(context.Background())

		// The function should handle malformed JSON gracefully
		if err != nil {
			validateJSONError(t, err)
		} else {
			// If successful, validate the result
			for _, tag := range tags {
				// Tags should not be excessively long (potential DoS protection)
				if len(tag) > 1000 {
					t.Errorf("Tag is excessively long: %d characters", len(tag))
				}
			}

			// Result should not have excessive number of tags
			if len(tags) > 10000 {
				t.Errorf("Excessive number of tags returned: %d", len(tags))
			}
		}
	})
}

// FuzzBounceQueryParams tests URL parameter encoding in GetBounces
func FuzzBounceQueryParams(f *testing.F) {
	// Seed corpus with various parameter combinations
	f.Add(int64(10), int64(0), "type", "HardBounce")
	f.Add(int64(50), int64(100), "tag", "newsletter")
	f.Add(int64(1), int64(999999), "email", "test@example.com")
	f.Add(int64(25), int64(50), "subject", "Test Subject")
	f.Add(int64(0), int64(0), "messageID", "abc123")
	f.Add(int64(-5), int64(-10), "fromdate", "2023-01-01")
	f.Add(int64(1000000), int64(1000000), "todate", "2023-12-31")

	f.Fuzz(func(t *testing.T, count, offset int64, optKey, optValue string) {
		// Create options map with fuzzed values
		options := make(map[string]interface{})
		if optKey != "" {
			options[optKey] = optValue
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate query parameters are properly encoded
			query := r.URL.RawQuery

			// Check for potential injection attempts
			if strings.Contains(query, "<script") || strings.Contains(query, "javascript:") {
				t.Errorf("Potential XSS in query parameters: %s", query)
			}

			// Verify URL parameters are properly escaped
			parsedURL, err := url.ParseQuery(query)
			if err != nil {
				t.Errorf("Query parameters are not properly encoded: %v", err)
				return
			}

			// Verify count and offset are present and numeric
			countParam := parsedURL.Get("count")
			offsetParam := parsedURL.Get("offset")

			if countParam == "" || offsetParam == "" {
				t.Errorf("Missing required parameters: count=%s, offset=%s", countParam, offsetParam)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{"TotalCount":0,"Bounces":[]}`)
		}))
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		// Test GetBounces with fuzzed parameters
		bounces, total, err := client.GetBounces(context.Background(), count, offset, options)

		if err != nil {
			t.Logf("Request failed (may be expected): %v", err)
		} else {
			// Validate response structure
			if bounces == nil {
				t.Errorf("Bounces slice should not be nil")
			}
			if total < 0 {
				t.Errorf("Total count should not be negative: %d", total)
			}
		}
	})
}

// FuzzBounceQueryParamsInjection tests for potential injection vulnerabilities
// in query parameter construction
func FuzzBounceQueryParamsInjection(f *testing.F) {
	// Seed corpus with potential injection attempts
	f.Add("type", "test&malicious=true")
	f.Add("tag", "test%26malicious%3Dtrue")
	f.Add("email", "test@evil.com&admin=true")
	f.Add("subject", "test\r\nX-Injected: header")
	f.Add("messageID", "../../../etc/passwd")
	f.Add("custom", "<script>alert('xss')</script>")
	f.Add("param", "' OR '1'='1")
	f.Add("value", "test\x00null\x00byte")

	f.Fuzz(func(t *testing.T, paramName, paramValue string) {
		// Skip empty parameter names
		if paramName == "" {
			return
		}

		options := map[string]interface{}{
			paramName: paramValue,
		}

		server := createBounceTestServer(t, paramName, paramValue)
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		_, _, err := client.GetBounces(context.Background(), 10, 0, options)
		validateErrorResponse(t, err, paramValue)
	})
}

// FuzzBounceJSONStructure tests bounce data structure parsing
func FuzzBounceJSONStructure(f *testing.F) {
	// Seed corpus with different bounce response structures
	f.Add(`{"TotalCount":1,"Bounces":[{"ID":123,"Type":"HardBounce","Email":"test@example.com"}]}`)
	f.Add(`{"TotalCount":0,"Bounces":[]}`)
	f.Add(`{"TotalCount":"invalid","Bounces":null}`)
	f.Add(`{"Bounces":[{"ID":"string","Type":123,"Email":true}]}`)
	f.Add(`{"TotalCount":999999999999999,"Bounces":[{}]}`)
	f.Add(`{}`)
	f.Add(`null`)

	f.Fuzz(func(t *testing.T, bounceJSON string) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, bounceJSON)
		}))
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		bounces, total, err := client.GetBounces(context.Background(), 10, 0, nil)

		if err != nil {
			validateJSONError(t, err)
		} else {
			validateBounceResponse(t, bounces, total)
		}
	})
}
