package postmark

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// SenderSignature contains the brief details of a sender signature associated with your account.
type SenderSignature struct {
	// Domain associated with sender signature.
	Domain string `json:"Domain"`
	// Email address associated with sender signature.
	FromEmail string `json:"EmailAddress"`
	// Reply-To email associated with sender signature.
	ReplyToEmail string `json:"ReplyToEmailAddress"`
	// From name of sender signature.
	Name string `json:"Name"`
	// Indicates whether this sender signature has been confirmed.
	Confirmed bool `json:"Confirmed"`
	// Unique ID of sender signature.
	ID int64 `json:"ID"`
}

// SenderSignatureDetails contains the full details of a sender signature associated with your account.
type SenderSignatureDetails struct {
	// Domain associated with sender signature.
	Domain string `json:"Domain"`
	// Email address associated with sender signature.
	FromEmail string `json:"EmailAddress"`
	// Reply-To email associated with sender signature.
	ReplyToEmail string `json:"ReplyToEmailAddress"`
	// From name of sender signature.
	Name string `json:"Name"`
	// Indicates whether this sender signature has been confirmed.
	Confirmed bool `json:"Confirmed"`
	// Deprecated: See our [blog post](https://postmarkapp.com/blog/why-we-no-longer-ask-for-spf-records) to learn
	// why this field was deprecated.
	SPFVerified bool `json:"SPFVerified"`
	// Host name used for the SPF configuration.
	SPFHost string `json:"SPFHost"`
	// Value that can be optionally setup with your DNS host.
	// See our [blog post](https://postmarkapp.com/blog/why-we-no-longer-ask-for-spf-records) to learn why this field is no longer necessary.
	SPFTextValue string `json:"SPFTextValue"`
	// Specifies whether DKIM has ever been verified for the domain or not. Once DKIM is verified, this response will
	// stay true, even if the record is later removed from DNS.
	DKIMVerified bool `json:"DKIMVerified"`
	// DKIM is using a strength weaker than 1024 bit. If so, it’s possible to request a new DKIM using the
	// [RequestNewDKIM](https://postmarkapp.com/developer/api/signatures-api#request-dkim) function below.
	WeakDKIM bool `json:"WeakDKIM"`
	// DNS TXT host being used to validate messages sent in.
	DKIMHost string `json:"DKIMHost"`
	// DNS TXT value being used to validate messages sent in.
	DKIMTextValue string `json:"DKIMTextValue"`
	// If a DKIM renewal has been initiated or this DKIM is from a new Sender Signature, this field will show the pending
	// DKIM DNS TXT host which has yet to be setup and confirmed at your registrar or DNS host.
	DKIMPendingHost string `json:"DKIMPendingHost"`
	// Similar to the DKIMPendingHost field, this will show the DNS TXT value waiting to be confirmed at your registrar
	// or DNS host.
	DKIMPendingTextValue string `json:"DKIMPendingTextValue"`
	// Once a new DKIM has been confirmed at your registrar or DNS host, Postmark will revoke the old DKIM host in
	// preparation for removing it permanently from the system.
	DKIMRevokedHost string `json:"DKIMRevokedHost"`
	// Similar to DKIMRevokedHost, this field will show the DNS TXT value that will soon be removed from the Postmark system.
	DKIMRevokedTextValue string `json:"DKIMRevokedTextValue"`
	// Indicates whether you may safely delete the old DKIM DNS TXT records at your registrar or DNS host.
	// The new DKIM is now safely in use.
	SafeToRemoveRevokedKeyFromDNS bool `json:"SafeToRemoveRevokedKeyFromDNS"`
	// While DKIM renewal or new DKIM operations are being conducted or setup, this field will indicate Pending. After
	// all DNS TXT records are up to date and any pending renewal operations are finished, it will indicate Verified.
	DKIMUpdateStatus string `json:"DKIMUpdateStatus"`
	// The custom Return-Path domain for this signature. For more information about this field, please
	// [read our support page](http://support.postmarkapp.com/article/910-adding-a-custom-return-path-domain).
	ReturnPathDomain string `json:"ReturnPathDomain"`
	// The verification state of the Return-Path domain. Tells you if the Return-Path is actively being used or
	// still needs further action to be used.
	ReturnPathDomainVerified bool `json:"ReturnPathDomainVerified"`
	// The CNAME DNS record that Postmark expects to find at the ReturnPathDomain value.
	ReturnPathDomainCNAMEValue string `json:"ReturnPathDomainCNAMEValue"`
	// Unique ID of sender signature.
	ID int64 `json:"ID"`
	// The text of the personal note sent to the recipient.
	ConfirmationPersonalNote string `json:"ConfirmationPersonalNote"`
}

// SenderSignatureCreateRequest is the request body for creating a new sender signature
type SenderSignatureCreateRequest struct {
	// From email associated with sender signature.
	FromEmail string `json:"FromEmail" binding:"required"`
	// From name associated with sender signature.
	Name string `json:"Name" binding:"required"`
	// Override for reply-to address.
	ReplyToEmail string `json:"ReplyToEmail"`
	// A custom value for the Return-Path domain. It is an optional field, but it must be a subdomain of your From
	// Email domain and must have a CNAME record that points to pm.mtasv.net. For more information about this field,
	// please [read our support page](http://support.postmarkapp.com/article/910-adding-a-custom-return-path-domain).
	ReturnPathDomain string `json:"ReturnPathDomain"`
	// Optional. A way to provide a note to the recipient of the confirmation email to have context of what Postmark is.
	// Max length of 400 characters.
	ConfirmationPersonalNote string `json:"ConfirmationPersonalNote"`
}

// SenderSignatureEditRequest is the request body for editing an existing sender signature
type SenderSignatureEditRequest struct {
	// From name associated with sender signature.
	Name string `json:"Name" binding:"required"`
	// Override for reply-to address.
	ReplyToEmail string `json:"ReplyToEmail"`
	// A custom value for the Return-Path domain. It is an optional field, but it must be a subdomain of your From
	// Email domain and must have a CNAME record that points to pm.mtasv.net. For more information about this field,
	// please [read our support page](http://support.postmarkapp.com/article/910-adding-a-custom-return-path-domain).
	ReturnPathDomain string `json:"ReturnPathDomain"`
	// Optional. A way to provide a note to the recipient of the confirmation email to have context of what Postmark is.
	// Max length of 400 characters.
	ConfirmationPersonalNote string `json:"ConfirmationPersonalNote"`
}

// SenderSignaturesList is just a list of SenderSignatures as they are in the response
type SenderSignaturesList struct {
	TotalCount       int
	SenderSignatures []SenderSignature
}

// GetSenderSignatures gets a list of sender signatures containing brief details associated with your account,
// limited by count and paged by offset
func (client *Client) GetSenderSignatures(ctx context.Context, count, offset int64) (SenderSignaturesList, error) {
	res := SenderSignaturesList{}

	values := &url.Values{}
	values.Add("count", fmt.Sprintf("%d", count))
	values.Add("offset", fmt.Sprintf("%d", offset))

	err := client.doRequest(ctx, parameters{
		Method:    "GET",
		Path:      fmt.Sprintf("senders?%s", values.Encode()),
		TokenType: accountToken,
	}, &res)
	return res, err
}

// GetSenderSignature gets all the details for a specific sender signature.
func (client *Client) GetSenderSignature(ctx context.Context, signatureID int) (SenderSignatureDetails, error) {
	var res SenderSignatureDetails
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodGet,
		Path:      fmt.Sprintf("senders/%d", signatureID),
		TokenType: accountToken,
	}, &res)
	return res, err
}

// CreateSenderSignature creates a new sender signature and returns the full details of the new sender signature.
func (client *Client) CreateSenderSignature(ctx context.Context, request SenderSignatureCreateRequest) (SenderSignatureDetails, error) {
	var res SenderSignatureDetails
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPost,
		Path:      "senders",
		Payload:   request,
		TokenType: accountToken,
	}, &res)
	return res, err
}

// EditSenderSignature updates an existing sender signature and returns the full details of the updated sender signature.
func (client *Client) EditSenderSignature(ctx context.Context, signatureID int, request SenderSignatureEditRequest) (SenderSignatureDetails, error) {
	var res SenderSignatureDetails
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPut,
		Path:      fmt.Sprintf("senders/%d", signatureID),
		Payload:   request,
		TokenType: accountToken,
	}, &res)
	return res, err
}

// DeleteSenderSignature removes a sender from the server.
func (client *Client) DeleteSenderSignature(ctx context.Context, signatureID int) error {
	res := APIError{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodDelete,
		Path:      fmt.Sprintf("senders/%d", signatureID),
		TokenType: accountToken,
	}, &res)

	if res.ErrorCode != 0 {
		return res
	}

	return err
}

// ResendSenderSignatureConfirmation resends the confirmation email for a sender signature.
func (client *Client) ResendSenderSignatureConfirmation(ctx context.Context, signatureID int) error {
	res := APIError{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPost,
		Path:      fmt.Sprintf("senders/%d/resend", signatureID),
		TokenType: accountToken,
	}, &res)

	if res.ErrorCode != 0 {
		return res
	}

	return err
}
