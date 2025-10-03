package postmark

import (
	"context"
	"fmt"
	"net/mail"
	"time"
)

// InboundMessage - a message received from the Postmark server
type InboundMessage struct {
	// From - The sender email address.
	From string
	// FromName - The sender name.
	FromName string
	// FromFull - Sender email address and name.
	FromFull Recipient
	// To - Inbound address the message was sent to.
	To string
	// ToFull - Slice of all TO recipients
	ToFull []Recipient
	// CcFull - Slice of all CC recipients
	CcFull []Recipient
	// Cc - Cc recipient email address.
	Cc string
	// ReplyTo - Reply to override email address.
	ReplyTo string
	// OriginalRecipient - Receiver (RCPT TO) address this webhook is for.
	OriginalRecipient string
	// Subject - Email subject
	Subject string
	// Date - Timestamp
	Date string
	// MailboxHash - Custom hash that the email was sent to.
	MailboxHash string
	// TextBody - Plain text email message.
	TextBody string
	// HTMLBody - HTML email message.
	HTMLBody string `json:"HtmlBody"`
	// Tag - Tag name
	Tag string
	// Headers - List of objects that each represent a header name and value.
	Headers []Header
	// Attachments - List of objects that each represent an attachment.
	Attachments []Attachment
	// MessageID - Unique ID of the message.
	MessageID string
	// BlockedReason - Reason message was blocked.
	BlockedReason string
	// Status - Status of message in your Postmark activity.
	Status string
}

// Time returns a parsed time.Time struct
func (x InboundMessage) Time() (time.Time, error) {
	return mail.ParseDate(x.Date)
}

// GetInboundMessage fetches a specific inbound message via serverID
func (client *Client) GetInboundMessage(ctx context.Context, messageID string) (InboundMessage, error) {
	res := InboundMessage{}
	err := client.get(ctx, fmt.Sprintf("messages/inbound/%s/details", messageID), &res)
	return res, err
}

type inboundMessagesResponse struct {
	TotalCount int64            `json:"TotalCount"`
	Messages   []InboundMessage `json:"InboundMessages"`
}

// GetInboundMessages fetches a list of inbound message on the server
// It returns a InboundMessage slice, the total message count, and any error that occurred
// http://developer.postmarkapp.com/developer-api-messages.html#inbound-message-search
func (client *Client) GetInboundMessages(ctx context.Context, count, offset int64, options map[string]interface{}) ([]InboundMessage, int64, error) {
	res := inboundMessagesResponse{}

	if options == nil {
		options = make(map[string]interface{})
	}

	options["count"] = count
	options["offset"] = offset

	err := client.get(ctx, buildURL("messages/inbound", options), &res)

	return res.Messages, res.TotalCount, err
}

// BypassInboundMessage - Bypass rules for a blocked inbound message
func (client *Client) BypassInboundMessage(ctx context.Context, messageID string) error {
	res := APIError{}
	err := client.put(ctx, fmt.Sprintf("messages/inbound/%s/bypass", messageID), nil, &res)
	if err != nil {
		return err
	}
	if res.ErrorCode != 0 {
		return res
	}
	return nil
}

// RetryInboundMessage - Retry a failed inbound message for processing
func (client *Client) RetryInboundMessage(ctx context.Context, messageID string) error {
	res := APIError{}
	err := client.put(ctx, fmt.Sprintf("messages/inbound/%s/retry", messageID), nil, &res)
	if err != nil {
		return err
	}
	if res.ErrorCode != 0 {
		return res
	}
	return nil
}