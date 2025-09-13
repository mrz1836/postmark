package postmark

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Helper function to get minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// validateTemplateInjectionPatterns checks for template injection patterns
func validateTemplateInjectionPatterns(t *testing.T, templateFields []string) {
	for _, field := range templateFields {
		if strings.Contains(strings.ToLower(field), "<script") {
			t.Logf("Script tag detected in template field: %s", field[:minInt(50, len(field))])
		}

		ssiPatterns := []string{"{{#system", "{{exec", "{{eval", "{{import", "{{require", "<%", "%>", "${"}
		for _, pattern := range ssiPatterns {
			if strings.Contains(strings.ToLower(field), pattern) {
				t.Logf("Potential template injection pattern: %s", pattern)
			}
		}

		if len(field) > 100000 {
			t.Logf("Very large template field: %d characters", len(field))
		}
	}
}

// createValidationResponse creates a mock validation response
func createValidationResponse(validation ValidateTemplateBody) ValidateTemplateResponse {
	return ValidateTemplateResponse{
		AllContentIsValid: true,
		HTMLBody: Validation{
			ContentIsValid:   true,
			ValidationErrors: []ValidationError{},
			RenderedContent:  validation.HTMLBody,
		},
		TextBody: Validation{
			ContentIsValid:   true,
			ValidationErrors: []ValidationError{},
			RenderedContent:  validation.TextBody,
		},
		Subject: Validation{
			ContentIsValid:   true,
			ValidationErrors: []ValidationError{},
			RenderedContent:  validation.Subject,
		},
		SuggestedTemplateModel: validation.TestRenderModel,
	}
}

// validateTemplateErrorLeakage checks if errors leak template content
func validateTemplateErrorLeakage(t *testing.T, err error, fields []string) {
	if err == nil {
		return
	}
	for _, field := range fields {
		if strings.Contains(err.Error(), field) && len(field) > 100 {
			t.Errorf("Error message contains large template content")
		}
	}
}

// validateQueryInjection checks for injection attempts in query parameters
func validateQueryInjection(t *testing.T, query string) {
	if strings.Contains(query, "\r") || strings.Contains(query, "\n") {
		t.Errorf("Header injection attempt in query parameters")
	}
	if strings.Contains(query, "<script") || strings.Contains(query, "javascript:") {
		t.Errorf("Script injection attempt in query parameters")
	}
}

// createTemplateValidationServer creates a test server for template validation
func createTemplateValidationServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
			return
		}

		var sentValidation ValidateTemplateBody
		if err := json.Unmarshal(body, &sentValidation); err != nil { //nolint:musttag // ValidateTemplateBody has JSON tags where needed
			t.Errorf("Failed to unmarshal validation JSON: %v", err)
			return
		}

		templateFields := []string{sentValidation.Subject, sentValidation.TextBody, sentValidation.HTMLBody}
		validateTemplateInjectionPatterns(t, templateFields)

		response := createValidationResponse(sentValidation)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		respJSON, _ := json.Marshal(response) //nolint:musttag // ValidateTemplateResponse has appropriate JSON tags
		_, _ = w.Write(respJSON)
	}))
}

// validateTemplateResponse validates template validation response
func validateTemplateResponse(t *testing.T, response ValidateTemplateResponse) {
	if len(response.HTMLBody.RenderedContent) > 1000000 {
		t.Errorf("Rendered HTML content is excessively large: %d chars", len(response.HTMLBody.RenderedContent))
	}
	if len(response.TextBody.RenderedContent) > 1000000 {
		t.Errorf("Rendered text content is excessively large: %d chars", len(response.TextBody.RenderedContent))
	}
}

// validateQueryParameters validates query parameters for templates
func validateQueryParameters(t *testing.T, parsedQuery url.Values, templateType, layoutTemplate string) {
	countParam := parsedQuery.Get("count")
	offsetParam := parsedQuery.Get("offset")
	if countParam == "" || offsetParam == "" {
		t.Errorf("Missing required parameters")
	}

	if templateType != "" && parsedQuery.Get("TemplateType") == "" {
		t.Errorf("TemplateType parameter missing when expected")
	}

	if layoutTemplate != "" && parsedQuery.Get("LayoutTemplate") == "" {
		t.Errorf("LayoutTemplate parameter missing when expected")
	}
}

// createTemplateQueryServer creates a test server for template query testing
func createTemplateQueryServer(t *testing.T, templateType, layoutTemplate string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.RawQuery
		validateQueryInjection(t, query)

		parsedQuery, err := url.ParseQuery(query)
		if err != nil {
			t.Errorf("Query parameters not properly encoded: %v", err)
			return
		}

		validateQueryParameters(t, parsedQuery, templateType, layoutTemplate)

		response := templatesResponse{
			TotalCount: 1,
			Templates: []TemplateInfo{
				{
					TemplateID:     123,
					Name:           "Test Template",
					Active:         true,
					TemplateType:   templateType,
					LayoutTemplate: layoutTemplate,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		respJSON, _ := json.Marshal(response) //nolint:musttag // templatesResponse has appropriate JSON tags
		_, _ = w.Write(respJSON)
	}))
}

// validateTemplateQueryResponse validates template query response
func validateTemplateQueryResponse(t *testing.T, templates []TemplateInfo, total int64) {
	if templates == nil {
		t.Errorf("Templates slice should not be nil")
	}
	if total < 0 {
		t.Errorf("Total count should not be negative: %d", total)
	}
}

// validateTemplatedEmail validates templated email for security issues
func validateTemplatedEmail(t *testing.T, sentEmail TemplatedEmail) {
	if sentEmail.TemplateID < 0 {
		t.Logf("Negative template ID: %d", sentEmail.TemplateID)
	}
	if sentEmail.TemplateID > 999999999 {
		t.Logf("Extremely large template ID: %d", sentEmail.TemplateID)
	}

	if strings.Contains(sentEmail.TemplateAlias, "\r") || strings.Contains(sentEmail.TemplateAlias, "\n") {
		t.Errorf("Header injection in template alias")
	}
	if strings.Contains(sentEmail.TemplateAlias, "../") || strings.Contains(sentEmail.TemplateAlias, "..\\") {
		t.Logf("Path traversal attempt in template alias: %s", sentEmail.TemplateAlias)
	}

	modelJSON, err := json.Marshal(sentEmail.TemplateModel)
	if err == nil && len(modelJSON) > 100000 {
		t.Logf("Large template model: %d bytes", len(modelJSON))
	}

	for key, value := range sentEmail.TemplateModel {
		if strValue, ok := value.(string); ok {
			if strings.Contains(strValue, "<script") || strings.Contains(strValue, "javascript:") {
				t.Logf("Potential script injection in model value for key: %s", key)
			}
		}
	}
}

// createTemplatedEmailServer creates a test server for templated email testing
func createTemplatedEmailServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
			return
		}

		var sentEmail TemplatedEmail
		if err := json.Unmarshal(body, &sentEmail); err != nil { //nolint:musttag // TemplatedEmail has appropriate JSON tags
			t.Errorf("Failed to unmarshal templated email JSON: %v", err)
			return
		}

		validateTemplatedEmail(t, sentEmail)

		response := EmailResponse{
			To:        sentEmail.To,
			MessageID: "test-message-id",
			ErrorCode: 0,
			Message:   "OK",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		respJSON, _ := json.Marshal(response) //nolint:musttag // EmailResponse has appropriate JSON tags
		_, _ = w.Write(respJSON)
	}))
}

// validateTemplatedEmailResponse validates templated email response
func validateTemplatedEmailResponse(t *testing.T, response EmailResponse) {
	if response.ErrorCode < 0 {
		t.Errorf("Error code should not be negative: %d", response.ErrorCode)
	}
	if len(response.MessageID) > 1000 {
		t.Errorf("MessageID is excessively long: %d characters", len(response.MessageID))
	}
}

// normalizeTemplateBatchSize normalizes batch size to prevent resource exhaustion
func normalizeTemplateBatchSize(batchSize int) int {
	if batchSize < 0 {
		return 0
	}
	if batchSize > 100 {
		return 100
	}
	return batchSize
}

// createTemplateBatchEmails creates a batch of templated test emails
func createTemplateBatchEmails(batchSize int, templateID int64, templateAlias string) []TemplatedEmail {
	var emails []TemplatedEmail
	for i := 0; i < batchSize; i++ {
		email := TemplatedEmail{
			TemplateID:    templateID,
			TemplateAlias: templateAlias,
			TemplateModel: map[string]interface{}{
				"index": i,
				"name":  fmt.Sprintf("User %d", i),
			},
			From: "sender@example.com",
			To:   fmt.Sprintf("user%d@example.com", i),
		}
		emails = append(emails, email)
	}
	return emails
}

// createTemplateBatchServer creates a test server for batch templated email testing
func createTemplateBatchServer(t *testing.T, expectedBatchSize int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
			return
		}

		var batchData map[string]interface{}
		if err := json.Unmarshal(body, &batchData); err != nil {
			t.Errorf("Failed to unmarshal batch JSON: %v", err)
			return
		}

		if _, ok := batchData["Messages"]; !ok && expectedBatchSize > 0 {
			t.Errorf("Expected Messages field in batch request")
		}

		var responses []EmailResponse
		for i := 0; i < expectedBatchSize; i++ {
			responses = append(responses, EmailResponse{
				To:        fmt.Sprintf("user%d@example.com", i),
				MessageID: fmt.Sprintf("batch-msg-%d", i),
				ErrorCode: 0,
				Message:   "OK",
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		respJSON, _ := json.Marshal(responses) //nolint:musttag // EmailResponse slice has appropriate JSON tags
		_, _ = w.Write(respJSON)
	}))
}

// validateTemplateBatchResponse validates batch templated email response
func validateTemplateBatchResponse(t *testing.T, responses []EmailResponse, expectedBatchSize int) {
	if len(responses) != expectedBatchSize {
		t.Errorf("Response count mismatch: expected %d, got %d", expectedBatchSize, len(responses))
	}

	for _, resp := range responses {
		if len(resp.MessageID) > 1000 {
			t.Errorf("MessageID too long: %d characters", len(resp.MessageID))
		}
		if resp.ErrorCode < 0 {
			t.Errorf("Negative error code: %d", resp.ErrorCode)
		}
	}
}

// FuzzTemplateValidation tests the ValidateTemplate function with various inputs
func FuzzTemplateValidation(f *testing.F) {
	// Seed corpus with various template body structures
	f.Add("Hello {{name}}", "Hello {{name}}", "<h1>Hello {{name}}</h1>", `{"name":"World"}`)
	f.Add("", "", "", `{}`)
	f.Add("Subject", "Text body", "<html><body>HTML body</body></html>", `{"var":"value"}`)
	f.Add("{{#each items}}{{this}}{{/each}}", "{{#each items}}{{this}}{{/each}}", "", `{"items":["a","b","c"]}`)
	f.Add("{{missing}}", "{{undefined}}", "{{invalid}}", `{"defined":"value"}`)
	f.Add("Very long template content that might cause issues", "", "", `{}`)
	f.Add("{{>partial}}", "{{>partial}}", "{{>partial}}", `{"data":"test"}`)
	f.Add("Unicode: test-fire", "Unicode: test-fire", "<p>Unicode: test-fire</p>", `{"unicode":"test-fire"}`)

	f.Fuzz(func(t *testing.T, subject, textBody, htmlBody, modelJSON string) {
		var testModel map[string]interface{}
		if err := json.Unmarshal([]byte(modelJSON), &testModel); err != nil {
			testModel = map[string]interface{}{}
		}

		validateBody := ValidateTemplateBody{
			Subject:                    subject,
			TextBody:                   textBody,
			HTMLBody:                   htmlBody,
			TestRenderModel:            testModel,
			InlineCSSForHTMLTestRender: true,
		}

		server := createTemplateValidationServer(t)
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		response, err := client.ValidateTemplate(context.Background(), validateBody)
		validateTemplateErrorLeakage(t, err, []string{subject, textBody, htmlBody})

		if err == nil {
			validateTemplateResponse(t, response)
		}
	})
}

// FuzzTemplateQueryParams tests query parameter encoding in GetTemplatesFiltered
func FuzzTemplateQueryParams(f *testing.F) {
	// Seed corpus with various template filter parameters
	f.Add(int64(10), int64(0), "Standard", "")
	f.Add(int64(50), int64(25), "Layout", "base-layout")
	f.Add(int64(100), int64(100), "", "")
	f.Add(int64(-1), int64(-1), "Invalid", "invalid-layout")
	f.Add(int64(0), int64(999999), "Standard", "very-long-layout-name-that-might-cause-issues")
	f.Add(int64(1000000), int64(0), "Layout\r\nInjected", "layout")
	f.Add(int64(25), int64(50), "test", "unicode-layout-test")

	f.Fuzz(func(t *testing.T, count, offset int64, templateType, layoutTemplate string) {
		server := createTemplateQueryServer(t, templateType, layoutTemplate)
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		templates, total, err := client.GetTemplatesFiltered(context.Background(), count, offset, templateType, layoutTemplate)

		if err != nil {
			if strings.Contains(err.Error(), templateType) && len(templateType) > 50 {
				t.Errorf("Error contains large template type parameter")
			}
			if strings.Contains(err.Error(), layoutTemplate) && len(layoutTemplate) > 50 {
				t.Errorf("Error contains large layout template parameter")
			}
		} else {
			validateTemplateQueryResponse(t, templates, total)
		}
	})
}

// FuzzTemplatedEmail tests templated email sending with various template models
func FuzzTemplatedEmail(f *testing.F) {
	// Seed corpus with different template model structures
	f.Add(int64(123), "", `{"name":"John","email":"john@example.com"}`)
	f.Add(int64(0), "welcome-template", `{"user":{"name":"Jane","preferences":{"email":true}}}`)
	f.Add(int64(456), "", `{}`)
	f.Add(int64(789), "complex-template", `{"items":[{"id":1,"name":"Item 1"},{"id":2,"name":"Item 2"}]}`)
	f.Add(int64(-1), "invalid-template", `{"invalid":"json"`)
	f.Add(int64(999), "", `{"large_data":"`+strings.Repeat("x", 1000)+`"}`)
	f.Add(int64(100), "unicode", `{"unicode":"test-fire","nested":{"deep":"value"}}`)

	f.Fuzz(func(t *testing.T, templateID int64, templateAlias, modelJSON string) {
		var templateModel map[string]interface{}
		if err := json.Unmarshal([]byte(modelJSON), &templateModel); err != nil {
			templateModel = map[string]interface{}{}
		}

		templatedEmail := TemplatedEmail{
			TemplateID:    templateID,
			TemplateAlias: templateAlias,
			TemplateModel: templateModel,
			From:          "sender@example.com",
			To:            "recipient@example.com",
			InlineCSS:     true,
		}

		server := createTemplatedEmailServer(t)
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		response, err := client.SendTemplatedEmail(context.Background(), templatedEmail)

		if err != nil {
			if strings.Contains(err.Error(), modelJSON) && len(modelJSON) > 100 {
				t.Errorf("Error message contains template model JSON")
			}
		} else {
			validateTemplatedEmailResponse(t, response)
		}
	})
}

// FuzzTemplateBatch tests batch templated email sending
func FuzzTemplateBatch(f *testing.F) {
	// Seed corpus with different batch configurations
	f.Add(int(1), int64(100), "template-alias")
	f.Add(int(5), int64(200), "")
	f.Add(int(50), int64(300), "batch-template")
	f.Add(int(0), int64(400), "empty-batch")
	f.Add(int(-1), int64(500), "negative")

	f.Fuzz(func(t *testing.T, batchSize int, templateID int64, templateAlias string) {
		normalizedSize := normalizeTemplateBatchSize(batchSize)
		emails := createTemplateBatchEmails(normalizedSize, templateID, templateAlias)

		server := createTemplateBatchServer(t, normalizedSize)
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		responses, err := client.SendTemplatedEmailBatch(context.Background(), emails)

		if err != nil {
			t.Logf("Batch templated email failed (may be expected): %v", err)
		} else {
			validateTemplateBatchResponse(t, responses, normalizedSize)
		}
	})
}
