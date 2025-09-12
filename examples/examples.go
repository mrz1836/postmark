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
}
