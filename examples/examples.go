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

	// Example 6: Get message clicks with filtering
	clicks, clickCount, err := client.GetOutboundMessagesClicks(context.Background(), 100, 0, map[string]interface{}{
		"tag":           "welcome-template",
		"recipient":     "john.doe@example.com",
		"client_name":   "Chrome",
		"platform":      "Desktop",
		"messagestream": "outbound",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d clicks out of %d total", len(clicks), clickCount)
	for _, click := range clicks {
		log.Printf("Click: %s clicked %s using %s", click.Recipient, click.OriginalLink, click.Client["Name"])
	}

	// Example 7: Get clicks for a specific message
	messageClicks, messageClickCount, err := client.GetOutboundMessageClicks(context.Background(), emailResponse.MessageID, 50, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Message %s has %d clicks", emailResponse.MessageID, messageClickCount)
	for _, click := range messageClicks {
		log.Printf("Click at %s: %s from %s", click.ReceivedAt.Format("2006-01-02 15:04:05"), click.OriginalLink, click.Geo["City"])
	}

	// Example 8: Get filtered templates (only Layout templates)
	templates, count, err := client.GetTemplatesFiltered(context.Background(), 50, 0, "Layout", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d layout templates out of %d total", len(templates), count)

	// Example 9: Push templates between servers (requires account token)
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

	// Example 8: Bounce API examples
	demonstrateBounceAPI(client)
}

// demonstrateBounceAPI shows examples of using the Bounce API
func demonstrateBounceAPI(client *postmark.Client) {
	// Get delivery stats to understand bounce metrics
	deliveryStats, err := client.GetDeliveryStats(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Inactive emails: %d", deliveryStats.InactiveMails)
	for _, bounceType := range deliveryStats.Bounces {
		log.Printf("Bounce type '%s': %d bounces", bounceType.Name, bounceType.Count)
	}

	// Get bounces with filtering options
	bounces, totalCount, err := client.GetBounces(context.Background(), 50, 0, map[string]interface{}{
		"type":        "HardBounce",
		"inactive":    true,
		"emailFilter": "@example.com",
		"tag":         "password-reset",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d bounces out of %d total", len(bounces), totalCount)
	for _, bounce := range bounces {
		log.Printf("Bounce ID %d: %s (%s) - %s", bounce.ID, bounce.Email, bounce.Type, bounce.Description)
	}

	// Process specific bounce details if available
	if len(bounces) > 0 {
		processBounceDetails(client, bounces[0].ID)
	}

	// Get tags that have generated bounces
	bouncedTags, err := client.GetBouncedTags(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Tags with bounces (%d): %v", len(bouncedTags), bouncedTags)
}

// processBounceDetails demonstrates detailed bounce operations
func processBounceDetails(client *postmark.Client, bounceID int64) {
	// Get detailed bounce information
	bounceDetails, err := client.GetBounce(context.Background(), bounceID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Bounce details for ID %d: Type=%s, MessageID=%s, Subject=%s",
		bounceDetails.ID, bounceDetails.Type, bounceDetails.MessageID, bounceDetails.Subject)

	// Get bounce dump if available
	if bounceDetails.DumpAvailable {
		dumpContent, err := client.GetBounceDump(context.Background(), bounceDetails.ID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Bounce dump for ID %d: %d characters", bounceDetails.ID, len(dumpContent))
	}

	// Activate bounce if possible
	if bounceDetails.CanActivate {
		activatedBounce, message, err := client.ActivateBounce(context.Background(), bounceDetails.ID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Bounce activation result: %s (ID: %d, Now Active: %v)",
			message, activatedBounce.ID, !activatedBounce.Inactive)
	}
}
