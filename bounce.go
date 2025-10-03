package postmark

import (
	"context"
	"fmt"
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
	err := client.get(ctx, "deliverystats", &res)
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

	if options == nil {
		options = make(map[string]interface{})
	}

	options["count"] = count
	options["offset"] = offset

	err := client.get(ctx, buildURL("bounces", options), &res)
	return res.Bounces, res.TotalCount, err
}

// GetBounce fetches a single bounce with bounceID
func (client *Client) GetBounce(ctx context.Context, bounceID int64) (Bounce, error) {
	res := Bounce{}
	err := client.get(ctx, fmt.Sprintf("bounces/%v", bounceID), &res)
	return res, err
}

type dumpResponse struct {
	Body string
}

// GetBounceDump fetches an SMTP data dump for a single bounce
func (client *Client) GetBounceDump(ctx context.Context, bounceID int64) (string, error) {
	res := dumpResponse{}
	err := client.get(ctx, fmt.Sprintf("bounces/%v/dump", bounceID), &res)
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
	err := client.put(ctx, fmt.Sprintf("bounces/%v/activate", bounceID), nil, &res)
	return res.Bounce, res.Message, err
}

// GetBouncedTags retrieves a list of tags that have generated bounced emails
func (client *Client) GetBouncedTags(ctx context.Context) ([]string, error) {
	var res []string
	err := client.get(ctx, "bounces/tags", &res)
	return res, err
}