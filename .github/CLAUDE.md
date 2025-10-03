# CLAUDE.md - Postmark Go Library

## 🎯 Quick Start

**Postmark Go Library** - Unofficial Go client for Postmark email API with comprehensive coverage and minimal dependencies.

**Key Info:**
- Go 1.18+, single dependency: `testify` for testing
- Context-aware API with Server/Account token auth
- 38 Go files covering complete Postmark API surface
- Test-driven with custom TestRouter for HTTP mocking

## 🏗️ Architecture & Patterns

### Client Structure
```go
type Client struct {
    HTTPClient   *http.Client
    ServerToken  string    // Server-level operations
    AccountToken string    // Account-level operations
    BaseURL      string    // API endpoint
}
```

### API Implementation Pattern
Every API endpoint follows this consistent pattern:
```go
func (client *Client) MethodName(ctx context.Context, payload PayloadType) (ResponseType, error) {
    res := ResponseType{}
    err := client.doRequest(ctx, parameters{
        Method:    "POST|GET|PUT|DELETE",
        Path:      "api/endpoint",
        Payload:   payload,          // nil for GET
        TokenType: serverToken,      // or accountToken
    }, &res)
    return res, err
}
```

### File Organization
- `postmark.go` - Core Client, doRequest, error handling
- `email.go` - Send emails, batch operations
- `templates.go` - Template management and templated emails
- `data_removals.go` - GDPR data removal requests
- `webhooks.go` - Webhook configuration
- `bounce.go`, `domains.go`, `message_streams.go`, etc. - Specialized APIs

## 🧪 Testing Patterns

### Test Structure (using testify suite)
```go
func (s *PostmarkTestSuite) TestMethodName() {
    tests := []struct {
        name         string
        responseJSON string
        wantErr      bool
        // expected fields
    }{
        {
            name: "successful operation",
            responseJSON: `{"field": "value"}`,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        s.Run(tt.name, func() {
            s.mux.Post("/endpoint", func(w http.ResponseWriter, _ *http.Request) {
                _, _ = w.Write([]byte(tt.responseJSON))
            })

            result, err := s.client.MethodName(context.Background(), payload)

            if tt.wantErr {
                s.Require().Error(err)
            } else {
                s.Require().NoError(err)
                // assertions
            }
        })
    }
}
```

### TestRouter Usage
- `s.mux.Get()`, `s.mux.Post()`, `s.mux.Put()`, `s.mux.Delete()`
- Supports path parameters: `/domains/:domainID`
- Custom TestRouter in `test_router.go`

## 🛠️ Common Tasks

### Adding New API Endpoint
1. **Create struct types** for request/response in appropriate file
2. **Add method** following the standard pattern
3. **Write tests** using TestRouter pattern
4. **Update README.md** API coverage checklist if needed

### JSON Struct Guidelines
- Use `json:"fieldName,omitempty"` for optional fields
- Match Postmark API field names exactly
- Use proper Go naming (e.g., `HTMLBody` for `HtmlBody`)
- Add struct tags for ID fields: `json:"TemplateID"`

### Error Handling
- API errors use `APIError` struct with ErrorCode and Message
- Context cancellation supported throughout
- Method-specific errors like `ErrEmailFailed`

## 🔧 Build & Test Commands

**Using MAGE-X:**
```bash
magex test          # Run tests (fast)
magex test:race     # Run with race detector
magex bench         # Run benchmarks
magex help          # View all commands
```

## 📁 Key Files

| File                   | Purpose                                     |
|------------------------|---------------------------------------------|
| `postmark.go`          | Core client, doRequest method, auth headers |
| `email.go`             | Email sending, batch operations             |
| `templates.go`         | Template CRUD, templated email sending      |
| `data_removals.go`     | GDPR data removal API                       |
| `webhooks.go`          | Webhook management                          |
| `test_router.go`       | Custom HTTP router for testing              |
| `examples/examples.go` | Usage examples for all major features       |

## ⚠️ Important Notes

- **Token Types**: Use `serverToken` (default) or `accountToken` based on API requirements
- **Context**: Always pass context.Context as first parameter
- **Testing**: Use `PostmarkTestSuite` pattern, not standalone tests
- **Dependencies**: Keep minimal - only add if absolutely necessary
- **API Coverage**: This library has near-complete Postmark API coverage
- **Conventions**: Follow existing patterns exactly for consistency

## 🚀 Examples

**Send Email:**
```go
client := postmark.NewClient("[SERVER-TOKEN]", "[ACCOUNT-TOKEN]")
email := postmark.Email{
    From: "no-reply@example.com", To: "user@example.com",
    Subject: "Test", HTMLBody: "<p>Hello</p>", TrackOpens: true,
}
res, err := client.SendEmail(context.Background(), email)
```

**Create Template:**
```go
template := postmark.Template{
    Name: "Welcome", Subject: "Welcome {{name}}!",
    HTMLBody: "<p>Hello {{name}}!</p>", Alias: "welcome-template",
}
res, err := client.CreateTemplate(context.Background(), template)
```
