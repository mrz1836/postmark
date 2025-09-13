package postmark

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// testEmail is a test-specific struct with explicit JSON tags for linter compliance
type testEmail struct {
	From        string       `json:"From,omitempty"`
	To          string       `json:"To,omitempty"`
	Cc          string       `json:"Cc,omitempty"`
	Bcc         string       `json:"Bcc,omitempty"`
	Subject     string       `json:"Subject,omitempty"`
	Tag         string       `json:"Tag,omitempty"`
	HTMLBody    string       `json:"HtmlBody,omitempty"`
	TextBody    string       `json:"TextBody,omitempty"`
	ReplyTo     string       `json:"ReplyTo,omitempty"`
	Headers     []Header     `json:"Headers,omitempty"`
	TrackOpens  bool         `json:"TrackOpens,omitempty"`
	TrackLinks  string       `json:"TrackLinks,omitempty"`
	Attachments []Attachment `json:"Attachments,omitempty"`
}

// testEmailResponse is a test-specific struct with explicit JSON tags for linter compliance
type testEmailResponse struct {
	To          string    `json:"To"`
	SubmittedAt time.Time `json:"SubmittedAt"`
	MessageID   string    `json:"MessageID"`
	ErrorCode   int       `json:"ErrorCode"`
	Message     string    `json:"Message"`
}

// convertToEmail converts testEmail to Email
func (te testEmail) convertToEmail() Email {
	return Email{
		From:        te.From,
		To:          te.To,
		Cc:          te.Cc,
		Bcc:         te.Bcc,
		Subject:     te.Subject,
		Tag:         te.Tag,
		HTMLBody:    te.HTMLBody,
		TextBody:    te.TextBody,
		ReplyTo:     te.ReplyTo,
		Headers:     te.Headers,
		TrackOpens:  te.TrackOpens,
		TrackLinks:  te.TrackLinks,
		Attachments: te.Attachments,
	}
}

// validateEmailSecurityIssues checks for common email security issues
func validateEmailSecurityIssues(t *testing.T, emailFields []string) {
	for _, field := range emailFields {
		if strings.Contains(field, "\r") || strings.Contains(field, "\n") {
			t.Logf("Potential header injection detected in email field: %s", field)
		}
		if strings.Contains(field, "\x00") {
			t.Logf("Null byte detected in email field")
		}
		if len(field) > 320 {
			t.Logf("Email field exceeds RFC limit: %d characters", len(field))
		}
	}
}

// validateErrorLeakage checks if errors contain sensitive information
func validateErrorLeakage(t *testing.T, err error, sensitiveData string, maxLen int) {
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), sensitiveData) && len(sensitiveData) > maxLen {
		t.Errorf("Error message contains large sensitive data")
	}
}

// createEmailTestServer creates a standard test server for email testing
func createEmailTestServer(t *testing.T, validateFunc func(*testing.T, Email)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
			return
		}

		var sentTestEmail testEmail
		if err := json.Unmarshal(body, &sentTestEmail); err != nil { //nolint:musttag // testEmail has explicit JSON tags
			t.Errorf("Failed to unmarshal email JSON: %v", err)
			return
		}

		if validateFunc != nil {
			sentEmail := sentTestEmail.convertToEmail()
			validateFunc(t, sentEmail)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"To":"test@example.com","SubmittedAt":"2023-01-01T00:00:00Z","MessageID":"123","ErrorCode":0,"Message":"OK"}`)
	}))
}

// validateEmailAddressFields validates email address fields for security issues
func validateEmailAddressFields(t *testing.T, sentEmail Email) {
	emailFields := []string{sentEmail.From, sentEmail.To, sentEmail.Cc, sentEmail.Bcc, sentEmail.ReplyTo}
	validateEmailSecurityIssues(t, emailFields)
}

// validateEmailHeaders validates email headers for security issues
func validateEmailHeaders(t *testing.T, sentEmail Email) {
	for _, header := range sentEmail.Headers {
		if strings.Contains(header.Name, "\r") || strings.Contains(header.Name, "\n") {
			t.Logf("Header injection in name: %s", header.Name)
		}
		if strings.Contains(header.Value, "\r") || strings.Contains(header.Value, "\n") {
			t.Logf("Header injection in value: %s", header.Value)
		}
		if strings.Contains(header.Name, "\x00") || strings.Contains(header.Value, "\x00") {
			t.Logf("Null byte in header")
		}
		if len(header.Name) > 1000 || len(header.Value) > 10000 {
			t.Logf("Large header detected: name=%d, value=%d chars", len(header.Name), len(header.Value))
		}
	}
}

// validateEmailAttachments validates email attachments for security issues
func validateEmailAttachments(t *testing.T, sentEmail Email) {
	for _, att := range sentEmail.Attachments {
		if strings.Contains(att.Name, "..") || strings.Contains(att.Name, "/") || strings.Contains(att.Name, "\\") {
			t.Logf("Potential path traversal in filename: %s", att.Name)
		}
		if strings.Contains(att.Name, "\x00") {
			t.Logf("Null byte in filename")
		}
		if att.Content != "" {
			if _, err := base64.StdEncoding.DecodeString(att.Content); err != nil {
				t.Logf("Invalid base64 content: %v", err)
			}
		}
		if len(att.Content) > 1000000 {
			t.Logf("Large attachment detected: %d base64 characters", len(att.Content))
		}
		if len(att.Name) > 255 {
			t.Logf("Long filename: %d characters", len(att.Name))
		}
	}
}

// normalizeBatchSize normalizes batch size to prevent resource exhaustion
func normalizeBatchSize(batchSize int) int {
	if batchSize < 0 {
		return 0
	}
	if batchSize > 200 {
		return 200
	}
	return batchSize
}

// createBatchEmails creates a batch of test emails
func createBatchEmails(batchSize int, emailDomain string) []Email {
	var emails []Email
	for i := 0; i < batchSize; i++ {
		emails = append(emails, Email{
			From:     "sender@example.com",
			To:       fmt.Sprintf("user%d@%s", i, emailDomain),
			Subject:  fmt.Sprintf("Test Email %d", i),
			TextBody: fmt.Sprintf("Test body %d", i),
		})
	}
	return emails
}

// createBatchTestServer creates a test server for batch email testing
func createBatchTestServer(t *testing.T, expectedBatchSize int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
			return
		}

		var sentTestEmails []testEmail
		if unmarshalErr := json.Unmarshal(body, &sentTestEmails); unmarshalErr != nil { //nolint:musttag // testEmail has explicit JSON tags
			t.Errorf("Failed to unmarshal batch email JSON: %v", unmarshalErr)
			return
		}

		if len(sentTestEmails) != expectedBatchSize {
			t.Errorf("Batch size mismatch: expected %d, got %d", expectedBatchSize, len(sentTestEmails))
		}

		var testResponses []testEmailResponse
		for i, email := range sentTestEmails {
			testResponses = append(testResponses, testEmailResponse{
				To:          email.To,
				SubmittedAt: time.Now(),
				MessageID:   fmt.Sprintf("msg-%d", i),
				ErrorCode:   0,
				Message:     "OK",
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		respJSON, err := json.Marshal(testResponses)
		if err != nil {
			t.Errorf("Failed to marshal response JSON: %v", err)
			return
		}
		_, _ = w.Write(respJSON)
	}))
}

// validateBatchResponse validates batch email response
func validateBatchResponse(t *testing.T, responses []EmailResponse, expectedCount int) {
	if len(responses) != expectedCount {
		t.Errorf("Response count mismatch: expected %d, got %d", expectedCount, len(responses))
	}

	for _, resp := range responses {
		if len(resp.MessageID) > 1000 {
			t.Errorf("MessageID too long: %d characters", len(resp.MessageID))
		}
	}
}

// FuzzEmailAddressValidation tests email address handling in various fields
func FuzzEmailAddressValidation(f *testing.F) {
	// Seed corpus with valid and invalid email formats
	f.Add("user@domain.com")
	f.Add("test@example.co.uk")
	f.Add("user+tag@domain.com")
	f.Add("user.name@domain.org")
	f.Add("user123@domain123.net")
	f.Add("invalid-email")
	f.Add("@domain.com")
	f.Add("user@")
	f.Add("user@@domain.com")
	f.Add("user@domain..com")
	f.Add("")
	f.Add("very-long-email-address-that-exceeds-normal-limits@very-long-domain-name-that-might-cause-issues.com")
	f.Add("unicode-test@domain.com")
	f.Add("user@domain.com, user2@domain.com")       // Multiple emails
	f.Add("user@domain.com\r\nBCC: evil@hacker.com") // Header injection
	f.Add("user@domain.com\x00null@byte.com")

	f.Fuzz(func(t *testing.T, emailAddr string) {
		email := Email{
			From:     emailAddr,
			To:       emailAddr,
			Cc:       emailAddr,
			Bcc:      emailAddr,
			ReplyTo:  emailAddr,
			Subject:  "Test Subject",
			TextBody: "Test body",
		}

		server := createEmailTestServer(t, validateEmailAddressFields)
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		response, err := client.SendEmail(context.Background(), email)
		validateErrorLeakage(t, err, emailAddr, 100)

		if err == nil {
			if response.ErrorCode != 0 {
				t.Logf("Email send reported error code: %d - %s", response.ErrorCode, response.Message)
			}
			if len(response.MessageID) > 1000 {
				t.Errorf("MessageID is excessively long: %d characters", len(response.MessageID))
			}
		}
	})
}

// FuzzEmailHeaders tests custom header handling
func FuzzEmailHeaders(f *testing.F) {
	// Seed corpus with various header combinations
	f.Add("X-Custom", "value")
	f.Add("Content-Type", "text/html")
	f.Add("Reply-To", "custom@domain.com")
	f.Add("X-Priority", "1")
	f.Add("List-Unsubscribe", "<mailto:unsub@domain.com>")
	f.Add("X-Long-Header-Name-That-Might-Cause-Issues", "value")
	f.Add("", "empty-name")
	f.Add("Header", "")
	f.Add("Header\r\nInjection", "value")
	f.Add("Header", "value\r\nX-Injected: header")
	f.Add("Unicode-Header", "test-fire")
	f.Add("Null\x00Byte", "value\x00null")

	f.Fuzz(func(t *testing.T, headerName, headerValue string) {
		headers := []Header{
			{Name: headerName, Value: headerValue},
		}

		email := Email{
			From:     "test@example.com",
			To:       "recipient@example.com",
			Subject:  "Test Subject",
			TextBody: "Test body",
			Headers:  headers,
		}

		server := createEmailTestServer(t, validateEmailHeaders)
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		_, err := client.SendEmail(context.Background(), email)
		validateErrorLeakage(t, err, headerValue, 100)
	})
}

// FuzzEmailAttachments tests attachment handling
func FuzzEmailAttachments(f *testing.F) {
	// Seed corpus with various attachment scenarios
	f.Add("document.pdf", "application/pdf", "dGVzdCBjb250ZW50") // "test content" in base64
	f.Add("image.jpg", "image/jpeg", "")
	f.Add("", "text/plain", "dGVzdA==")
	f.Add("file with spaces.txt", "text/plain", "invalid-base64")
	f.Add("very-long-filename-that-might-cause-issues-with-processing.txt", "text/plain", "dGVzdA==")
	f.Add("unicode-file.txt", "text/plain", "dGVzdA==")
	f.Add("file\x00null.txt", "application/octet-stream", "dGVzdA==")
	f.Add("../../../etc/passwd", "text/plain", "dGVzdA==")

	f.Fuzz(func(t *testing.T, fileName, contentType, content string) {
		attachment := Attachment{
			Name:        fileName,
			ContentType: contentType,
			Content:     content,
			ContentID:   "cid:test",
		}

		email := Email{
			From:        "test@example.com",
			To:          "recipient@example.com",
			Subject:     "Test with Attachment",
			TextBody:    "Test body",
			Attachments: []Attachment{attachment},
		}

		server := createEmailTestServer(t, validateEmailAttachments)
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		_, err := client.SendEmail(context.Background(), email)
		validateErrorLeakage(t, err, content, 100)
	})
}

// FuzzEmailBatch tests batch email sending
func FuzzEmailBatch(f *testing.F) {
	// Seed corpus with different batch sizes and configurations
	f.Add(int(1), "single@example.com")
	f.Add(int(5), "batch@example.com")
	f.Add(int(50), "large@example.com")      // Max batch size
	f.Add(int(100), "oversized@example.com") // Over limit
	f.Add(int(0), "empty@example.com")
	f.Add(int(-1), "negative@example.com")

	f.Fuzz(func(t *testing.T, batchSize int, emailDomain string) {
		normalizedSize := normalizeBatchSize(batchSize)
		emails := createBatchEmails(normalizedSize, emailDomain)

		server := createBatchTestServer(t, normalizedSize)
		defer server.Close()

		client := &Client{
			HTTPClient:   &http.Client{},
			ServerToken:  "test-token",
			AccountToken: "test-account-token",
			BaseURL:      server.URL,
		}

		responses, err := client.SendEmailBatch(context.Background(), emails)

		if err != nil {
			t.Logf("Batch send failed (may be expected for edge cases): %v", err)
		} else {
			validateBatchResponse(t, responses, len(emails))
		}
	})
}
