package postmark

import (
	"context"
	"net/http"
	"testing"

	"goji.io/pat"
)

func TestGetDomain(t *testing.T) {
	responseJSON := `{
  "Name": "postmarkapp.com",
  "SPFVerified": true,
  "SPFHost": "postmarkapp.com",
  "SPFTextValue": "v=spf1 a mx include:spf.mtasv.net ~all",
  "DKIMVerified": false,
  "WeakDKIM": false,
  "DKIMHost": "jan2013pm._domainkey.postmarkapp.com",
  "DKIMTextValue": "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDJ...",
  "DKIMPendingHost": "20131031155228pm._domainkey.postmarkapp.com",
  "DKIMPendingTextValue": "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCFn...",
  "DKIMRevokedHost": "",
  "DKIMRevokedTextValue": "",
  "SafeToRemoveRevokedKeyFromDNS": false,
  "DKIMUpdateStatus": "Pending",
  "ReturnPathDomain": "pm-bounces.postmarkapp.com",
  "ReturnPathDomainVerified": false,
  "ReturnPathDomainCNAMEValue": "pm.mtasv.net",
  "ID": 1234
}`

	tMux.HandleFunc(pat.Get("/domains/:domainID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := client.GetDomain(context.Background(), "1234")
	if err != nil {
		t.Fatalf("GetDomain: %s", err.Error())
	}

	if res.Name != "postmarkapp.com" {
		t.Fatalf("GetDomain: wrong name!: %s", res.Name)
	}
}

func TestCreateDomain(t *testing.T) {
	responseJSON := `{
  "Name": "example.com",
  "SPFVerified": false,
  "SPFHost": "example.com",
  "SPFTextValue": "v=spf1 a mx include:spf.mtasv.net ~all",
  "DKIMVerified": false,
  "WeakDKIM": false,
  "DKIMHost": "",
  "DKIMTextValue": "",
  "DKIMPendingHost": "20131031155228pm._domainkey.example.com",
  "DKIMPendingTextValue": "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCFn...",
  "DKIMRevokedHost": "",
  "DKIMRevokedTextValue": "",
  "SafeToRemoveRevokedKeyFromDNS": false,
  "DKIMUpdateStatus": "Pending",
  "ReturnPathDomain": "pm-bounces.example.com",
  "ReturnPathDomainVerified": false,
  "ReturnPathDomainCNAMEValue": "pm.mtasv.net",
  "ID": 1234
}`

	tMux.HandleFunc(pat.Post("/domains"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := client.CreateDomain(context.Background(), DomainCreateRequest{
		Name:             "example.com",
		ReturnPathDomain: "pm-bounces.example.com",
	})

	if err != nil {
		t.Fatalf("CreateDomain: %s", err.Error())
	}

	if res.Name != "example.com" {
		t.Fatalf("CreateDomain: wrong name!: %s", res.Name)
	}

	if res.ReturnPathDomain != "pm-bounces.example.com" {
		t.Fatalf("CreateDomain: wrong ReturnPathDomain!: %s", res.ReturnPathDomain)
	}
}

func TestEditDomain(t *testing.T) {
	responseJSON := `{
  "Name": "example.com",
  "SPFVerified": false,
  "SPFHost": "example.com",
  "SPFTextValue": "v=spf1 a mx include:spf.mtasv.net ~all",
  "DKIMVerified": false,
  "WeakDKIM": false,
  "DKIMHost": "20160921046319pm._domainkey.example.com",
  "DKIMTextValue": "k=rsa; p=MIGfMA0GDRrFQJc5dZEBAQUAA4GNADCBiQKBgQCFn...",
  "DKIMPendingHost": "20131031155228pm._domainkey.example.com",
  "DKIMPendingTextValue": "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCFn...",
  "DKIMRevokedHost": "",
  "DKIMRevokedTextValue": "",
  "SafeToRemoveRevokedKeyFromDNS": false,
  "DKIMUpdateStatus": "Pending",
  "ReturnPathDomain": "pm-bounces.example.com",
  "ReturnPathDomainVerified": false,
  "ReturnPathDomainCNAMEValue": "pm.mtasv.net",
  "ID": 1234
}`

	tMux.HandleFunc(pat.Put("/domains/:domainID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := client.EditDomain(context.Background(), "1234", DomainEditRequest{
		ReturnPathDomain: "pm-bounces.example.com",
	})
	if err != nil {
		t.Fatalf("EditDomain: %s", err.Error())
	}

	if res.Name != "example.com" {
		t.Fatalf("EditDomain: wrong name!: %s", res.Name)
	}

	if res.ReturnPathDomain != "pm-bounces.example.com" {
		t.Fatalf("EditDomain: wrong ReturnPathDomain!: %s", res.ReturnPathDomain)
	}
}

func TestDeleteDomain(t *testing.T) {
	responseJSON := `{
	  "ErrorCode": 0,
	  "Message": "Domain example.com removed."
	}`

	tMux.HandleFunc(pat.Delete("/domains/:domainID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	// Success
	err := client.DeleteDomain(context.Background(), "1234")
	if err != nil {
		t.Fatalf("DeleteDomain: %s", err.Error())
	}

	// Failure
	responseJSON = `{
	  "ErrorCode": 402,
	  "Message": "Invalid JSON"
	}`

	err = client.DeleteDomain(context.Background(), "1234")
	if err == nil {
		t.Fatalf("DeleteDomain: should have failed")
	}
}
