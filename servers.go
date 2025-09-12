package postmark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Server represents a server registered in your Postmark account
type Server struct {
	// ID of server
	ID int64 `json:"ID"`
	// Name of server
	Name string `json:"Name"`
	// APITokens associated with server.
	APITokens []string `json:"ApiTokens"`
	// ServerLink to your server overview page in Postmark.
	ServerLink string `json:"ServerLink"`
	// Color of the server in the rack screen. Purple Blue Turquoise Green Red Yellow Grey
	Color string `json:"Color"`
	// SMTPAPIActivated specifies whether SMTP is enabled on this server.
	SMTPAPIActivated bool `json:"SmtpApiActivated"`
	// RawEmailEnabled allows raw email to be sent with inbound.
	RawEmailEnabled bool `json:"RawEmailEnabled"`
	// DeliveryType specifies the type of environment for your server: Live or Sandbox, defaults to Live
	DeliveryType string `json:"DeliveryType"`
	// InboundAddress is the inbound email address
	InboundAddress string `json:"InboundAddress"`
	// InboundHookURL to POST to every time an inbound event occurs.
	InboundHookURL string `json:"InboundHookUrl"`
	// Deprecated: Use the Bounce Webhook API instead.
	BounceHookURL string `json:"BounceHookUrl"`
	// Deprecated: Use the Open Tracking Webhook API instead.
	OpenHookURL string `json:"OpenHookUrl"`
	// Deprecated: Use the Delivery Webhook API instead.
	DeliveryHookURL string `json:"DeliveryHookUrl"`
	// PostFirstOpenOnly - If set to true, only the first open by a particular recipient will initiate the open webhook. Any
	// subsequent opens of the same email by the same recipient will not initiate the webhook.
	PostFirstOpenOnly bool `json:"PostFirstOpenOnly"`
	// TrackOpens indicates if all emails being sent through this server have open tracking enabled.
	TrackOpens bool `json:"TrackOpens"`
	// TrackLinks specifies link tracking in emails: None, HtmlAndText, HtmlOnly, TextOnly, defaults to "None"
	TrackLinks string `json:"TrackLinks"`
	// IncludeBounceContentInHook determines if bounce content is included in webhook.
	IncludeBounceContentInHook bool `json:"IncludeBounceContentInHook"`
	// InboundDomain is the inbound domain for MX setup
	InboundDomain string `json:"InboundDomain"`
	// InboundHash is the inbound hash of your inbound email address.
	InboundHash string `json:"InboundHash"`
	// InboundSpamThreshold is the maximum spam score for an inbound message before it's blocked.
	InboundSpamThreshold int64 `json:"InboundSpamThreshold"`
	// EnableSMTPAPIErrorHooks specifies whether SMTP API Errors will be included with bounce webhooks.
	EnableSMTPAPIErrorHooks bool `json:"EnableSmtpApiErrorHooks"`
}

// ServerCreateRequest represents the fields to create a server
type ServerCreateRequest struct {
	// Name of server
	Name string `json:"Name" binding:"required"`
	// Color of the server in the server list, for quick identification. Purple Blue Turquoise Green Red Yellow Grey Orange
	Color string `json:"Color"`
	// SMTPAPIActivated specifies whether SMTP is enabled on this server.
	SMTPAPIActivated bool `json:"SmtpApiActivated"`
	// When enabled, the raw email content will be included with inbound webhook payloads under the RawEmail key.
	RawEmailEnabled bool `json:"RawEmailEnabled"`
	// Specifies the type of environment for your server. Possible options: Live Sandbox. Defaults to Live if not
	// specified. This cannot be changed after the server has been created.
	DeliveryType string `json:"DeliveryType"`
	// URL to POST to every time an inbound event occurs.
	InboundHookURL string `json:"InboundHookUrl"`
	// Deprecated: Use the Bounce Webhook API instead.
	BounceHookURL string `json:"BounceHookUrl"`
	// Deprecated: Use the Open Tracking Webhook API instead.
	OpenHookURL string `json:"OpenHookUrl"`
	// Deprecated: Use the Delivery Webhook API instead.
	DeliveryHookURL string `json:"DeliveryHookUrl"`
	// Deprecated: Use the Click Webhook API instead.
	ClickHookURL string `json:"ClickHookUrl"`
	// PostFirstOpenOnly - If set to true, only the first open by a particular recipient will initiate the open webhook. Any
	// subsequent opens of the same email by the same recipient will not initiate the webhook.
	PostFirstOpenOnly bool `json:"PostFirstOpenOnly"`
	// InboundDomain is the inbound domain for MX setup
	InboundDomain string `json:"InboundDomain"`
	// InboundSpamThreshold is the maximum spam score for an inbound message before it's blocked.
	InboundSpamThreshold int64 `json:"InboundSpamThreshold"`
	// TrackOpens indicates if all emails being sent through this server have open tracking enabled.
	TrackOpens bool `json:"TrackOpens"`
	// TrackLinks specifies link tracking in emails: None, HtmlAndText, HtmlOnly, TextOnly, defaults to "None"
	TrackLinks string `json:"TrackLinks"`
	// IncludeBounceContentInHook determines if bounce content is included in webhook.
	IncludeBounceContentInHook bool `json:"IncludeBounceContentInHook"`
	// EnableSMTPAPIErrorHooks specifies whether SMTP API Errors will be included with bounce webhooks.
	EnableSMTPAPIErrorHooks bool `json:"EnableSmtpApiErrorHooks"`
}

// ServerEditRequest represents the fields that can be updated for a server
type ServerEditRequest struct {
	// Name of server
	Name string `json:"Name" binding:"required"`
	// Color of the server in the server list, for quick identification. Purple Blue Turquoise Green Red Yellow Grey Orange
	Color string `json:"Color"`
	// SMTPAPIActivated specifies whether SMTP is enabled on this server.
	SMTPAPIActivated bool `json:"SmtpApiActivated"`
	// When enabled, the raw email content will be included with inbound webhook payloads under the RawEmail key.
	RawEmailEnabled bool `json:"RawEmailEnabled"`
	// URL to POST to every time an inbound event occurs.
	InboundHookURL string `json:"InboundHookUrl"`
	// Deprecated: Use the Bounce Webhook API instead.
	BounceHookURL string `json:"BounceHookUrl"`
	// Deprecated: Use the Open Tracking Webhook API instead.
	OpenHookURL string `json:"OpenHookUrl"`
	// Deprecated: Use the Delivery Webhook API instead.
	DeliveryHookURL string `json:"DeliveryHookUrl"`
	// Deprecated: Use the Click Webhook API instead.
	ClickHookURL string `json:"ClickHookUrl"`
	// PostFirstOpenOnly - If set to true, only the first open by a particular recipient will initiate the open webhook. Any
	// subsequent opens of the same email by the same recipient will not initiate the webhook.
	PostFirstOpenOnly bool `json:"PostFirstOpenOnly"`
	// InboundDomain is the inbound domain for MX setup
	InboundDomain string `json:"InboundDomain"`
	// InboundSpamThreshold is the maximum spam score for an inbound message before it's blocked.
	InboundSpamThreshold int64 `json:"InboundSpamThreshold"`
	// TrackOpens indicates if all emails being sent through this server have open tracking enabled.
	TrackOpens bool `json:"TrackOpens"`
	// TrackLinks specifies link tracking in emails: None, HtmlAndText, HtmlOnly, TextOnly, defaults to "None"
	TrackLinks string `json:"TrackLinks"`
	// IncludeBounceContentInHook determines if bounce content is included in webhook.
	IncludeBounceContentInHook bool `json:"IncludeBounceContentInHook"`
	// EnableSMTPAPIErrorHooks specifies whether SMTP API Errors will be included with bounce webhooks.
	EnableSMTPAPIErrorHooks bool `json:"EnableSmtpApiErrorHooks"`
}

// ServersList is just a list of Server as they are in the response
type ServersList struct {
	TotalCount int
	Servers    []Server
}

// MarshalJSON customizes the JSON representation of the Server struct by setting default values for specific fields.
func (s Server) MarshalJSON() ([]byte, error) {
	type Aux Server

	// If TrackLinks is empty, set it to "None"
	trackLinks := s.TrackLinks
	if trackLinks == "" {
		trackLinks = "None"
	}

	// If DeliveryType is empty, set it to default value "Live"
	deliveryType := s.DeliveryType
	if deliveryType == "" {
		deliveryType = "Live"
	}

	return json.Marshal(&struct {
		Aux

		TrackLinks   string `json:"TrackLinks"`
		DeliveryType string `json:"DeliveryType"`
	}{
		Aux:          Aux(s),
		TrackLinks:   trackLinks,
		DeliveryType: deliveryType,
	})
}

// GetServer fetches a specific server via serverID
func (client *Client) GetServer(ctx context.Context, serverID int64) (Server, error) {
	res := Server{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodGet,
		Path:      fmt.Sprintf("servers/%d", serverID),
		TokenType: accountToken,
	}, &res)
	return res, err
}

// GetServers fetches a list of servers on the account, limited by count and paged by offset
// Optionally filter by a specific server name. Note that this is a string search, so MyServer will match
// MyServer, MyServer Production, and MyServer Test.
func (client *Client) GetServers(ctx context.Context, count, offset int64, name string) (ServersList, error) {
	res := ServersList{}

	values := &url.Values{}
	values.Add("count", fmt.Sprintf("%d", count))
	values.Add("offset", fmt.Sprintf("%d", offset))

	if name != "" {
		values.Add("name", name)
	}

	err := client.doRequest(ctx, parameters{
		Method:    "GET",
		Path:      fmt.Sprintf("servers?%s", values.Encode()),
		TokenType: accountToken,
	}, &res)
	return res, err
}

// EditServer updates details for a specific server with serverID
func (client *Client) EditServer(ctx context.Context, serverID int64, request ServerEditRequest) (Server, error) {
	res := Server{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPut,
		Path:      fmt.Sprintf("servers/%d", serverID),
		TokenType: accountToken,
		Payload:   request,
	}, &res)
	return res, err
}

// CreateServer creates a server
func (client *Client) CreateServer(ctx context.Context, request ServerCreateRequest) (Server, error) {
	res := Server{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPost,
		Path:      "servers",
		TokenType: accountToken,
		Payload:   request,
	}, &res)
	return res, err
}

// DeleteServer removes a server.
func (client *Client) DeleteServer(ctx context.Context, serverID int64) error {
	res := APIError{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodDelete,
		Path:      fmt.Sprintf("servers/%d", serverID),
		TokenType: accountToken,
	}, &res)

	if res.ErrorCode != 0 {
		return res
	}

	return err
}
