package postmark

import (
	"context"
	"fmt"
	"net/http"

	"goji.io/pat"
)

func (s *PostmarkTestSuite) TestGetOutboundMessage() {
	responseJSON := `{
	  "TextBody": "Thank you for your order...",
	  "HtmlBody": "<p>Thank you for your order...</p>",
	  "Body": "SMTP dump data",
	  "Tag": "product-orders",
	  "MessageID": "07311c54-0687-4ab9-b034-b54b5bad88ba",
	  "To": [
		{
		  "Email": "john.doe@yahoo.com",
		  "Name": null
		}
	  ],
	  "Cc": [],
	  "Bcc": [],
	  "Recipients": [
		"john.doe@yahoo.com"
	  ],
	  "ReceivedAt": "2014-02-14T11:12:54.8054242-05:00",
	  "From": "\"Joe\" <joe@domain.com>",
	  "Subject": "Parts Order #5454",
	  "Attachments": ["test-file.txt"],
	  "Status": "Sent",
	  "MessageEvents": [
		{
		  "Recipient": "john.doe@yahoo.com",
		  "Type": "Delivered",
		  "ReceivedAt": "2014-02-14T11:13:10.8054242-05:00",
		  "Details": {
			"DeliveryMessage": "smtp;250 2.0.0 OK l10si21599969igu.63 - gsmtp",
			"DestinationServer": "yahoo-smtp-in.l.yahoo.com (433.899.888.26)",
			"DestinationIP": "173.194.74.256"
		  }
		},
		{
		  "Recipient": "john.doe@yahoo.com",
		  "Type": "Opened",
		  "ReceivedAt": "2014-02-14T11:20:10.8054242-05:00",
		  "Details": {
			"Summary": "Email opened with Mozilla/5.0 (Windows NT 5.1; rv:11.0) Gecko Firefox/11.0 (via ggpht.com GoogleImageProxy)"
		  }
		},
		{
		  "Recipient": "badrecipient@example.com",
		  "Type": "Bounced",
		  "ReceivedAt": "2014-02-14T11:20:15.8054242-05:00",
		  "Details": {
			"Summary": "smtp;550 5.1.1 The email account that you tried to reach does not exist. Please try double-checking the recipient's email address for typos or unnecessary spaces.",
			"BounceID": "374814878"
		  }
		}
	  ]
	}`

	s.mux.HandleFunc(pat.Get("/messages/outbound/07311c54-0687-4ab9-b034-b54b5bad88ba/details"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetOutboundMessage(context.Background(), "07311c54-0687-4ab9-b034-b54b5bad88ba")
	s.Require().NoError(err)

	s.Equal("07311c54-0687-4ab9-b034-b54b5bad88ba", res.MessageID, "GetOutboundMessage: wrong MessageID")
}

func (s *PostmarkTestSuite) TestGetOutboundMessageDump() {
	dump := `From: \"John Doe\" <john.doe@yahoo.com> \r\nTo: \"john.doe@yahoo.com\" <john.doe@yahoo.com>\r\nReply-To: joe@domain.com\r\nDate: Fri, 14 Feb 2014 11:12:56 -0500\r\nSubject: Parts Order #5454\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: quoted-printable\r\nX-Mailer: aspNetEmail ver 4.0.0.22\r\nX-Job: 44013_34141\r\nX-virtual-MTA: shared1\r\nX-Complaints-To: abuse@postmarkapp.com\r\nX-PM-RCPT: |bTB8NDQwMTN8MzQxNDF8anBAd2lsZGJpdC5jb20=|\r\nX-PM-Tag: product-orders\r\nX-PM-Message-Id: 07311c54-0687-4ab9-b034-b54b5bad88ba\r\nMessage-ID: <SC-ORD-MAIL4390fbe08b95f4257984dcaed896b4730@SC-ORD-MAIL4>\r\n\r\nThank you for your order=2E=2E=2E\r\n`

	responseJSON := fmt.Sprintf(`{"Body": "%s"}`, dump)

	s.mux.HandleFunc(pat.Get("/messages/outbound/07311c54-0687-4ab9-b034-b54b5bad88ba/dump"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	_, err := s.client.GetOutboundMessageDump(context.Background(), "07311c54-0687-4ab9-b034-b54b5bad88ba")
	s.Require().NoError(err)
}

func (s *PostmarkTestSuite) TestGetOutboundMessages() {
	responseJSON := `{
	  "TotalCount": 194,
		"Messages": [
		  {
			"Tag": "Invitation",
			"MessageID": "0ac29aee-e1cd-480d-b08d-4f48548ff48d",
			"To": [
			  {
				"Email": "john.doe@yahoo.com",
				"Name": null
			  }
			],
			"Cc": [],
			"Bcc": [],
			"Recipients": [
			  "john.doe@yahoo.com"
			],
			"ReceivedAt": "2014-02-20T07:25:02.8782715-05:00",
			"From": "\"Joe\" <joe@domain.com>",
			"Subject": "staging",
			"Attachments": [],
			"Status": "Sent"
		  }
		]
	}`

	s.mux.HandleFunc(pat.Get("/messages/outbound"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	_, total, err := s.client.GetOutboundMessages(context.Background(), 100, 0, map[string]interface{}{
		"recipient": "john.doe@yahoo.com",
		"tag":       "welcome",
		"status":    "",
		"todate":    "2015-01-12",
		"fromdate":  "2015-01-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(194), total, "GetOutboundMessages: wrong total")
}

func (s *PostmarkTestSuite) TestGetOutboundMessagesOpens() {
	responseJSON := `{
		"TotalCount": 1,
		"Opens": [
		  {
			"FirstOpen": true,
			"Client": {
			  "Name": "Chrome 34.0.1847.131",
			  "Company": "Google Inc.",
			  "Family": "Chrome"
			},
			"OS": {
			  "Name": "OS X 10.7 Lion",
			  "Company": "Apple Computer, Inc.",
			  "Family": "OS X"
			},
			"Platform": "WebMail",
			"UserAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.131 Safari/537.36",
			"ReadSeconds": 16,
			"Geo": {
			  "CountryISOCode": "RS",
			  "Country": "Serbia",
			  "RegionISOCode": "VO",
			  "Region": "Autonomna Pokrajina Vojvodina",
			  "City": "Novi Sad",
			  "Zip": "21000",
			  "Coords": "45.2517,19.8369",
			  "IP": "188.2.95.4"
			},
			"MessageID": "927e56d4-dc66-4070-bbf0-1db76c2ae14b",
			"ReceivedAt": "2014-04-30T05:04:23.8768746-04:00",
			"Tag": "welcome-user",
			"Recipient": "john.doe@yahoo.com"
		  }
		]

	}`
	s.mux.HandleFunc(pat.Get("/messages/outbound/opens"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	_, total, err := s.client.GetOutboundMessagesOpens(context.Background(), 100, 0, map[string]interface{}{
		"recipient": "john.doe@yahoo.com",
	})
	s.Require().NoError(err)

	s.Equal(int64(1), total, "GetOutboundMessagesOpens: wrong total")
}

func (s *PostmarkTestSuite) TestGetOutboundMessageOpens() {
	responseJSON := `{
		"TotalCount": 1,
	  "Opens": [
		{
		  "Client": {
			"Name": "Chrome 34.0.1847.131",
			"Company": "Google Inc.",
			"Family": "Chrome"
		  },
		  "OS": {
			"Name": "OS X 10.7 Lion",
			"Company": "Apple Computer, Inc.",
			"Family": "OS X"
		  },
		  "Platform": "WebMail",
		  "UserAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.131 Safari/537.36",
		  "ReadSeconds": 16,
		  "Geo": {
			"CountryISOCode": "RS",
			"Country": "Serbia",
			"RegionISOCode": "VO",
			"Region": "Autonomna Pokrajina Vojvodina",
			"City": "Novi Sad",
			"Zip": "21000",
			"Coords": "45.2517,19.8369",
			"IP": "188.2.95.4"
		  },
		  "MessageID": "927e56d4-dc66-4070-bbf0-1db76c2ae14b",
		  "ReceivedAt": "2014-04-30T05:04:23.8768746-04:00",
		  "Tag": "welcome-user",
		  "Recipient": "john.doe@yahoo.com"
		}
	  ]
	}`

	s.mux.HandleFunc(pat.Get("/messages/outbound/opens/927e56d4-dc66-4070-bbf0-1db76c2ae14b"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	_, total, err := s.client.GetOutboundMessageOpens(context.Background(), "927e56d4-dc66-4070-bbf0-1db76c2ae14b", 100, 0)
	s.Require().NoError(err)

	s.Equal(int64(1), total, "GetOutboundMessageOpens: wrong total")
}
