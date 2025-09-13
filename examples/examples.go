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

	// Example 10: Message Streams API examples
	demonstrateMessageStreamsAPI(client)

	// Example 11: Domains API examples
	demonstrateDomainsAPI(client)

	// Example 12: Bounce API examples
	demonstrateBounceAPI(client)

	// Example 13: Stats API examples
	demonstrateStatsAPI(client)

	// Example 14: Inbound Rules Triggers API examples
	demonstrateInboundRulesTriggersAPI(client)
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

// demonstrateMessageStreamsAPI shows examples of using the Message Streams API
func demonstrateMessageStreamsAPI(client *postmark.Client) {
	// List all message streams
	messageStreams, err := client.ListMessageStreams(context.Background(), "All", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d message streams", len(messageStreams))
	for _, stream := range messageStreams {
		log.Printf("Stream: %s (%s) - Type: %s", stream.Name, stream.ID, stream.MessageStreamType)
	}

	// List only broadcast streams including archived ones
	broadcastStreams, err := client.ListMessageStreams(context.Background(), "Broadcasts", true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d broadcast streams (including archived)", len(broadcastStreams))

	// Create a new broadcast message stream
	description := "Marketing and newsletter emails"
	createRequest := postmark.CreateMessageStreamRequest{
		ID:                "marketing-broadcasts",
		Name:              "Marketing Broadcasts",
		Description:       &description,
		MessageStreamType: postmark.BroadcastMessageStreamType,
		SubscriptionManagementConfiguration: postmark.MessageStreamSubscriptionManagementConfiguration{
			UnsubscribeHandlingType: postmark.PostmarkUnsubscribeHandlingType,
		},
	}

	newStream, err := client.CreateMessageStream(context.Background(), createRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created new stream: %s (ID: %s)", newStream.Name, newStream.ID)

	// Get details of a specific message stream
	streamDetails, err := client.GetMessageStream(context.Background(), newStream.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Stream details - Name: %s, Type: %s, Created: %s",
		streamDetails.Name, streamDetails.MessageStreamType, streamDetails.CreatedAt)

	// Edit the message stream
	newDescription := "Updated marketing and promotional emails"
	editRequest := postmark.EditMessageStreamRequest{
		Name:        "Updated Marketing Broadcasts",
		Description: &newDescription,
		SubscriptionManagementConfiguration: postmark.MessageStreamSubscriptionManagementConfiguration{
			UnsubscribeHandlingType: postmark.PostmarkUnsubscribeHandlingType,
		},
	}

	updatedStream, err := client.EditMessageStream(context.Background(), newStream.ID, editRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Updated stream: %s", updatedStream.Name)

	// Archive the message stream
	archiveResponse, err := client.ArchiveMessageStream(context.Background(), newStream.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Archived stream %s, will be purged on: %s",
		archiveResponse.ID, archiveResponse.ExpectedPurgeDate)

	// Unarchive the message stream
	unarchivedStream, err := client.UnarchiveMessageStream(context.Background(), newStream.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Unarchived stream: %s (Archived: %v)",
		unarchivedStream.Name, unarchivedStream.ArchivedAt != nil)
}

// demonstrateDomainsAPI shows examples of using the Domains API
func demonstrateDomainsAPI(client *postmark.Client) {
	// List all domains
	domains, err := client.GetDomains(context.Background(), 50, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d domains out of %d total", len(domains.Domains), domains.TotalCount)
	for _, domain := range domains.Domains {
		log.Printf("Domain: %s (ID: %d) - DKIM: %v, ReturnPath: %v",
			domain.Name, domain.ID, domain.DKIMVerified, domain.ReturnPathDomainVerified)
	}

	// Create a new domain
	createRequest := postmark.DomainCreateRequest{
		Name:             "example.com",
		ReturnPathDomain: "bounces.example.com",
	}

	newDomain, err := client.CreateDomain(context.Background(), createRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created domain: %s (ID: %d)", newDomain.Name, newDomain.ID)

	// Get detailed information about the domain
	domainDetails, err := client.GetDomain(context.Background(), newDomain.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Domain details - Name: %s, DKIM Host: %s, Return Path: %s",
		domainDetails.Name, domainDetails.DKIMHost, domainDetails.ReturnPathDomain)

	// Edit the domain (update return path)
	editRequest := postmark.DomainEditRequest{
		ReturnPathDomain: "pm-bounces.example.com",
	}

	updatedDomain, err := client.EditDomain(context.Background(), newDomain.ID, editRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Updated domain return path to: %s", updatedDomain.ReturnPathDomain)

	// Verify DKIM status
	dkimVerified, err := client.VerifyDKIMStatus(context.Background(), newDomain.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DKIM verification for %s: %v (Status: %s)",
		dkimVerified.Name, dkimVerified.DKIMVerified, dkimVerified.DKIMUpdateStatus)

	// Verify Return-Path
	returnPathVerified, err := client.VerifyReturnPath(context.Background(), newDomain.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Return Path verification for %s: %v (CNAME: %s)",
		returnPathVerified.Name, returnPathVerified.ReturnPathDomainVerified, returnPathVerified.ReturnPathDomainCNAMEValue)

	// Rotate DKIM keys (if needed for security)
	rotatedDKIM, err := client.RotateDKIM(context.Background(), newDomain.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DKIM rotation initiated for %s - Pending Host: %s",
		rotatedDKIM.Name, rotatedDKIM.DKIMPendingHost)

	// Delete the domain (cleanup example)
	err = client.DeleteDomain(context.Background(), newDomain.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleted domain: %s", newDomain.Name)

	// Example: Sender Signatures API
	log.Println("\n=== Sender Signatures API Examples ===")

	// List sender signatures
	senderSignatures, err := client.GetSenderSignatures(context.Background(), 50, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Total sender signatures: %d", senderSignatures.TotalCount)
	for _, signature := range senderSignatures.SenderSignatures {
		log.Printf("Signature: %s <%s> (ID: %d) - Confirmed: %v",
			signature.Name, signature.FromEmail, signature.ID, signature.Confirmed)
	}

	// Create a new sender signature
	createSignatureRequest := postmark.SenderSignatureCreateRequest{
		FromEmail:                "noreply@example.com",
		Name:                     "Example Service",
		ReplyToEmail:             "support@example.com",
		ReturnPathDomain:         "pm-bounces.example.com",
		ConfirmationPersonalNote: "This is a sender signature for Example Service notifications.",
	}

	newSignature, err := client.CreateSenderSignature(context.Background(), createSignatureRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created sender signature: %s <%s> (ID: %d)",
		newSignature.Name, newSignature.FromEmail, newSignature.ID)

	// Get detailed information about the sender signature
	signatureDetails, err := client.GetSenderSignature(context.Background(), newSignature.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Signature details - Name: %s, Domain: %s, DKIM Status: %s, Confirmed: %v",
		signatureDetails.Name, signatureDetails.Domain, signatureDetails.DKIMUpdateStatus, signatureDetails.Confirmed)

	// Edit the sender signature
	editSignatureRequest := postmark.SenderSignatureEditRequest{
		Name:                     "Updated Example Service",
		ReplyToEmail:             "help@example.com",
		ReturnPathDomain:         "bounces.example.com",
		ConfirmationPersonalNote: "Updated sender signature for Example Service.",
	}

	updatedSignature, err := client.EditSenderSignature(context.Background(), newSignature.ID, editSignatureRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Updated sender signature name to: %s", updatedSignature.Name)

	// Resend confirmation email (if not confirmed)
	if !updatedSignature.Confirmed {
		err = client.ResendSenderSignatureConfirmation(context.Background(), newSignature.ID)
		if err != nil {
			log.Printf("Failed to resend confirmation: %v", err)
		} else {
			log.Printf("Resent confirmation email for signature: %s", updatedSignature.FromEmail)
		}
	}

	// Delete the sender signature (cleanup example)
	err = client.DeleteSenderSignature(context.Background(), newSignature.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleted sender signature: %s", updatedSignature.FromEmail)
}

func demonstrateStatsAPI(client *postmark.Client) {
	ctx := context.Background()
	options := map[string]interface{}{
		"fromdate": "2024-01-01",
		"todate":   "2024-12-31",
		"tag":      "newsletter",
	}

	demonstrateBasicStatsAPI(ctx, client, options)
	demonstrateClickStatsAPI(ctx, client, options)
}

func demonstrateBasicStatsAPI(ctx context.Context, client *postmark.Client, options map[string]interface{}) {
	if outboundStats, err := client.GetOutboundStats(ctx, options); err != nil {
		log.Printf("Error getting outbound stats: %v", err)
	} else {
		log.Printf("Outbound Stats - Sent: %d, Bounced: %d, Opens: %d, Unique Opens: %d",
			outboundStats.Sent, outboundStats.Bounced, outboundStats.Opens, outboundStats.UniqueOpens)
	}

	if sentCounts, err := client.GetSentCounts(ctx, options); err != nil {
		log.Printf("Error getting sent counts: %v", err)
	} else {
		log.Printf("Total Sent: %d, Days with data: %d", sentCounts.Sent, len(sentCounts.Days))
	}

	if bounceCounts, err := client.GetBounceCounts(ctx, options); err != nil {
		log.Printf("Error getting bounce counts: %v", err)
	} else {
		log.Printf("Hard Bounces: %d, Soft Bounces: %d, SMTP Errors: %d",
			bounceCounts.HardBounce, bounceCounts.SoftBounce, bounceCounts.SMTPApiError)
	}

	if spamCounts, err := client.GetSpamCounts(ctx, options); err != nil {
		log.Printf("Error getting spam counts: %v", err)
	} else {
		log.Printf("Total Spam Complaints: %d", spamCounts.SpamComplaint)
	}

	if trackedCounts, err := client.GetTrackedCounts(ctx, options); err != nil {
		log.Printf("Error getting tracked counts: %v", err)
	} else {
		log.Printf("Total Tracked Emails: %d", trackedCounts.Tracked)
	}

	if openCounts, err := client.GetOpenCounts(ctx, options); err != nil {
		log.Printf("Error getting open counts: %v", err)
	} else {
		log.Printf("Total Opens: %d, Unique Opens: %d", openCounts.Opens, openCounts.Unique)
	}

	if platformCounts, err := client.GetPlatformCounts(ctx, options); err != nil {
		log.Printf("Error getting platform counts: %v", err)
	} else {
		log.Printf("Opens by Platform - Desktop: %d, Mobile: %d, WebMail: %d, Unknown: %d",
			platformCounts.Desktop, platformCounts.Mobile, platformCounts.WebMail, platformCounts.Unknown)
	}

	if emailClientCounts, err := client.GetEmailClientCounts(ctx, options); err != nil {
		log.Printf("Error getting email client counts: %v", err)
	} else {
		log.Printf("Opens by Email Client - Outlook: %d, Gmail: %d, AppleMail: %d, Yahoo: %d",
			emailClientCounts.Outlook, emailClientCounts.Gmail, emailClientCounts.AppleMail, emailClientCounts.Yahoo)
	}
}

func demonstrateClickStatsAPI(ctx context.Context, client *postmark.Client, options map[string]interface{}) {
	if clickCounts, err := client.GetClickCounts(ctx, options); err != nil {
		log.Printf("Error getting click counts: %v", err)
	} else {
		log.Printf("Total Clicks: %d, Unique Clicks: %d", clickCounts.Clicks, clickCounts.Unique)
	}

	if browserCounts, err := client.GetBrowserFamilyCounts(ctx, options); err != nil {
		log.Printf("Error getting browser family counts: %v", err)
	} else {
		log.Printf("Clicks by Browser - Chrome: %d, Safari: %d, Firefox: %d, IE: %d",
			browserCounts.Chrome, browserCounts.Safari, browserCounts.Firefox, browserCounts.InternetExplorer)
	}

	if locationCounts, err := client.GetClickLocationCounts(ctx, options); err != nil {
		log.Printf("Error getting click location counts: %v", err)
	} else {
		log.Printf("Clicks by Location - HTML: %d, Text: %d", locationCounts.HTML, locationCounts.Text)
	}

	if clickPlatformCounts, err := client.GetClickPlatformCounts(ctx, options); err != nil {
		log.Printf("Error getting click platform counts: %v", err)
	} else {
		log.Printf("Clicks by Platform - Desktop: %d, Mobile: %d, WebMail: %d, Unknown: %d",
			clickPlatformCounts.Desktop, clickPlatformCounts.Mobile, clickPlatformCounts.WebMail, clickPlatformCounts.Unknown)
	}
}

// demonstrateInboundRulesTriggersAPI shows examples of using the Inbound Rules Triggers API
func demonstrateInboundRulesTriggersAPI(client *postmark.Client) {
	ctx := context.Background()

	// List existing inbound rule triggers
	triggers, totalCount, err := client.GetInboundRuleTriggers(ctx, 50, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d inbound rule triggers out of %d total", len(triggers), totalCount)
	for _, trigger := range triggers {
		log.Printf("Trigger ID %d: %s", trigger.ID, trigger.Rule)
	}

	// Create inbound rule trigger for blocking specific email
	emailTrigger, err := client.CreateInboundRuleTrigger(ctx, "spam@example.com")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created email trigger: ID %d blocking '%s'", emailTrigger.ID, emailTrigger.Rule)

	// Create inbound rule trigger for blocking entire domain
	domainTrigger, err := client.CreateInboundRuleTrigger(ctx, "*.spammer.com")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created domain trigger: ID %d blocking '%s'", domainTrigger.ID, domainTrigger.Rule)

	// Delete the created triggers (cleanup)
	err = client.DeleteInboundRuleTrigger(ctx, emailTrigger.ID)
	if err != nil {
		log.Printf("Failed to delete email trigger: %v", err)
	} else {
		log.Printf("Deleted email trigger: ID %d", emailTrigger.ID)
	}

	err = client.DeleteInboundRuleTrigger(ctx, domainTrigger.ID)
	if err != nil {
		log.Printf("Failed to delete domain trigger: %v", err)
	} else {
		log.Printf("Deleted domain trigger: ID %d", domainTrigger.ID)
	}

	// List triggers again to confirm deletion
	triggersAfter, totalCountAfter, err := client.GetInboundRuleTriggers(ctx, 50, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("After cleanup: %d inbound rule triggers out of %d total", len(triggersAfter), totalCountAfter)
}
