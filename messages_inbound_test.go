package postmark

import (
	"context"
	"net/http"
	"net/http/httptest"
)

func (s *PostmarkTestSuite) TestGetInboundMessage() {
	responseJSON := `{
		"From": "dart-zzzzz@yandex.ru",
		  "FromName": "Dart Zzzzz",
		  "FromFull": {
			"Email": "dart-zzzzz@yandex.ru",
			"Name": "Dart Zzzzz"
		  },
		  "To": "ad8a4d0842c486355a33a7f019caab51@inbound.postmarkapp.com",
		  "ToFull": [
			{
			  "Email": "ad8a4d0842c486355a33a7f019caab51@inbound.postmarkapp.com",
			  "Name": ""
			}
		  ],
		  "CcFull": [],
		  "Cc": "",
		  "ReplyTo": "",
		  "OriginalRecipient": "ad8a4d0842c486355a33a7f019caab51@inbound.postmarkapp.com",
		  "Subject": "Тест.",
		  "Date": "Thu, 13 Feb 2014 17:48:22 +0300",
		  "MailboxHash": "",
		  "TextBody": "stuff stuff.",
		  "HtmlBody": "",
		  "Tag": "",
		  "Headers": [
			{
			  "Name": "X-Spam-Checker-Version",
			  "Value": "SpamAssassin 3.3.1 (2010-03-16) on sc-ord-inbound1"
			},
			{
			  "Name": "X-Spam-Status",
			  "Value": "No"
			},
			{
			  "Name": "X-Spam-Score",
			  "Value": "0.7"
			},
			{
			  "Name": "X-Spam-Tests",
			  "Value": "DKIM_SIGNED,DKIM_VALID,DKIM_VALID_AU,FREEMAIL_FROM,FSL_HELO_BARE_IP_2,RCVD_IN_DNSWL_LOW,SPF_PASS"
			},
			{
			  "Name": "Received-SPF",
			  "Value": "Pass (sender SPF authorized) identity=mailfrom; client-ip=95.108.130.92; helo=forward14.mail.yandex.net; envelope-from=dart-zzzzz@yandex.ru; receiver=ad8a4d0842c486355a33a7f019caab51@inbound.postmarkapp.com"
			},
			{
			  "Name": "DKIM-Signature",
			  "Value": "v=1; a=rsa-sha256; c=relaxed/relaxed; d=yandex.ru; s=mail;t=1392302902; bh=4mN45y6KsGBYQjvZYsA49+gc9iuptslitnW5OR+Gg0M=;h=From:To:Subject:Date;b=StRtIzi3pvGDORwJkDc49RGqcgvlFvUEqAXi8RoHGu3LvHQmZs0F2pRdqc5UYt1gO OvLSKhlDslDkACdSJQAkj6EF99gXgiLItWo7hNfbv03qDlIq27f8vCZN5Uw0DY5shQ mVatnZbP/L01YP1pTXQONaalDFJ4ByRjjrWDrFVI="
			},
			{
			  "Name": "Envelope-From",
			  "Value": "Dart-zzzzz@yandex.ua"
			},
			{
			  "Name": "MIME-Version",
			  "Value": "1.0"
			},
			{
			  "Name": "Message-Id",
			  "Value": "<51351392302902@web19j.yandex.ru>"
			},
			{
			  "Name": "X-Mailer",
			  "Value": "Yamail [ http://yandex.ru ] 5.0"
			},
			{
			  "Name": "Content-Transfer-Encoding",
			  "Value": "8bit"
			}
		  ],
		  "Attachments": [],
		  "MessageID": "cc5727a0-ea30-4e79-baea-aa43c9628ac4",
		  "BlockedReason": "Inbound request blocked by domain rule: badsender@example.com",
		  "Status": "Blocked"
	}`

	s.mux.Get("/messages/inbound/cc5727a0-ea30-4e79-baea-aa43c9628ac4/details", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetInboundMessage(context.Background(), "cc5727a0-ea30-4e79-baea-aa43c9628ac4")
	s.Require().NoError(err)

	s.Equal("cc5727a0-ea30-4e79-baea-aa43c9628ac4", res.MessageID, "GetInboundMessage: wrong MessageID")

	_, err = res.Time()
	s.Require().NoError(err, "GetInboundMessage: date couldn't be parsed: %s", res.Date)
}

func (s *PostmarkTestSuite) TestGetInboundMessages() {
	responseJSON := `{
		"TotalCount": 7,
	   	"InboundMessages": [
		 {
		   "From": "dart-zzzzz@yandex.ru",
		   "FromName": "Dart Zzzzz",
		   "FromFull": {
			 "Email": "dart-zzzzz@yandex.ru",
			 "Name": "Dart Zzzzz"
		   },
		   "To": "ad8a4d0842c486355a33a7f019caab51@inbound.postmarkapp.com",
		   "ToFull": [
			 {
			   "Email": "ad8a4d0842c486355a33a7f019caab51@inbound.postmarkapp.com",
			   "Name": ""
			 }
		   ],
		   "CcFull": [],
		   "Cc": "",
		   "ReplyTo": "",
		   "OriginalRecipient": "ad8a4d0842c486355a33a7f019caab51@inbound.postmarkapp.com",
		   "Subject": "Тест.",
		   "Date": "Thu, 13 Feb 2014 17:48:22 +0300",
		   "MailboxHash": "",
		   "Tag": "",
		   "Attachments": [],
		   "MessageID": "cc5727a0-ea30-4e79-baea-aa43c9628ac4",
		   "Status": "Blocked"
		 }
	   ]
	}`

	s.mux.Get("/messages/inbound", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	_, total, err := s.client.GetInboundMessages(context.Background(), 100, 0, map[string]interface{}{
		"recipient": "john.doe@yahoo.com",
		"fromdate":  "2015-02-01",
		"todate":    "2015-03-01",
		"status":    "blocked",
	})
	s.Require().NoError(err)

	s.Equal(int64(7), total, "GetInboundMessages: wrong total")
}

func (s *PostmarkTestSuite) TestGetInboundMessagesNilOptions() {
	// Create separate mux/server to avoid conflicts with TestGetInboundMessages
	errorMux := NewTestRouter()
	errorServer := httptest.NewServer(errorMux)
	defer errorServer.Close()

	errorClient := NewClient("server-token", "account-token")
	errorClient.BaseURL = errorServer.URL

	responseJSON := `{
		"TotalCount": 1,
		"InboundMessages": []
	}`

	errorMux.Get("/messages/inbound", func(w http.ResponseWriter, req *http.Request) {
		// Verify count and offset are still in query params even with nil options
		s.Equal("50", req.URL.Query().Get("count"))
		s.Equal("0", req.URL.Query().Get("offset"))
		_, _ = w.Write([]byte(responseJSON))
	})

	_, total, err := errorClient.GetInboundMessages(context.Background(), 50, 0, nil)
	s.Require().NoError(err)
	s.Equal(int64(1), total)
}

func (s *PostmarkTestSuite) TestBypassInboundMessage() {
	tests := []struct {
		name         string
		responseJSON string
		statusCode   int
		wantErr      bool
	}{
		{
			name: "success",
			responseJSON: `{
				"ErrorCode": 0,
				"Message": "Successfully bypassed message: 792a3e9d-0078-40df-a6b0-fc78f87bf277."
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "api error code",
			responseJSON: `{
				"ErrorCode": 701,
				"Message": "This message was not found or cannot be bypassed."
			}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "http error",
			responseJSON: `{
				"ErrorCode": 500,
				"Message": "Internal Server Error"
			}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mux.Put("/messages/inbound/792a3e9d-0078-40df-a6b0-fc78f87bf277/bypass", func(w http.ResponseWriter, _ *http.Request) {
				if tt.statusCode != http.StatusOK {
					w.WriteHeader(tt.statusCode)
				}
				_, _ = w.Write([]byte(tt.responseJSON))
			})

			err := s.client.BypassInboundMessage(context.Background(), "792a3e9d-0078-40df-a6b0-fc78f87bf277")

			if tt.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *PostmarkTestSuite) TestRetryInboundMessage() {
	tests := []struct {
		name         string
		responseJSON string
		statusCode   int
		wantErr      bool
	}{
		{
			name: "success",
			responseJSON: `{
				"ErrorCode": 0,
				"Message": "Successfully rescheduled failed message: 041e3d29-737d-491e-9a13-a94d3rjkjka13."
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "api error code",
			responseJSON: `{
				"ErrorCode": 701,
				"Message": "This message was not found or cannot be retried."
			}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "http error",
			responseJSON: `{
				"ErrorCode": 500,
				"Message": "Internal Server Error"
			}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mux.Put("/messages/inbound/041e3d29-737d-491e-9a13-a94d3rjkjka13/retry", func(w http.ResponseWriter, _ *http.Request) {
				if tt.statusCode != http.StatusOK {
					w.WriteHeader(tt.statusCode)
				}
				_, _ = w.Write([]byte(tt.responseJSON))
			})

			err := s.client.RetryInboundMessage(context.Background(), "041e3d29-737d-491e-9a13-a94d3rjkjka13")

			if tt.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
