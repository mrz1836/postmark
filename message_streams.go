package postmark

import (
	"context"
	"fmt"
)

// MessageStreamType is an Enum representing the type of a message stream.
type MessageStreamType string

// MessageStreamUnsubscribeHandlingType is an Enum with the possible values for
// the unsubscribe handling in a message stream.
type MessageStreamUnsubscribeHandlingType string

const (
	// InboundMessageStreamType indicates a message stream is for inbound messages.
	InboundMessageStreamType MessageStreamType = "Inbound"
	// BroadcastMessageStreamType indicates a message stream is for broadcast messages.
	BroadcastMessageStreamType MessageStreamType = "Broadcasts"
	// TransactionalMessageStreamType indicates a message stream is for transactional messages.
	TransactionalMessageStreamType MessageStreamType = "Transactional"

	// NoneUnsubscribeHandlingType indicates a message stream unsubscribe
	// handling will be performed by the user.
	NoneUnsubscribeHandlingType MessageStreamUnsubscribeHandlingType = "None"
	// PostmarkUnsubscribeHandlingType indicates a message stream unsubscribe
	// handling will be performed by postmark.
	PostmarkUnsubscribeHandlingType MessageStreamUnsubscribeHandlingType = "Postmark"
	// CustomUnsubscribeHandlingType indicates a message stream unsubscribe
	// handling is custom.
	CustomUnsubscribeHandlingType MessageStreamUnsubscribeHandlingType = "Custom"
)

// MessageStreamSubscriptionManagementConfiguration is the configuration for
// subscriptions to the message stream.
type MessageStreamSubscriptionManagementConfiguration struct {
	// The unsubscribe management option used for the Stream. Broadcast Message
	// Streams require unsubscribe management, Postmark is default. For Inbound
	// and Transactional Streams default is none.
	UnsubscribeHandlingType MessageStreamUnsubscribeHandlingType `json:"UnsubscribeHandlingType"`
}

// MessageStream holes the configuration for a message stream on a server.
// https://postmarkapp.com/developer/api/message-streams-api
type MessageStream struct {
	// ID of message stream.
	ID string `json:"ID"`
	// ID of server the message stream is associated with.
	ServerID int `json:"ServerID"`
	// Name of message stream.
	Name string `json:"Name"`
	// Description of message stream. This value can be null.
	Description *string `json:"Description,omitempty"`
	// Type of message stream.
	MessageStreamType MessageStreamType `json:"MessageStreamType"`
	// Timestamp when message stream was created.
	CreatedAt string `json:"CreatedAt"`
	// Timestamp when message stream was last updated. This value can be null.
	UpdatedAt *string `json:"UpdatedAt,omitempty"`
	// Timestamp when message stream was archived. This value can be null.
	ArchivedAt *string `json:"ArchivedAt,omitempty"`
	// Archived streams are deleted 45 days after archiving date. Until this
	// date, it can be restored. This value is null if the stream is not
	// archived.
	ExpectedPurgeDate *string `json:"ExpectedPurgeDate,omitempty"`
	// Subscription management options for the Stream
	SubscriptionManagementConfiguration MessageStreamSubscriptionManagementConfiguration `json:"SubscriptionManagementConfiguration"`
}

// ListMessageStreams returns all message streams for a server.
// messageStreamType must be one of "All", "Inbound", "Transactional",
// "Broadcasts" and defaults to "All".
func (client *Client) ListMessageStreams(ctx context.Context, messageStreamType string, includeArchived bool) ([]MessageStream, error) {
	switch messageStreamType {
	case "Inbound", "Transactional", "Broadcasts":
		break
	default:
		messageStreamType = "All"
	}

	var res struct {
		MessageStreams []MessageStream
	}

	err := client.get(ctx, fmt.Sprintf("message-streams?MessageStreamType=%s&IncludeArchivedStreams=%t", messageStreamType, includeArchived), &res)

	return res.MessageStreams, err
}

// GetMessageStream retrieves a specific message stream by the message stream's ID.
func (client *Client) GetMessageStream(ctx context.Context, id string) (MessageStream, error) {
	var res MessageStream
	err := client.get(ctx, fmt.Sprintf("message-streams/%s", id), &res)
	return res, err
}

// EditMessageStreamRequest is the request body for EditMessageStream. It
// contains only a subset of the fields of MessageStream.
type EditMessageStreamRequest struct {
	// Name of message stream.
	Name string `json:"Name"`
	// Description of message stream. This value can be null.
	Description *string `json:"Description,omitempty"`
	// Subscription management options for the Stream
	SubscriptionManagementConfiguration MessageStreamSubscriptionManagementConfiguration `json:"SubscriptionManagementConfiguration"`
}

// EditMessageStream updates a message stream.
func (client *Client) EditMessageStream(ctx context.Context, id string, req EditMessageStreamRequest) (MessageStream, error) {
	var res MessageStream
	err := client.patch(ctx, fmt.Sprintf("message-streams/%s", id), req, &res)
	return res, err
}

// CreateMessageStreamRequest is the request body for CreateMessageStream. It
// contains only a subset of the fields of MessageStream.
type CreateMessageStreamRequest struct {
	// ID of message stream.
	ID string `json:"ID"`
	// Name of message stream.
	Name string `json:"Name"`
	// Description of message stream. This value can be null.
	Description *string `json:"Description,omitempty"`
	// Type of message stream.
	MessageStreamType MessageStreamType `json:"MessageStreamType"`
	// Subscription management options for the Stream
	SubscriptionManagementConfiguration MessageStreamSubscriptionManagementConfiguration `json:"SubscriptionManagementConfiguration"`
}

// CreateMessageStream makes a new message stream. It will be created on the
// server of the token used by this Client.
func (client *Client) CreateMessageStream(ctx context.Context, req CreateMessageStreamRequest) (MessageStream, error) {
	var res MessageStream
	err := client.post(ctx, "message-streams", req, &res)
	return res, err
}

// ArchiveMessageStreamResponse is the response body for ArchiveMessageStream.
type ArchiveMessageStreamResponse struct {
	// ID of message stream.
	ID string `json:"ID"`
	// Server ID of message stream.
	ServerID int `json:"ServerID"`
	// Expected purge date of message stream. Stream is deleted 45 days after
	// archiving date. Until this date, it can be restored.
	ExpectedPurgeDate string `json:"ExpectedPurgeDate"`
}

// ArchiveMessageStream archives a message stream. Archived streams are deleted
// after 45 days, but they can be restored until that point.
func (client *Client) ArchiveMessageStream(ctx context.Context, id string) (ArchiveMessageStreamResponse, error) {
	var res ArchiveMessageStreamResponse
	err := client.post(ctx, fmt.Sprintf("message-streams/%s/archive", id), nil, &res)
	return res, err
}

// UnarchiveMessageStream unarchives a message stream if it has not been deleted yet.
// The ArchivedAt value will be null after calling this method.
func (client *Client) UnarchiveMessageStream(ctx context.Context, id string) (MessageStream, error) {
	var res MessageStream
	err := client.post(ctx, fmt.Sprintf("message-streams/%s/unarchive", id), nil, &res)
	return res, err
}
