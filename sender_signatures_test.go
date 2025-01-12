package postmark

import (
	"context"
	"net/http"
	"testing"

	"goji.io/pat"
)

func TestGetSenderSignatures(t *testing.T) {
	responseJSON := `{
	"TotalCount": 2,
	"SenderSignatures": [
	  {
		"Domain": "wildbit.com",
		"EmailAddress": "jp@wildbit.com",
		"ReplyToEmailAddress": "info@wildbit.com",
		"Name": "JP Toto",
		"Confirmed": true,
		"ID": 36735
	  },
	  {
		"Domain": "example.com",
		"EmailAddress": "jp@example.com",
		"ReplyToEmailAddress": "",
		"Name": "JP Toto",
		"Confirmed": true,
		"ID": 81605
	  }
	]
  }`

	tMux.HandleFunc(pat.Get("/senders"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := client.GetSenderSignatures(context.Background(), 50, 0)
	if err != nil {
		t.Fatalf("GetSenderSignatures: %s", err.Error())
	}

	if res.TotalCount != 2 {
		t.Fatalf("GetSenderSignatures: wrong TotalCount!")
	}
}

func TestGetSenderSignature(t *testing.T) {
	responseJSON := `{
  "Domain": "postmarkapp.com",
  "EmailAddress": "jp@postmarkapp.com",
  "ReplyToEmailAddress": "info@postmarkapp.com",
  "Name": "JP Toto",
  "Confirmed": true,
  "SPFVerified": true,
  "SPFHost": "postmarkapp.com",
  "SPFTextValue": "v=spf1 a mx include:spf.mtasv.net ~all",
  "DKIMVerified": false,
  "WeakDKIM": false,
  "DKIMHost": "jan2013.pm._domainkey.postmarkapp.com",
  "DKIMTextValue": "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDJ...",
  "DKIMPendingHost": "20131031155228.pm._domainkey.postmarkapp.com",
  "DKIMPendingTextValue": "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCFn...",
  "DKIMRevokedHost": "",
  "DKIMRevokedTextValue": "",
  "SafeToRemoveRevokedKeyFromDNS": false,
  "DKIMUpdateStatus": "Pending",
  "ReturnPathDomain": "pm-bounces.postmarkapp.com",
  "ReturnPathDomainVerified": false,
  "ReturnPathDomainCNAMEValue": "pm.mtasv.net",
  "ID": 1234,
  "ConfirmationPersonalNote": "This is a note visible to the recipient to provide context of what Postmark is."
}`

	tMux.HandleFunc(pat.Get("/senders/:signatureID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := client.GetSenderSignature(context.Background(), 1234)
	if err != nil {
		t.Fatalf("ServerSignature: %s", err.Error())
	}

	if res.Name != "JP Toto" {
		t.Fatalf("ServerSignature: wrong name!")
	}
}

func TestCreateSenderSignature(t *testing.T) {
	responseJSON := `{
  "Domain": "example.com",
  "EmailAddress": "john.doe@example.com",
  "ReplyToEmailAddress": "reply@example.com",
  "Name": "John Doe",
  "Confirmed": false,
  "SPFVerified": false,
  "SPFHost": "example.com",
  "SPFTextValue": "v=spf1 a mx include:spf.mtasv.net ~all",
  "DKIMVerified": false,
  "WeakDKIM": false,
  "DKIMHost": "",
  "DKIMTextValue": "",
  "DKIMPendingHost": "20140220130148.pm._domainkey.example.com",
  "DKIMPendingTextValue": "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCQ35xZciGB0g...",
  "DKIMRevokedHost": "",
  "DKIMRevokedTextValue": "",
  "SafeToRemoveRevokedKeyFromDNS": false,
  "DKIMUpdateStatus": "Pending",
  "ReturnPathDomain": "pm-bounces.example.com",
  "ReturnPathDomainVerified": true,
  "ReturnPathDomainCNAMEValue": "pm.mtasv.net",
  "ID": 1,
  "ConfirmationPersonalNote": "This is a note visible to the recipient to provide context of what Postmark is."
}`

	tMux.HandleFunc(pat.Post("/senders"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := client.CreateSenderSignature(context.Background(), SenderSignatureCreateRequest{
		FromEmail:                "john.doe@example.com",
		Name:                     "John Doe",
		ReplyToEmail:             "reply@example.com",
		ReturnPathDomain:         "pm-bounces.example.com",
		ConfirmationPersonalNote: "This is a note visible to the recipient to provide context of what Postmark is.",
	})
	if err != nil {
		t.Fatalf("CreateSenderSignature: %s", err.Error())
	}

	if res.Name != "John Doe" {
		t.Fatalf("CreateSenderSignature: wrong name!")
	}
}

func TestEditSenderSignature(t *testing.T) {
	responseJSON := `{
  "Domain": "example.com",
  "EmailAddress": "john.doe@example.com",
  "ReplyToEmailAddress": "jane.doe@example.com",
  "Name": "Jane Doe",
  "Confirmed": false,
  "SPFVerified": false,
  "SPFHost": "crazydomain.com",
  "SPFTextValue": "v=spf1 a mx include:spf.mtasv.net ~all",
  "DKIMVerified": false,
  "WeakDKIM": false,
  "DKIMHost": "",
  "DKIMTextValue": "",
  "DKIMPendingHost": "20140220130148.pm._domainkey.crazydomain.com",
  "DKIMPendingTextValue": "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCQ35xZciGB0g...",
  "DKIMRevokedHost": "",
  "DKIMRevokedTextValue": "",
  "SafeToRemoveRevokedKeyFromDNS": false,
  "DKIMUpdateStatus": "Pending",
  "ReturnPathDomain": "pm-bounces.example.com",
  "ReturnPathDomainVerified": true,
  "ReturnPathDomainCNAMEValue": "pm.mtasv.net",
  "ID": 1,
  "ConfirmationPersonalNote": "This is a note visible to the recipient to provide context of what Postmark is."
}`

	tMux.HandleFunc(pat.Put("/senders/:signatureID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := client.EditSenderSignature(context.Background(), 1, SenderSignatureEditRequest{
		Name:                     "Jane Doe",
		ReplyToEmail:             "jane.doe@example.com",
		ReturnPathDomain:         "pm-bounces.example.com",
		ConfirmationPersonalNote: "This is a note visible to the recipient to provide context of what Postmark is.",
	})
	if err != nil {
		t.Fatalf("EditSenderSignature: %s", err.Error())
	}

	if res.Name != "Jane Doe" {
		t.Fatalf("EditSenderSignature: wrong name!")
	}
}

func TestDeleteSenderSignature(t *testing.T) {
	responseJSON := `{
	  "ErrorCode": 0,
	  "Message": "SenderSignature 1234 removed."
	}`

	tMux.HandleFunc(pat.Delete("/senders/:signatureID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	// Success
	err := client.DeleteSenderSignature(context.Background(), 1234)
	if err != nil {
		t.Fatalf("DeleteSenderSignature: %s", err.Error())
	}

	// Failure
	responseJSON = `{
	  "ErrorCode": 402,
	  "Message": "Invalid JSON"
	}`

	err = client.DeleteSenderSignature(context.Background(), 1234)
	if err == nil {
		t.Fatalf("DeleteSenderSignature  should have failed")
	}
}
