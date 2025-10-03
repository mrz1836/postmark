package postmark

import (
	"context"
	"fmt"
	"time"
)

// DataRemovalRequest represents a request to remove recipient data from Postmark
type DataRemovalRequest struct {
	// Recipient: The email address of the recipient whose data should be removed
	Recipient string `json:"Recipient"`
}

// DataRemovalResponse represents the response from a data removal request
type DataRemovalResponse struct {
	// ID: Unique ID of the data removal request
	ID int64 `json:"ID"`
	// Recipient: Email address of the recipient whose data is being removed
	Recipient string `json:"Recipient"`
	// RequestedAt: Timestamp when the removal was requested
	RequestedAt time.Time `json:"RequestedAt"`
	// Status: Current status of the removal (Pending, Processing, Completed, Failed)
	Status string `json:"Status"`
	// CompletedAt: Timestamp when the removal was completed (if applicable)
	CompletedAt *time.Time `json:"CompletedAt,omitempty"`
}

// CreateDataRemoval creates a new data removal request
func (client *Client) CreateDataRemoval(ctx context.Context, request DataRemovalRequest) (DataRemovalResponse, error) {
	res := DataRemovalResponse{}
	err := client.postWithAccountToken(ctx, "data-removals", request, &res)
	return res, err
}

// GetDataRemovalStatus checks the status of a data removal request
func (client *Client) GetDataRemovalStatus(ctx context.Context, id int64) (DataRemovalResponse, error) {
	res := DataRemovalResponse{}
	err := client.getWithAccountToken(ctx, fmt.Sprintf("data-removals/%d", id), &res)
	return res, err
}