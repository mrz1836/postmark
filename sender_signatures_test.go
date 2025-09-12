package postmark

import (
	"context"
	"net/http"
)

func (s *PostmarkTestSuite) TestGetSenderSignatures() {
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

	s.mux.Get("/senders", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetSenderSignatures(context.Background(), 50, 0)
	s.Require().NoError(err)

	s.Equal(int(2), res.TotalCount, "GetSenderSignatures: wrong TotalCount")
}

func (s *PostmarkTestSuite) TestGetSenderSignature() {
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

	s.mux.Get("/senders/:signatureID", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetSenderSignature(context.Background(), 1234)
	s.Require().NoError(err)

	s.Equal("JP Toto", res.Name, "SenderSignature: wrong name")
}

func (s *PostmarkTestSuite) TestCreateSenderSignature() {
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

	s.mux.Post("/senders", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.CreateSenderSignature(context.Background(), SenderSignatureCreateRequest{
		FromEmail:                "john.doe@example.com",
		Name:                     "John Doe",
		ReplyToEmail:             "reply@example.com",
		ReturnPathDomain:         "pm-bounces.example.com",
		ConfirmationPersonalNote: "This is a note visible to the recipient to provide context of what Postmark is.",
	})
	s.Require().NoError(err)

	s.Equal("John Doe", res.Name, "CreateSenderSignature: wrong name")
}

func (s *PostmarkTestSuite) TestEditSenderSignature() {
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

	s.mux.Put("/senders/:signatureID", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.EditSenderSignature(context.Background(), 1, SenderSignatureEditRequest{
		Name:                     "Jane Doe",
		ReplyToEmail:             "jane.doe@example.com",
		ReturnPathDomain:         "pm-bounces.example.com",
		ConfirmationPersonalNote: "This is a note visible to the recipient to provide context of what Postmark is.",
	})
	s.Require().NoError(err)

	s.Equal("Jane Doe", res.Name, "EditSenderSignature: wrong name")
}

func (s *PostmarkTestSuite) TestDeleteSenderSignature() {
	responseJSON := `{
	  "ErrorCode": 0,
	  "Message": "SenderSignature 1234 removed."
	}`

	s.mux.Delete("/senders/:signatureID", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	// Success
	err := s.client.DeleteSenderSignature(context.Background(), 1234)
	s.Require().NoError(err)

	// Failure
	responseJSON = `{
	  "ErrorCode": 402,
	  "Message": "Invalid JSON"
	}`

	err = s.client.DeleteSenderSignature(context.Background(), 1234)
	s.Require().Error(err, "DeleteSenderSignature should have failed")
}

func (s *PostmarkTestSuite) TestResendSenderSignatureConfirmation() {
	tests := []struct {
		name         string
		responseJSON string
		wantErr      bool
		errContains  string
	}{
		{
			name: "successful resend confirmation",
			responseJSON: `{
				"ErrorCode": 0,
				"Message": "Confirmation resent to 'test@example.com'"
			}`,
			wantErr: false,
		},
		{
			name: "resend confirmation failure",
			responseJSON: `{
				"ErrorCode": 406,
				"Message": "You already have a confirmed signature with this email address."
			}`,
			wantErr:     true,
			errContains: "confirmed signature",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mux.Post("/senders/:signatureID/resend", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(tt.responseJSON))
			})

			err := s.client.ResendSenderSignatureConfirmation(context.Background(), 1234)

			if tt.wantErr {
				s.Require().Error(err, "ResendSenderSignatureConfirmation should fail")
				if tt.errContains != "" {
					s.Contains(err.Error(), tt.errContains, "Error should contain expected message")
				}
			} else {
				s.Require().NoError(err, "ResendSenderSignatureConfirmation should not fail")
			}
		})
	}
}
