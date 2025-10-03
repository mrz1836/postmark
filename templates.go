package postmark

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ErrHeaderInjection is returned when header injection is detected
var ErrHeaderInjection = errors.New("header injection detected: illegal characters in template alias")

// validateTemplateAlias checks for header injection attempts in template alias
func validateTemplateAlias(alias string) error {
	if strings.Contains(alias, "\r") || strings.Contains(alias, "\n") {
		return ErrHeaderInjection
	}
	return nil
}

// Template represents an email template on the server
type Template struct {
	// TemplateID: ID of template
	TemplateID int64 `json:"TemplateID"`
	// Name: Name of template
	Name string
	// Subject: The content to use for the Subject when this template is used to send email.
	Subject string
	// HTMLBody: The content to use for the HTMLBody when this template is used to send email.
	HTMLBody string `json:"HtmlBody"`
	// TextBody: The content to use for the TextBody when this template is used to send email.
	TextBody string
	// AssociatedServerID: The ID of the Server with which this template is associated.
	AssociatedServerID int64 `json:"AssociatedServerId"`
	// Active: Indicates that this template may be used for sending email.
	Active bool
	// Alias: Optional alias for the template.
	Alias string `json:",omitempty"`
	// TemplateType: Type of template (Standard or Layout)
	TemplateType string `json:",omitempty"`
	// LayoutTemplate: Layout template alias if using a layout
	LayoutTemplate string `json:",omitempty"`
}

// TemplateInfo is a limited set of template info returned via Index/Editing endpoints
type TemplateInfo struct {
	// TemplateID: ID of template
	TemplateID int64 `json:"TemplateID"`
	// Name: Name of template
	Name string
	// Active: Indicates that this template may be used for sending email.
	Active bool
	// Alias: Optional alias for the template.
	Alias string `json:",omitempty"`
	// TemplateType: Type of template (Standard or Layout)
	TemplateType string `json:",omitempty"`
	// LayoutTemplate: Layout template alias if using a layout
	LayoutTemplate string `json:",omitempty"`
}

// GetTemplate fetches a specific template via TemplateID
func (client *Client) GetTemplate(ctx context.Context, templateID string) (Template, error) {
	res := Template{}
	err := client.get(ctx, fmt.Sprintf("templates/%s", templateID), &res)
	return res, err
}

type templatesResponse struct {
	TotalCount int64
	Templates  []TemplateInfo
}

// GetTemplates fetches a list of templates on the server
// It returns a TemplateInfo slice, the total template count, and any error that occurred
// TemplateInfo only returns a subset of template attributes, use GetTemplate(id) to
// retrieve all template info.
func (client *Client) GetTemplates(ctx context.Context, count, offset int64) ([]TemplateInfo, int64, error) {
	return client.GetTemplatesFiltered(ctx, count, offset, "", "")
}

// GetTemplatesFiltered fetches a filtered list of templates on the server
// templateType: filter by template type ("Standard", "Layout", or "" for all)
// layoutTemplate: filter by layout template alias (or "" for all)
func (client *Client) GetTemplatesFiltered(ctx context.Context, count, offset int64, templateType, layoutTemplate string) ([]TemplateInfo, int64, error) {
	res := templatesResponse{}

	values := &url.Values{}
	values.Add("count", fmt.Sprintf("%d", count))
	values.Add("offset", fmt.Sprintf("%d", offset))

	if templateType != "" {
		values.Add("TemplateType", templateType)
	}
	if layoutTemplate != "" {
		values.Add("LayoutTemplate", layoutTemplate)
	}

	err := client.get(ctx, buildURLWithQuery("templates", *values), &res)
	return res.Templates, res.TotalCount, err
}

// CreateTemplate saves a new template to the server
func (client *Client) CreateTemplate(ctx context.Context, template Template) (TemplateInfo, error) {
	res := TemplateInfo{}
	err := client.post(ctx, "templates", template, &res)
	return res, err
}

// EditTemplate updates details for a specific template with templateID
func (client *Client) EditTemplate(ctx context.Context, templateID string, template Template) (TemplateInfo, error) {
	res := TemplateInfo{}
	err := client.put(ctx, fmt.Sprintf("templates/%s", templateID), template, &res)
	return res, err
}

// DeleteTemplate removes a template (with templateID) from the server
func (client *Client) DeleteTemplate(ctx context.Context, templateID string) error {
	res := APIError{}
	err := client.delete(ctx, fmt.Sprintf("templates/%s", templateID), &res)
	if err != nil {
		return err
	}
	if res.ErrorCode != 0 {
		return res
	}
	return nil
}

// ValidateTemplateBody contains the template/render model combination to be validated
type ValidateTemplateBody struct {
	Subject                    string
	TextBody                   string
	HTMLBody                   string `json:"HTMLBody"`
	TestRenderModel            map[string]interface{}
	InlineCSSForHTMLTestRender bool `json:"InlineCssForHtmlTestRender"`
}

// ValidateTemplateResponse contains information as to how the validation went
type ValidateTemplateResponse struct {
	AllContentIsValid      bool
	HTMLBody               Validation `json:"HTMLBody"`
	TextBody               Validation
	Subject                Validation
	SuggestedTemplateModel map[string]interface{}
}

// Validation contains the results of a field's validation
type Validation struct {
	ContentIsValid   bool
	ValidationErrors []ValidationError
	RenderedContent  string
}

// ValidationError contains information about the errors which occurred during validation for a given field
type ValidationError struct {
	Message           string
	Line              int
	CharacterPosition int
}

// ValidateTemplate validates the provided template/render model combination
func (client *Client) ValidateTemplate(ctx context.Context, validateTemplateBody ValidateTemplateBody) (ValidateTemplateResponse, error) {
	res := ValidateTemplateResponse{}
	err := client.post(ctx, "templates/validate", validateTemplateBody, &res)
	return res, err
}

// TemplatedEmail is used to send an email via a template
type TemplatedEmail struct {
	// TemplateID: REQUIRED if TemplateAlias is not specified. - The template id to use when sending this message.
	TemplateID int64 `json:"TemplateId,omitempty"`
	// TemplateAlias: REQUIRED if TemplateID is not specified. - The template alias to use when sending this message.
	TemplateAlias string `json:",omitempty"`
	// TemplateModel: The model to be applied to the specified template to generate HtmlBody, TextBody, and Subject.
	TemplateModel map[string]interface{} `json:",omitempty"`
	// InlineCSS: By default, if the specified template contains an HtmlBody, we will apply the style blocks as inline attributes to the rendered HTML content. You may opt out of this behavior by passing false for this request field.
	InlineCSS bool `json:"InlineCSS,omitempty"`
	// From: The sender email address. Must have a registered and confirmed Sender Signature.
	From string `json:",omitempty"`
	// To: REQUIRED Recipient email address. Multiple addresses are comma separated. Max 50.
	To string `json:",omitempty"`
	// Cc recipient email address. Multiple addresses are comma separated. Max 50.
	Cc string `json:",omitempty"`
	// Bcc recipient email address. Multiple addresses are comma separated. Max 50.
	Bcc string `json:",omitempty"`
	// Tag: Email tag that allows you to categorize outgoing emails and get detailed statistics.
	Tag string `json:",omitempty"`
	// Reply To override email address. Defaults to the Reply To set in the sender signature.
	ReplyTo string `json:",omitempty"`
	// Headers: List of custom headers to include.
	Headers []Header `json:",omitempty"`
	// TrackOpens: Activate open tracking for this email.
	TrackOpens bool `json:",omitempty"`
	// TrackLinks: Activate link tracking. Possible options: "None", "HtmlAndText", "HtmlOnly", "TextOnly".
	TrackLinks string `json:",omitempty"`
	// Attachments: List of attachments
	Attachments []Attachment `json:",omitempty"`
	// MessageStream: MessageStream will default to the outbound message stream ID (Default Transactional Stream) if no message stream ID is provided.
	MessageStream string `json:",omitempty"`
	// Metadata: Custom metadata key/value pairs.
	Metadata map[string]string `json:",omitempty"`
}

// SendTemplatedEmail sends an email using a template (TemplateID)
func (client *Client) SendTemplatedEmail(ctx context.Context, email TemplatedEmail) (EmailResponse, error) {
	// Validate TemplateAlias for header injection
	if err := validateTemplateAlias(email.TemplateAlias); err != nil {
		return EmailResponse{}, err
	}

	res := EmailResponse{}
	err := client.post(ctx, "email/withTemplate", email, &res)
	return res, err
}

// SendTemplatedEmailBatch sends batch email using a template (TemplateID)
func (client *Client) SendTemplatedEmailBatch(ctx context.Context, emails []TemplatedEmail) ([]EmailResponse, error) {
	// Validate TemplateAlias for header injection in all emails
	for i, email := range emails {
		if err := validateTemplateAlias(email.TemplateAlias); err != nil {
			return nil, fmt.Errorf("email %d: %w", i, err)
		}
	}

	var res []EmailResponse
	formatEmails := map[string]interface{}{
		"Messages": emails,
	}
	err := client.post(ctx, "email/batchWithTemplates", formatEmails, &res)
	return res, err
}

// PushTemplatesRequest contains the request data for pushing templates between servers
type PushTemplatesRequest struct {
	// SourceServerID: ID of the server to push templates from
	SourceServerID int64 `json:"SourceServerId"`
	// DestinationServerID: ID of the server to push templates to
	DestinationServerID int64 `json:"DestinationServerId"`
	// PerformChanges: Whether to actually perform the push (true) or just simulate it (false)
	PerformChanges bool `json:",omitempty"`
}

// PushedTemplate represents a template that was pushed between servers
type PushedTemplate struct {
	// TemplateID: ID of the template
	TemplateID int64 `json:"TemplateId"`
	// Name: Name of the template
	Name string
	// Alias: Alias of the template (if any)
	Alias string
	// Action: Action performed (Created, Updated, Skipped, etc.)
	Action string
}

// PushTemplatesResponse contains the results of pushing templates between servers
type PushTemplatesResponse struct {
	// TotalCount: Total number of templates processed
	TotalCount int64
	// Templates: Details of each template that was processed
	Templates []PushedTemplate
}

// PushTemplates pushes templates from one server to another
func (client *Client) PushTemplates(ctx context.Context, request PushTemplatesRequest) (PushTemplatesResponse, error) {
	res := PushTemplatesResponse{}
	err := client.putWithAccountToken(ctx, "templates/push", request, &res)
	return res, err
}