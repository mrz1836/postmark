package postmark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// BounceType represents a type of bounce, and how many bounces have occurred
// http://developer.postmarkapp.com/developer-api-bounce.html#bounce-types
type BounceType struct {
	// Type: bounce type identifier
	Type string
	// Name: full name of the bounce type
	Name string
	// Count: how many bounces have occurred
	Count int64
}

// DeliveryStats represents bounce stats
type DeliveryStats struct {
	// InactiveMails: Number of inactive emails
	InactiveMails int64
	// Bounces: List of bounce types with total counts.
	Bounces []BounceType
}

// GetDeliveryStats returns delivery stats for the server
func (client *Client) GetDeliveryStats(ctx context.Context) (DeliveryStats, error) {
	res := DeliveryStats{}
	path := "deliverystats"
	err := client.doRequest(ctx, parameters{
		Method:    "GET",
		Path:      path,
		TokenType: serverToken,
	}, &res)
	return res, err
}

// Bounce represents a specific delivery failure
type Bounce struct {
	// RecordType: Type of record (bounce)
	RecordType string
	// ID: ID of bounce
	ID int64
	// Type: Bounce type
	Type string
	// TypeCode: Bounce type code
	TypeCode int64
	// Name: Bounce type name
	Name string
	// Tag: Tag name
	Tag string
	// MessageID: ID of message
	MessageID string
	// MessageStream: Message stream ID
	MessageStream string
	// Description: Description of bounce
	Description string
	// Details: Details on the bounce
	Details string
	// Email: Email address that bounced
	Email string
	// BouncedAt: Timestamp of bounce
	BouncedAt time.Time
	// DumpAvailable: Specifies whether you can get a raw dump from this bounce. Postmark does not store bounce dumps older than 30 days.
	DumpAvailable bool
	// Inactive: Specifies if the bounce caused Postmark to deactivate this email.
	Inactive bool
	// CanActivate: Specifies whether you are able to reactivate this email.
	CanActivate bool
	// Subject: Email subject
	Subject string
	// Content: Raw email content
	Content string
}

type bouncesResponse struct {
	TotalCount int64
	Bounces    []Bounce
}

// GetBounces returns bounces for the server
// It returns a Bounce slice, the total bounce count, and any error that occurred
// Available options: http://developer.postmarkapp.com/developer-api-bounce.html#bounces
func (client *Client) GetBounces(ctx context.Context, count, offset int64, options map[string]interface{}) ([]Bounce, int64, error) {
	res := bouncesResponse{}

	values := &url.Values{}
	values.Add("count", fmt.Sprintf("%d", count))
	values.Add("offset", fmt.Sprintf("%d", offset))

	for k, v := range options {
		values.Add(k, fmt.Sprintf("%v", v))
	}

	path := fmt.Sprintf("bounces?%s", values.Encode())

	err := client.doRequest(ctx, parameters{
		Method:    "GET",
		Path:      path,
		TokenType: serverToken,
	}, &res)
	return res.Bounces, res.TotalCount, err
}

// GetBounce fetches a single bounce with bounceID
func (client *Client) GetBounce(ctx context.Context, bounceID int64) (Bounce, error) {
	res := Bounce{}
	path := fmt.Sprintf("bounces/%v", bounceID)
	err := client.doRequest(ctx, parameters{
		Method:    "GET",
		Path:      path,
		TokenType: serverToken,
	}, &res)
	return res, err
}

type dumpResponse struct {
	Body string
}

// GetBounceDump fetches an SMTP data dump for a single bounce
func (client *Client) GetBounceDump(ctx context.Context, bounceID int64) (string, error) {
	res := dumpResponse{}
	path := fmt.Sprintf("bounces/%v/dump", bounceID)
	err := client.doRequest(ctx, parameters{
		Method:    "GET",
		Path:      path,
		TokenType: serverToken,
	}, &res)
	return res.Body, err
}

type activateBounceResponse struct {
	Message string
	Bounce  Bounce
}

// ActivateBounce reactivates a bounce for resending. Returns the bounce, a
// message, and any error that occurs
// TODO: clarify this with Postmark
func (client *Client) ActivateBounce(ctx context.Context, bounceID int64) (Bounce, string, error) {
	res := activateBounceResponse{}
	path := fmt.Sprintf("bounces/%v/activate", bounceID)
	err := client.doRequest(ctx, parameters{
		Method:    "PUT",
		Path:      path,
		TokenType: serverToken,
	}, &res)
	return res.Bounce, res.Message, err
}

type bouncedTagsResponse struct {
	Tags []string `json:"tags"`
}

// GetBouncedTags retrieves a list of tags that have generated bounced emails
func (client *Client) GetBouncedTags(ctx context.Context) ([]string, error) {
	var raw json.RawMessage
	path := "bounces/tags"
	err := client.doRequest(ctx, parameters{
		Method:    "GET",
		Path:      path,
		TokenType: serverToken,
	}, &raw)
	if err != nil {
		return []string{}, err
	}

	// PM returns this payload in an impossible to unmarshal way
	// ["tag1","tag2","tag3"]. So let's rejigger it to make it possible.
	jsonString := fmt.Sprintf(`{"tags": %s}`, string(raw))
	res := bouncedTagsResponse{}
	err = json.Unmarshal([]byte(jsonString), &res)

	return res.Tags, err
}
