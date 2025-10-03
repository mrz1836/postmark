package postmark

import (
	"context"
	"fmt"
	"net/url"
)

// InboundRuleTrigger represents an inbound rule trigger
type InboundRuleTrigger struct {
	// ID: Unique ID of the trigger
	ID int64
	// Rule: Email address or domain to block
	Rule string
}

// InboundRulesTriggersResponse represents the response from listing inbound rule triggers
type InboundRulesTriggersResponse struct {
	// TotalCount: Total matching triggers
	TotalCount int64
	// InboundRules: List of inbound rules
	InboundRules []InboundRuleTrigger
}

// InboundRuleTriggerCreateRequest represents the request to create an inbound rule trigger
type InboundRuleTriggerCreateRequest struct {
	// Rule: Email address or domain to block (required)
	Rule string `json:"Rule"`
}

// GetInboundRuleTriggers fetches a list of inbound rule triggers on the server
// It returns a slice of InboundRuleTrigger, the total trigger count, and any error that occurred
func (client *Client) GetInboundRuleTriggers(ctx context.Context, count, offset int64) ([]InboundRuleTrigger, int64, error) {
	res := InboundRulesTriggersResponse{}

	values := &url.Values{}
	values.Add("count", fmt.Sprintf("%d", count))
	values.Add("offset", fmt.Sprintf("%d", offset))

	err := client.get(ctx, buildURLWithQuery("triggers/inboundrules", *values), &res)

	return res.InboundRules, res.TotalCount, err
}

// CreateInboundRuleTrigger creates an inbound rule trigger to block emails
func (client *Client) CreateInboundRuleTrigger(ctx context.Context, rule string) (InboundRuleTrigger, error) {
	res := InboundRuleTrigger{}

	requestData := InboundRuleTriggerCreateRequest{
		Rule: rule,
	}

	err := client.post(ctx, "triggers/inboundrules", requestData, &res)

	return res, err
}

// DeleteInboundRuleTrigger deletes an inbound rule trigger by ID
func (client *Client) DeleteInboundRuleTrigger(ctx context.Context, triggerID int64) error {
	res := APIError{}
	err := client.delete(ctx, fmt.Sprintf("triggers/inboundrules/%d", triggerID), &res)
	if err != nil {
		return err
	}
	if res.ErrorCode != 0 {
		return res
	}
	return nil
}
