// Package main is an example of how to use the Postmark Go client.
package main

import (
	"context"
	"log"

	"github.com/mrz1836/postmark"
)

func main() {
	client := postmark.NewClient("[SERVER-TOKEN]", "[ACCOUNT-TOKEN]")

	// Example 1: Send email with InlineCss feature
	email := postmark.Email{
		From:       "no-reply@example.com",
		To:         "tito@example.com",
		Subject:    "Reset your password",
		HTMLBody:   "<style>body { color: blue; }</style><body>Your password reset link</body>",
		TextBody:   "Your password reset link",
		Tag:        "pw-reset",
		TrackOpens: true,
		InlineCSS:  true, // This will inline CSS styles into HTML attributes
	}

	res, err := client.SendEmail(context.Background(), email)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Email sent with ID: %s", res.MessageID)

	// Example 2: Create a data removal request
	dataRemovalRequest := postmark.DataRemovalRequest{
		Recipient: "user@example.com",
	}

	removalResponse, err := client.CreateDataRemoval(context.Background(), dataRemovalRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Data removal request created with ID: %d", removalResponse.ID)

	// Example 3: Check data removal status
	status, err := client.GetDataRemovalStatus(context.Background(), removalResponse.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Data removal status: %s", status.Status)

	// Example 4: Create a template with alias
	template := postmark.Template{
		Name:           "Welcome Email",
		Alias:          "welcome-template",
		TemplateType:   "Standard",
		Subject:        "Welcome to {{company_name}}, {{user_name}}!",
		HTMLBody:       "<html><body><h1>Welcome {{user_name}}!</h1><p>Thanks for joining {{company_name}}.</p></body></html>",
		TextBody:       "Welcome {{user_name}}!\n\nThanks for joining {{company_name}}.",
		LayoutTemplate: "", // No layout template for this example
	}

	templateResponse, err := client.CreateTemplate(context.Background(), template)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Template created with ID: %d, Alias: %s", templateResponse.TemplateID, templateResponse.Alias)

	// Example 5: Send email using template alias
	templatedEmail := postmark.TemplatedEmail{
		TemplateAlias: "welcome-template", // Using alias instead of ID
		TemplateModel: map[string]interface{}{
			"user_name":    "John Doe",
			"company_name": "ACME Corporation",
		},
		From: "no-reply@example.com",
		To:   "john.doe@example.com",
	}

	emailResponse, err := client.SendTemplatedEmail(context.Background(), templatedEmail)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Templated email sent with ID: %s", emailResponse.MessageID)

	// Example 6: Get filtered templates (only Layout templates)
	templates, count, err := client.GetTemplatesFiltered(context.Background(), 50, 0, "Layout", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d layout templates out of %d total", len(templates), count)

	// Example 7: Push templates between servers (requires account token)
	pushRequest := postmark.PushTemplatesRequest{
		SourceServerID:      1001, // Replace with actual server IDs
		DestinationServerID: 1002,
		PerformChanges:      false, // Set to true to actually perform the push
	}

	pushResponse, err := client.PushTemplates(context.Background(), pushRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Push simulation completed: %d templates would be affected", pushResponse.TotalCount)
	for _, pushedTemplate := range pushResponse.Templates {
		log.Printf("Template '%s' would be %s", pushedTemplate.Name, pushedTemplate.Action)
	}
}
