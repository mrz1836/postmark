package postmark

import (
	"context"
	"fmt"
	"net/http"
)

// WebhookHTTPAuth is an optional set of auth configuration to use when calling
// the webhook.
type WebhookHTTPAuth struct {
	// HTTP Auth username.
	Username string `json:"Username"`
	// HTTP Auth password.
	Password string `json:"Password"`
}

// WebhookTriggerEnabled holds configuration for webhooks which can only be
// enabled or disabled.
type WebhookTriggerEnabled struct {
	// Specifies whether this webhook is enabled.
	Enabled bool `json:"Enabled"`
}

// WebhookTriggerIncContent holds configuration for webhooks which can be
// enabled/disabled and optionally include message contents.
type WebhookTriggerIncContent struct {
	WebhookTriggerEnabled

	// Specifies whether the full content of the email is included in webhook POST.
	IncludeContent bool `json:"IncludeContent"`
}

// WebhookTriggerOpen holds configuration for the Open webhook.
type WebhookTriggerOpen struct {
	WebhookTriggerEnabled

	PostFirstOpenOnly bool `json:"PostFirstOpenOnly"`
}

// WebhookTrigger holds configuration for when this webhook should be called.
type WebhookTrigger struct {
	// List of open webhook details.
	Open WebhookTriggerOpen `json:"Open"`
	// List of click webhook details.
	Click WebhookTriggerEnabled `json:"Click"`
	// List of delivery webhook details.
	Delivery WebhookTriggerEnabled `json:"Delivery"`
	// List of bounce webhook details.
	Bounce WebhookTriggerIncContent `json:"Bounce"`
	// List of spam complaint webhook details.
	SpamComplaint WebhookTriggerIncContent `json:"SpamComplaint"`
	// List of subscription change webhook details.
	SubscriptionChange WebhookTriggerEnabled `json:"SubscriptionChange"`
}

// Webhook is a configured webhook on a message stream.
// https://postmarkapp.com/developer/api/webhooks-api#get-a-webhook
type Webhook struct {
	// ID of webhook.
	ID int `json:"ID,omitempty"`
	// Your webhook URL.
	URL string `json:"Url"`
	// The stream this webhook is associated with.
	MessageStream string `json:"MessageStream"`
	// Optional. HTTP Auth username and password.
	HTTPAuth *WebhookHTTPAuth `json:"HttpAuth,omitempty"`
	// Optional. List of custom headers included.
	HTTPHeaders []Header `json:"HttpHeaders,omitempty"`
	// List of different possible triggers a webhook can be enabled/disabled for.
	Triggers WebhookTrigger `json:"Triggers"`
}

// ListWebhooks returns all webhooks for a message stream. If the message stream
// is empty it will return all webhooks for the server. A non-existent message
// stream will result in an error.
func (client *Client) ListWebhooks(ctx context.Context, messageStream string) ([]Webhook, error) {
	msgStreamParam := ""
	if messageStream != "" {
		msgStreamParam = fmt.Sprintf("?MessageStream=%s", messageStream)
	}

	var res struct {
		Webhooks []Webhook
	}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodGet,
		Path:      "webhooks" + msgStreamParam,
		TokenType: serverToken,
	}, &res)
	return res.Webhooks, err
}

// GetWebhook retrieves a specific webhook by the webhook's ID.
func (client *Client) GetWebhook(ctx context.Context, id int) (Webhook, error) {
	var res Webhook
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodGet,
		Path:      fmt.Sprintf("webhooks/%d", id),
		TokenType: serverToken,
	}, &res)
	return res, err
}

// CreateWebhook makes a new Webhook. Do not specify the ID in the provided webhook. The
// returned webhook if successful will include the ID of the created webhook.
func (client *Client) CreateWebhook(ctx context.Context, webhook Webhook) (Webhook, error) {
	var res Webhook
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPost,
		Path:      "webhooks",
		Payload:   webhook,
		TokenType: serverToken,
	}, &res)
	return res, err
}

// EditWebhook alters an existing webhook. Do not specify the ID in the provided webhook. The
// returned webhook if successful will be the resulting state of after the edit.
func (client *Client) EditWebhook(ctx context.Context, id int, webhook Webhook) (Webhook, error) {
	var res Webhook
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPut,
		Path:      fmt.Sprintf("webhooks/%d", id),
		Payload:   webhook,
		TokenType: serverToken,
	}, &res)
	return res, err
}

// DeleteWebhook removes a webhook from the server.
func (client *Client) DeleteWebhook(ctx context.Context, id int) error {
	res := APIError{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodDelete,
		Path:      fmt.Sprintf("webhooks/%d", id),
		TokenType: serverToken,
	}, &res)

	if res.ErrorCode != 0 {
		return res
	}

	return err
}
