package postmark

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Domain contains the details of a domain in Postmark. https://postmarkapp.com/developer/api/domains-api
type Domain struct {
	// Name of the domain.
	Name string `binding:"required"`
	// Host name used for the SPF configuration.
	SPFHost string
	// Value that can be optionally setup with your DNS host. See our
	// [blog post](https://postmarkapp.com/blog/why-we-no-longer-ask-for-spf-records) to learn why this field is no
	// longer necessary.
	SPFTextValue string
	// Specifies whether DKIM has ever been verified for the domain or not. Once DKIM is verified, this response will
	// stay true, even if the record is later removed from DNS.
	DKIMVerified bool
	// DKIM is using a strength weaker than 1024 bit. If so, it’s possible to request a new DKIM using the
	// [RequestNewDKIM](https://postmarkapp.com/developer/api/domains-api#rotate-dkim) function below.
	WeakDKIM bool
	// DNS TXT host being used to validate messages sent in
	DKIMHost string
	// DNS TXT value being used to validate messages sent in.
	DKIMTextValue string
	// If a DKIM rotation has been initiated or this DKIM is from a new Domain, this field will show the pending
	// DKIM DNS TXT host which has yet to be setup and confirmed at your registrar or DNS host.
	DKIMPendingHost string
	// Similar to the DKIMPendingHost field, this will show the DNS TXT value waiting to be confirmed at your
	// registrar or DNS host.
	DKIMPendingTextValue string
	// Once a new DKIM has been confirmed at your registrar or DNS host, Postmark will revoke the old DKIM host in
	// preparation for removing it permanently from the system.
	DKIMRevokedHost string
	// Similar to DKIMRevokedHost, this field will show the DNS TXT value that will soon be removed from the Postmark system.
	DKIMRevokedTextValue string
	// Indicates whether you may safely delete the old DKIM DNS TXT records at your registrar or DNS host. The new
	// DKIM is now safely in use.
	SafeToRemoveRevokedKeyFromDNS bool
	// While DKIM renewal or new DKIM operations are being conducted or setup, this field will indicate Pending.
	// After all DNS TXT records are up to date and any pending renewal operations are finished, it will indicate Verified.
	DKIMUpdateStatus string
	// The custom Return-Path for this domain, please [read our support page](http://support.postmarkapp.com/article/910-adding-a-custom-return-path-domain).
	ReturnPathDomain string
	// The verification state of the Return-Path domain. Tells you if the Return-Path is actively being used or still
	// needs further action to be used.
	ReturnPathDomainVerified bool
	// The CNAME DNS record that Postmark expects to find at the ReturnPathDomain value.
	ReturnPathDomainCNAMEValue string
	// Unique ID of the Domain.
	ID int64
}

// DomainCreateRequest is the request body to create a domain
type DomainCreateRequest struct {
	// Name of the domain.
	Name string `binding:"required"`
	// A custom value for the Return-Path domain. It is an optional field, but it must be a subdomain of your
	// From Email domain and must have a CNAME record that points to pm.mtasv.net. For more information about this
	// field, please read our [support page](http://support.postmarkapp.com/article/910-adding-a-custom-return-path-domain).
	ReturnPathDomain string
}

// DomainEditRequest is the request body to edit a domain
type DomainEditRequest struct {
	// A custom value for the Return-Path domain. It is an optional field, but it must be a subdomain of your
	// From Email domain and must have a CNAME record that points to pm.mtasv.net. For more information about this
	// field, please read our [support page](http://support.postmarkapp.com/article/910-adding-a-custom-return-path-domain).
	ReturnPathDomain string
}

// DomainsList is just a list of Domains as they are in the response
type DomainsList struct {
	TotalCount int
	Domains    []Domain
}

// GetDomains gets a list of domains, limited by count and paged by offset
func (client *Client) GetDomains(ctx context.Context, count, offset int64) (DomainsList, error) {
	res := DomainsList{}

	values := &url.Values{}
	values.Add("count", fmt.Sprintf("%d", count))
	values.Add("offset", fmt.Sprintf("%d", offset))

	err := client.doRequest(ctx, parameters{
		Method:    "GET",
		Path:      fmt.Sprintf("domains?%s", values.Encode()),
		TokenType: accountToken,
	}, &res)
	return res, err
}

// GetDomain fetches a specific domain via domainID
func (client *Client) GetDomain(ctx context.Context, domainID string) (Domain, error) {
	res := Domain{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodGet,
		Path:      fmt.Sprintf("domains/%s", domainID),
		TokenType: accountToken,
	}, &res)
	return res, err
}

// EditDomain updates details for a specific domain with domainID
func (client *Client) EditDomain(ctx context.Context, domainID string, request DomainEditRequest) (Domain, error) {
	res := Domain{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPut,
		Path:      fmt.Sprintf("domains/%s", domainID),
		TokenType: accountToken,
		Payload:   request,
	}, &res)
	return res, err
}

// CreateDomain creates a domain
func (client *Client) CreateDomain(ctx context.Context, request DomainCreateRequest) (Domain, error) {
	res := Domain{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPost,
		Path:      "domains",
		TokenType: accountToken,
		Payload:   request,
	}, &res)
	return res, err
}

// DeleteDomain deletes a specific domain via domainID
func (client *Client) DeleteDomain(ctx context.Context, domainID string) error {
	res := APIError{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodDelete,
		Path:      fmt.Sprintf("domains/%s", domainID),
		TokenType: accountToken,
	}, &res)

	if res.ErrorCode != 0 {
		return res
	}

	return err
}

// VerifyDKIMStatus verifies DKIM keys for the specified domain.
func (client *Client) VerifyDKIMStatus(ctx context.Context, domainID string) (Domain, error) {
	res := Domain{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPut,
		Path:      fmt.Sprintf("domains/%s/verifyDkim", domainID),
		TokenType: accountToken,
	}, &res)
	return res, err
}

// VerifyReturnPath verifies Return-Path DNS record for the specified domain.
func (client *Client) VerifyReturnPath(ctx context.Context, domainID string) (Domain, error) {
	res := Domain{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPut,
		Path:      fmt.Sprintf("domains/%s/verifyReturnPath", domainID),
		TokenType: accountToken,
	}, &res)
	return res, err
}

// RotateDKIM Creates a new DKIM key to replace your current key. Until the new DNS entries are confirmed, the pending
// values will be in DKIMPendingHost and DKIMPendingTextValue fields. After the new DKIM value is verified in DNS,
// the pending values will migrate to DKIMTextValue and DKIMPendingTextValue and Postmark will begin to sign emails with
// the new DKIM key.
func (client *Client) RotateDKIM(ctx context.Context, domainID string) (Domain, error) {
	res := Domain{}
	err := client.doRequest(ctx, parameters{
		Method:    http.MethodPost,
		Path:      fmt.Sprintf("domains/%s/rotatedkim", domainID),
		TokenType: accountToken,
	}, &res)
	return res, err
}
