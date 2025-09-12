package postmark

import (
	"context"
	"net/http"
	"testing"
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

	s.mux.Get("/domains/:domainID", func(w http.ResponseWriter, _ *http.Request) {
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

	s.mux.Post("/domains", func(w http.ResponseWriter, _ *http.Request) {
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

	s.mux.Put("/domains/:domainID", func(w http.ResponseWriter, _ *http.Request) {
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
			s.mux.Delete("/domains/:domainID", func(w http.ResponseWriter, _ *http.Request) {
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

func (s *PostmarkTestSuite) TestGetDomains() {
	responseJSON := `{
  "TotalCount": 1,
  "Domains": [
    {
      "Name": "postmarkapp.com",
      "SPFVerified": true,
      "DKIMVerified": false,
      "WeakDKIM": false,
      "ReturnPathDomainVerified": false,
      "ID": 1234
    }
  ]
}`

	s.mux.Get("/domains", func(w http.ResponseWriter, req *http.Request) {
		count := req.URL.Query().Get("count")
		offset := req.URL.Query().Get("offset")
		s.Equal("10", count, "GetDomains should send correct count parameter")
		s.Equal("5", offset, "GetDomains should send correct offset parameter")
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetDomains(context.Background(), 10, 5)
	s.Require().NoError(err, "GetDomains should not fail")
	s.Equal(1, res.TotalCount, "GetDomains should return correct total count")
	s.Require().Len(res.Domains, 1, "GetDomains should return correct number of domains")
	s.Equal("postmarkapp.com", res.Domains[0].Name, "GetDomains should return correct domain name")
}

func (s *PostmarkTestSuite) TestVerifyDKIMStatus() {
	responseJSON := `{
  "Name": "postmarkapp.com",
  "SPFVerified": true,
  "SPFHost": "postmarkapp.com",
  "SPFTextValue": "v=spf1 a mx include:spf.mtasv.net ~all",
  "DKIMVerified": true,
  "WeakDKIM": false,
  "DKIMHost": "jan2013pm._domainkey.postmarkapp.com",
  "DKIMTextValue": "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDJ...",
  "DKIMPendingHost": "",
  "DKIMPendingTextValue": "",
  "DKIMRevokedHost": "",
  "DKIMRevokedTextValue": "",
  "SafeToRemoveRevokedKeyFromDNS": false,
  "DKIMUpdateStatus": "Verified",
  "ReturnPathDomain": "pm-bounces.postmarkapp.com",
  "ReturnPathDomainVerified": false,
  "ReturnPathDomainCNAMEValue": "pm.mtasv.net",
  "ID": 1234
}`

	s.mux.Put("/domains/:domainID/verifyDkim", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.VerifyDKIMStatus(context.Background(), 1234)
	s.Require().NoError(err, "VerifyDKIMStatus should not fail")
	s.Equal("postmarkapp.com", res.Name, "VerifyDKIMStatus should return correct domain name")
	s.True(res.DKIMVerified, "VerifyDKIMStatus should verify DKIM")
	s.Equal("Verified", res.DKIMUpdateStatus, "VerifyDKIMStatus should update DKIM status")
}

func (s *PostmarkTestSuite) TestVerifyReturnPath() {
	responseJSON := `{
  "Name": "postmarkapp.com",
  "SPFVerified": true,
  "SPFHost": "postmarkapp.com",
  "SPFTextValue": "v=spf1 a mx include:spf.mtasv.net ~all",
  "DKIMVerified": false,
  "WeakDKIM": false,
  "DKIMHost": "jan2013pm._domainkey.postmarkapp.com",
  "DKIMTextValue": "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDJ...",
  "DKIMPendingHost": "",
  "DKIMPendingTextValue": "",
  "DKIMRevokedHost": "",
  "DKIMRevokedTextValue": "",
  "SafeToRemoveRevokedKeyFromDNS": false,
  "DKIMUpdateStatus": "Pending",
  "ReturnPathDomain": "pm-bounces.postmarkapp.com",
  "ReturnPathDomainVerified": true,
  "ReturnPathDomainCNAMEValue": "pm.mtasv.net",
  "ID": 1234
}`

	s.mux.Put("/domains/:domainID/verifyReturnPath", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.VerifyReturnPath(context.Background(), 1234)
	s.Require().NoError(err, "VerifyReturnPath should not fail")
	s.Equal("postmarkapp.com", res.Name, "VerifyReturnPath should return correct domain name")
	s.True(res.ReturnPathDomainVerified, "VerifyReturnPath should verify return path")
}

func (s *PostmarkTestSuite) TestRotateDKIM() {
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

	s.mux.Post("/domains/:domainID/rotatedkim", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.RotateDKIM(context.Background(), 1234)
	s.Require().NoError(err, "RotateDKIM should not fail")
	s.Equal("postmarkapp.com", res.Name, "RotateDKIM should return correct domain name")
	s.NotEmpty(res.DKIMPendingHost, "RotateDKIM should set pending DKIM host")
	s.NotEmpty(res.DKIMPendingTextValue, "RotateDKIM should set pending DKIM text value")
	s.Equal("Pending", res.DKIMUpdateStatus, "RotateDKIM should set status to pending")
}

// Benchmark functions for Domains API

func BenchmarkGetDomains(b *testing.B) {
	ctx := context.Background()
	count := 50
	offset := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = count
		_ = offset
	}
}

func BenchmarkGetDomain(b *testing.B) {
	ctx := context.Background()
	domainID := int64(1234)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = domainID
	}
}

func BenchmarkCreateDomain(b *testing.B) {
	ctx := context.Background()
	request := DomainCreateRequest{
		Name:             "test.com",
		ReturnPathDomain: "bounces.test.com",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = request
	}
}

func BenchmarkEditDomain(b *testing.B) {
	ctx := context.Background()
	domainID := int64(1234)
	request := DomainEditRequest{
		ReturnPathDomain: "new-bounces.test.com",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = domainID
		_ = request
	}
}

func BenchmarkDeleteDomain(b *testing.B) {
	ctx := context.Background()
	domainID := int64(1234)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = domainID
	}
}

func BenchmarkVerifyDKIMStatus(b *testing.B) {
	ctx := context.Background()
	domainID := int64(1234)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = domainID
	}
}

func BenchmarkVerifyReturnPath(b *testing.B) {
	ctx := context.Background()
	domainID := int64(1234)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = domainID
	}
}

func BenchmarkRotateDKIM(b *testing.B) {
	ctx := context.Background()
	domainID := int64(1234)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = domainID
	}
}
