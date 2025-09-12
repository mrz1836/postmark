package postmark

import (
	"context"
	"net/http"

	"goji.io/pat"
)

func (s *PostmarkTestSuite) TestGetDomain() {
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

	s.mux.HandleFunc(pat.Get("/domains/:domainID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetDomain(context.Background(), 1234)
	s.Require().NoError(err, "GetDomain should not fail")
	s.Equal("postmarkapp.com", res.Name, "GetDomain should return correct domain name")
}

func (s *PostmarkTestSuite) TestCreateDomain() {
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

	s.mux.HandleFunc(pat.Post("/domains"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.CreateDomain(context.Background(), DomainCreateRequest{
		Name:             "example.com",
		ReturnPathDomain: "pm-bounces.example.com",
	})
	s.Require().NoError(err, "CreateDomain should not fail")
	s.Equal("example.com", res.Name, "CreateDomain should return correct domain name")
	s.Equal("pm-bounces.example.com", res.ReturnPathDomain, "CreateDomain should return correct return path domain")
}

func (s *PostmarkTestSuite) TestEditDomain() {
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

	s.mux.HandleFunc(pat.Put("/domains/:domainID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.EditDomain(context.Background(), 1234, DomainEditRequest{
		ReturnPathDomain: "pm-bounces.example.com",
	})
	s.Require().NoError(err, "EditDomain should not fail")
	s.Equal("example.com", res.Name, "EditDomain should return correct domain name")
	s.Equal("pm-bounces.example.com", res.ReturnPathDomain, "EditDomain should return correct return path domain")
}

func (s *PostmarkTestSuite) TestDeleteDomain() {
	tests := []struct {
		name         string
		responseJSON string
		wantErr      bool
	}{
		{
			name: "successful domain deletion",
			responseJSON: `{
				"ErrorCode": 0,
				"Message": "Domain example.com removed."
			}`,
			wantErr: false,
		},
		{
			name: "domain deletion failure",
			responseJSON: `{
				"ErrorCode": 402,
				"Message": "Invalid JSON"
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mux.HandleFunc(pat.Delete("/domains/:domainID"), func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(tt.responseJSON))
			})

			err := s.client.DeleteDomain(context.Background(), 1234)

			if tt.wantErr {
				s.Require().Error(err, "DeleteDomain should fail")
			} else {
				s.Require().NoError(err, "DeleteDomain should not fail")
			}
		})
	}
}
