package postmark

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Email is exactly what it sounds like
type Email struct {
	// From: REQUIRED The sender email address. Must have a registered and confirmed Sender Signature.
	From string `json:",omitempty"`
	// To: REQUIRED Recipient email address. Multiple addresses are comma separated. Max 50.
	To string `json:",omitempty"`
	// Cc recipient email address. Multiple addresses are comma separated. Max 50.
	Cc string `json:",omitempty"`
	// Bcc recipient email address. Multiple addresses are comma separated. Max 50.
	Bcc string `json:",omitempty"`
	// Subject: Email subject
	Subject string `json:",omitempty"`
	// Tag: Email tag that allows you to categorize outgoing emails and get detailed statistics.
	Tag string `json:",omitempty"`
	// HTMLBody: HTML email message. REQUIRED, If no TextBody specified
	HTMLBody string `json:"HtmlBody,omitempty"`
	// TextBody: Plain text email message. REQUIRED, If no HTMLBody specified
	TextBody string `json:",omitempty"`
	// ReplyTo: Reply To override email address. Defaults to the Reply To set in the sender signature.
	ReplyTo string `json:",omitempty"`
	// Headers: List of custom headers to include.
	Headers []Header `json:",omitempty"`
	// TrackOpens: Activate open tracking for this email.
	TrackOpens bool `json:",omitempty"`
	// TrackLinks:Activate link tracking for links in the HTML or Text bodies of this email. Possible options: None HtmlAndText HtmlOnly TextOnly
	TrackLinks string `json:",omitempty"`
	// Attachments: List of attachments
	Attachments []Attachment `json:",omitempty"`
	// Metadata: metadata
	Metadata map[string]string `json:",omitempty"`
	// MessageStream: MessageStream will default to the outbound message stream ID (Default Transactional Stream) if no message stream ID is provided.
	MessageStream string `json:",omitempty"`
}

// Header - an email header
type Header struct {
	// Name: header name
	Name string `json:",omitempty"`
	// Value: header value
	Value string `json:",omitempty"`
}

// Attachment is an optional encoded file to send along with an email
type Attachment struct {
	// Name: attachment name
	Name string
	// Content: Base64 encoded attachment data
	Content string
	// ContentType: attachment MIME type
	ContentType string
	// ContentId: populate for inlining images with the images cid
	ContentID string `json:",omitempty"`
}

// ErrEmailFailed is returned when email sending fails
var ErrEmailFailed = errors.New("email send failed")

// EmailResponse holds info in response to a send/send-batch request
// Even if API request comes back successful, check the ErrorCode to see if there might be a delivery problem
type EmailResponse struct {
	// To: Recipient email address
	To string
	// SubmittedAt: Timestamp
	SubmittedAt time.Time
	// MessageID: ID of message
	MessageID string
	// ErrorCode: see error codes here (https://postmarkapp.com/developer/api/overview#error-codes)
	ErrorCode int64
	// Message: Response message
	Message string
}

// SendEmail sends, well, an email.
func (client *Client) SendEmail(ctx context.Context, email Email) (EmailResponse, error) {
	res := EmailResponse{}
	err := client.doRequest(ctx, parameters{
		Method:    "POST",
		Path:      "email",
		Payload:   email,
		TokenType: serverToken,
	}, &res)

	if res.ErrorCode != 0 {
		return res, fmt.Errorf("%w: %v %s", ErrEmailFailed, res.ErrorCode, res.Message)
	}

	return res, err
}

// SendEmailBatch sends multiple emails together
// Individual emails in the batch can error, so it would be wise to
// range over the responses and sniff for errors
func (client *Client) SendEmailBatch(ctx context.Context, emails []Email) ([]EmailResponse, error) {
	var res []EmailResponse
	err := client.doRequest(ctx, parameters{
		Method:    "POST",
		Path:      "email/batch",
		Payload:   emails,
		TokenType: serverToken,
	}, &res)
	return res, err
}
