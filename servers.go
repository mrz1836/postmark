package postmark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Server represents a server registered in your Postmark account
type Server struct {
	// ID of server
	ID int64
	// Name of server
	Name string
	// APITokens associated with server.
	APITokens []string `json:"ApiTokens"`
	// ServerLink to your server overview page in Postmark.
	ServerLink string
	// Color of the server in the rack screen. Purple Blue Turquoise Green Red Yellow Grey
	Color string
	// SMTPAPIActivated specifies whether SMTP is enabled on this server.
	SMTPAPIActivated bool `json:"SmtpApiActivated"`
	// RawEmailEnabled allows raw email to be sent with inbound.
	RawEmailEnabled bool
	// DeliveryType specifies the type of environment for your server: Live or Sandbox, defaults to Live
	DeliveryType string `json:"DeliveryType"`
	// InboundAddress is the inbound email address
	InboundAddress string
	// InboundHookURL to POST to every time an inbound event occurs.
	InboundHookURL string `json:"InboundHookUrl"`
	// BounceHookURL to POST to every time a bounce event occurs.
	BounceHookURL string `json:"BounceHookUrl"`
	// OpenHookURL to POST to every time an open event occurs.
	OpenHookURL string `json:"OpenHookUrl"`
	// PostFirstOpenOnly - If set to true, only the first open by a particular recipient will initiate the open webhook. Any
	// subsequent opens of the same email by the same recipient will not initiate the webhook.
	PostFirstOpenOnly bool
	// TrackOpens indicates if all emails being sent through this server have open tracking enabled.
	TrackOpens bool
	// TrackLinks specifies link tracking in emails: None, HtmlAndText, HtmlOnly, TextOnly, defaults to "None"
	TrackLinks string `json:"TrackLinks"`
	// IncludeBounceContentInHook determines if bounce content is included in webhook.
	IncludeBounceContentInHook bool `json:"IncludeBounceContentInHook"`
	// InboundDomain is the inbound domain for MX setup
	InboundDomain string
	// InboundHash is the inbound hash of your inbound email address.
	InboundHash string
	// InboundSpamThreshold is the maximum spam score for an inbound message before it's blocked.
	InboundSpamThreshold int64
	// EnableSmtpApiErrorHooks specifies whether SMTP API Errors will be included with bounce webhooks.
	EnableSmtpApiErrorHooks bool `json:"EnableSmtpApiErrorHooks"`
}

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
func (client *Client) GetServer(ctx context.Context, serverID string) (Server, error) {
	res := Server{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodGet,
		Path:      fmt.Sprintf("servers/%s", serverID),
		TokenType: accountToken,
	}, &res)
	return res, err
}

// EditServer updates details for a specific server with serverID
func (client *Client) EditServer(ctx context.Context, serverID string, server Server) (Server, error) {
	// res := Server{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPut,
		Path:      fmt.Sprintf("servers/%s", serverID),
		TokenType: accountToken,
	}, &server)
	return server, err
}

// CreateServer creates a server
func (client *Client) CreateServer(ctx context.Context, server Server) (Server, error) {
	res := Server{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPost,
		Path:      "servers",
		TokenType: accountToken,
		Payload:   server,
	}, &res)
	return res, err
}
