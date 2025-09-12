package postmark

import (
	"context"
	"net/http"
)

func getTestEmail() Email {
	return Email{
		From:     "sender@example.com",
		To:       "receiver@example.com",
		Cc:       "copied@example.com",
		Bcc:      "blank-copied@example.com",
		Subject:  "Test",
		Tag:      "Invitation",
		HTMLBody: "<b>Hello</b>",
		TextBody: "Hello",
		ReplyTo:  "reply@example.com",
		Headers: []Header{
			{
				Name:  "CUSTOM-HEADER",
				Value: "value",
			},
		},
		TrackOpens: true,
		InlineCSS:  true,
		Attachments: []Attachment{
			{
				Name:        "readme.txt",
				Content:     "dGVzdCBjb250ZW50",
				ContentType: "text/plain",
			},
			{
				Name:        "report.pdf",
				Content:     "dGVzdCBjb250ZW50",
				ContentType: "application/octet-stream",
			},
		},
	}
}

func (s *PostmarkTestSuite) TestSendEmail() {
	tests := []struct {
		name         string
		responseJSON string
		wantErr      bool
		expectedID   string
	}{
		{
			name: "successful email send",
			responseJSON: `{
				"To": "receiver@example.com",
				"SubmittedAt": "2014-02-17T07:25:01.4178645-05:00",
				"MessageID": "0a129aee-e1cd-480d-b08d-4f48548ff48d",
				"ErrorCode": 0,
				"Message": "OK"
			}`,
			wantErr:    false,
			expectedID: "0a129aee-e1cd-480d-b08d-4f48548ff48d",
		},
		{
			name: "email send failure with error code",
			responseJSON: `{
				"To": "receiver@example.com",
				"SubmittedAt": "2014-02-17T07:25:01.4178645-05:00",
				"MessageID": "0a129aee-e1cd-480d-b08d-4f48548ff48d",
				"ErrorCode": 401,
				"Message": "Sender signature not confirmed"
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mux.Post("/email", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(tt.responseJSON))
			})

			res, err := s.client.SendEmail(context.Background(), getTestEmail())

			if tt.wantErr {
				s.Require().Error(err, "SendEmail should have failed")
			} else {
				s.Require().NoError(err, "SendEmail should not have failed")
				s.Equal(tt.expectedID, res.MessageID, "SendEmail returned wrong message ID")
			}
		})
	}
}

func (s *PostmarkTestSuite) TestSendEmailBatch() {
	responseJSON := `[
	  {
		"ErrorCode": 0,
		"Message": "OK",
		"MessageID": "b7bc2f4a-e38e-4336-af7d-e6c392c2f817",
		"SubmittedAt": "2010-11-26T12:01:05.1794748-05:00",
		"To": "receiver1@example.com"
	  },
	  {
		"ErrorCode": 0,
		"Message": "OK",
		"MessageID": "e2ecbbfc-fe12-463d-b933-9fe22915106d",
		"SubmittedAt": "2010-11-26T12:01:05.1794748-05:00",
		"To": "receiver2@example.com"
	  }
	]`

	s.mux.Post("/email/batch", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	testEmail := getTestEmail()
	res, err := s.client.SendEmailBatch(context.Background(), []Email{testEmail, testEmail})
	s.Require().NoError(err, "SendEmailBatch should not have failed")
	s.Len(res, 2, "SendEmailBatch should return 2 results")
}
