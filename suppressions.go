package postmark

import (
	"context"
	"fmt"
	"time"
)

// SuppressionReasonType - The reason type of suppression
type SuppressionReasonType string

// OriginType - The reason type of origin
type OriginType string

// SuppressionUpdateStatus - The status of suppression update
type SuppressionUpdateStatus string

const (
	// HardBounceReason means an email sent to the address returned a hard bounce.
	HardBounceReason SuppressionReasonType = "HardBounce"

	// SpamComplaintReason means the recipient marked an email as spam.
	SpamComplaintReason SuppressionReasonType = "SpamComplaint"

	// ManualSuppressionReason means the recipient followed an unsubscribe link.
	ManualSuppressionReason SuppressionReasonType = "ManualSuppression"

	// RecipientOrigin means the email was added to the suppression list
	// as a result of the recipient's own action, e.g. by following an unsubscribe link.
	RecipientOrigin OriginType = "Recipient"

	// CustomerOrigin means the email was added to the suppression list as
	// the result of action by the Postmark account holder (e.g. Postmark's
	// customer).
	CustomerOrigin OriginType = "Customer"

	// AdminOrigin means the email was added to the suppression list as
	// the result of action by Postmark staff.
	AdminOrigin OriginType = "Admin"

	// SuppressionUpdateStatusSuppressed means the server successfully suppressed the email address.
	SuppressionUpdateStatusSuppressed SuppressionUpdateStatus = "Suppressed"

	// SuppressionUpdateStatusDeleted means the server successfully deleted the suppression.
	SuppressionUpdateStatusDeleted SuppressionUpdateStatus = "Deleted"

	// SuppressionUpdateStatusFailed means the server failed to update the suppression.
	SuppressionUpdateStatusFailed SuppressionUpdateStatus = "Failed"
)

// Suppression contains a suppressed email address for a particular message stream.
type Suppression struct {
	// EmailAddress is the address that is suppressed (can't be emailed any more)
	EmailAddress string

	// SuppressionReason is why the email address was added to the suppression list.
	// Possible options: HardBounce, SpamComplaint, ManualSuppression
	SuppressionReason SuppressionReasonType

	// Origin describes who added the email address to the suppression list.
	// Possible options: Recipient, Customer, Admin.
	Origin OriginType

	// CreatedAt is when the email address was added to the suppression list.
	CreatedAt time.Time
}

// SuppressionResponse contains a status of suppression creation or deletion.
type SuppressionResponse struct {
	// EmailAddress is the address that is suppressed (can't be emailed any more)
	EmailAddress string

	// Status of suppression creation or deletion.
	Status SuppressionUpdateStatus

	// If address cannot be suppressed or deleted (Status: Failed), the cause for failure is listed.
	// Otherwise, this field is null.
	Message string
}

// suppressionsResponse - A message received from the Postmark server
type suppressionsResponse struct {
	// Suppressions - The slice of suppression email address.
	Suppressions []Suppression
}

// updateSuppressionsResponse
type updateSuppressionsResponse struct {
	// Suppressions - The slice of suppression email address.
	Suppressions []SuppressionResponse
}

// suppressionsRequest - A message to send to the Postmark server for creating or deleting suppressions
type suppressionsRequest struct {
	Suppressions []Suppression
}

// GetSuppressions fetches email addresses in the list of suppression dump on the server
// It returns a Suppressions slice, and any error that occurred
// https://postmarkapp.com/developer/api/suppressions-api#suppression-dump
func (client *Client) GetSuppressions(
	ctx context.Context,
	streamID string,
	options map[string]interface{},
) ([]Suppression, error) {
	path := fmt.Sprintf("message-streams/%s/suppressions/dump", streamID)

	res := suppressionsResponse{}
	err := client.get(ctx, buildURL(path, options), &res)
	return res.Suppressions, err
}

// CreateSuppressions creates email addresses in the suppression list on the server.
func (client *Client) CreateSuppressions(
	ctx context.Context,
	streamID string,
	suppressions []Suppression,
) ([]SuppressionResponse, error) {
	res := updateSuppressionsResponse{}
	path := fmt.Sprintf("message-streams/%s/suppressions", streamID)
	err := client.post(ctx, path, suppressionsRequest{Suppressions: suppressions}, &res)
	return res.Suppressions, err
}

// DeleteSuppressions deletes email addresses from the suppression list on the server.
// SpamComplaint suppressions cannot be deleted.
func (client *Client) DeleteSuppressions(
	ctx context.Context,
	streamID string,
	suppressions []Suppression,
) ([]SuppressionResponse, error) {
	res := updateSuppressionsResponse{}
	path := fmt.Sprintf("message-streams/%s/suppressions/delete", streamID)
	err := client.post(ctx, path, suppressionsRequest{Suppressions: suppressions}, &res)
	return res.Suppressions, err
}