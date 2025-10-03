package postmark

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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

	s.mux.Get("/messages/outbound/07311c54-0687-4ab9-b034-b54b5bad88ba/details", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetOutboundMessage(context.Background(), "07311c54-0687-4ab9-b034-b54b5bad88ba")
	s.Require().NoError(err)

	s.Equal("07311c54-0687-4ab9-b034-b54b5bad88ba", res.MessageID, "GetOutboundMessage: wrong MessageID")
}

func (s *PostmarkTestSuite) TestGetOutboundMessageDump() {
	dump := `From: \"John Doe\" <john.doe@yahoo.com> \r\nTo: \"john.doe@yahoo.com\" <john.doe@yahoo.com>\r\nReply-To: joe@domain.com\r\nDate: Fri, 14 Feb 2014 11:12:56 -0500\r\nSubject: Parts Order #5454\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: quoted-printable\r\nX-Mailer: aspNetEmail ver 4.0.0.22\r\nX-Job: 44013_34141\r\nX-virtual-MTA: shared1\r\nX-Complaints-To: abuse@postmarkapp.com\r\nX-PM-RCPT: |bTB8NDQwMTN8MzQxNDF8anBAd2lsZGJpdC5jb20=|\r\nX-PM-Tag: product-orders\r\nX-PM-Message-Id: 07311c54-0687-4ab9-b034-b54b5bad88ba\r\nMessage-ID: <SC-ORD-MAIL4390fbe08b95f4257984dcaed896b4730@SC-ORD-MAIL4>\r\n\r\nThank you for your order=2E=2E=2E\r\n`

	responseJSON := fmt.Sprintf(`{"Body": "%s"}`, dump)

	s.mux.Get("/messages/outbound/07311c54-0687-4ab9-b034-b54b5bad88ba/dump", func(w http.ResponseWriter, _ *http.Request) {
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

	s.mux.Get("/messages/outbound", func(w http.ResponseWriter, _ *http.Request) {
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
	s.mux.Get("/messages/outbound/opens", func(w http.ResponseWriter, _ *http.Request) {
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

	s.mux.Get("/messages/outbound/opens/927e56d4-dc66-4070-bbf0-1db76c2ae14b", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	_, total, err := s.client.GetOutboundMessageOpens(context.Background(), "927e56d4-dc66-4070-bbf0-1db76c2ae14b", 100, 0)
	s.Require().NoError(err)

	s.Equal(int64(1), total, "GetOutboundMessageOpens: wrong total")
}

func (s *PostmarkTestSuite) TestGetOutboundMessagesClicks() {
	responseJSON := `{
	  "TotalCount": 1,
	  "Clicks": [
		{
		  "RecordType": "Click",
		  "ClickLocation": "HTML",
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
		  "Platform": "Desktop",
		  "UserAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.131 Safari/537.36",
		  "OriginalLink": "http://example.com/click-me",
		  "Geo": {
			"CountryISOCode": "RS",
			"Country": "Serbia",
			"RegionISOCode": "VO",
			"Region": "Vojvodina",
			"City": "Novi Sad",
			"Zip": "21000",
			"Coords": "45.2517,19.8369",
			"IP": "188.2.95.4"
		  },
		  "MessageID": "927e56d4-dc66-4c01-a0be-645b4b6f5fd7",
		  "MessageStream": "outbound",
		  "ReceivedAt": "2014-02-14T11:13:10.8054242-05:00",
		  "Tag": "Invitation",
		  "Recipient": "john.doe@yahoo.com"
		}
	  ]
	}`

	s.mux.Get("/messages/outbound/clicks", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	clicks, count, err := s.client.GetOutboundMessagesClicks(context.Background(), 100, 0, map[string]interface{}{
		"tag":       "Invitation",
		"recipient": "john.doe@yahoo.com",
	})
	s.Require().NoError(err)

	s.EqualValues(1, count)
	s.Len(clicks, 1)

	click := clicks[0]
	s.Equal("927e56d4-dc66-4c01-a0be-645b4b6f5fd7", click.MessageID, "GetOutboundMessagesClicks: wrong MessageID")
	s.Equal("HTML", click.ClickLocation, "GetOutboundMessagesClicks: wrong ClickLocation")
	s.Equal("Chrome 34.0.1847.131", click.Client["Name"], "GetOutboundMessagesClicks: wrong Client Name")
	s.Equal("Google Inc.", click.Client["Company"], "GetOutboundMessagesClicks: wrong Client Company")
	s.Equal("Chrome", click.Client["Family"], "GetOutboundMessagesClicks: wrong Client Family")
	s.Equal("OS X 10.7 Lion", click.OS["Name"], "GetOutboundMessagesClicks: wrong OS Name")
	s.Equal("http://example.com/click-me", click.OriginalLink, "GetOutboundMessagesClicks: wrong OriginalLink")
	s.Equal("john.doe@yahoo.com", click.Recipient, "GetOutboundMessagesClicks: wrong Recipient")
	s.Equal("Invitation", click.Tag, "GetOutboundMessagesClicks: wrong Tag")
}

func (s *PostmarkTestSuite) TestGetOutboundMessageClicks() {
	responseJSON := `{
	  "TotalCount": 2,
	  "Clicks": [
		{
		  "RecordType": "Click",
		  "ClickLocation": "HTML",
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
		  "Platform": "Desktop",
		  "UserAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.131 Safari/537.36",
		  "OriginalLink": "http://example.com/click-me",
		  "Geo": {
			"CountryISOCode": "RS",
			"Country": "Serbia",
			"RegionISOCode": "VO",
			"Region": "Vojvodina",
			"City": "Novi Sad",
			"Zip": "21000",
			"Coords": "45.2517,19.8369",
			"IP": "188.2.95.4"
		  },
		  "MessageID": "927e56d4-dc66-4c01-a0be-645b4b6f5fd7",
		  "MessageStream": "outbound",
		  "ReceivedAt": "2014-02-14T11:13:10.8054242-05:00",
		  "Tag": "Invitation",
		  "Recipient": "john.doe@yahoo.com"
		},
		{
		  "RecordType": "Click",
		  "ClickLocation": "Text",
		  "Client": {
			"Name": "Safari 7.0.3",
			"Company": "Apple Computer, Inc.",
			"Family": "Safari"
		  },
		  "OS": {
			"Name": "OS X 10.9 Mavericks",
			"Company": "Apple Computer, Inc.",
			"Family": "OS X"
		  },
		  "Platform": "Desktop",
		  "UserAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.75.14 (KHTML, like Gecko) Version/7.0.3 Safari/537.75.14",
		  "OriginalLink": "http://example.com/another-link",
		  "Geo": {
			"CountryISOCode": "US",
			"Country": "United States",
			"RegionISOCode": "CA",
			"Region": "California",
			"City": "San Francisco",
			"Zip": "94102",
			"Coords": "37.7749,-122.4194",
			"IP": "192.168.1.1"
		  },
		  "MessageID": "927e56d4-dc66-4c01-a0be-645b4b6f5fd7",
		  "MessageStream": "outbound",
		  "ReceivedAt": "2014-02-14T11:15:10.8054242-05:00",
		  "Tag": "Invitation",
		  "Recipient": "jane.doe@gmail.com"
		}
	  ]
	}`

	s.mux.Get("/messages/outbound/clicks/927e56d4-dc66-4c01-a0be-645b4b6f5fd7", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	clicks, count, err := s.client.GetOutboundMessageClicks(context.Background(), "927e56d4-dc66-4c01-a0be-645b4b6f5fd7", 10, 0)
	s.Require().NoError(err)

	s.EqualValues(2, count)
	s.Len(clicks, 2)

	htmlClick := clicks[0]
	s.Equal("927e56d4-dc66-4c01-a0be-645b4b6f5fd7", htmlClick.MessageID, "GetOutboundMessageClicks: wrong MessageID for HTML click")
	s.Equal("HTML", htmlClick.ClickLocation, "GetOutboundMessageClicks: wrong ClickLocation for HTML click")
	s.Equal("http://example.com/click-me", htmlClick.OriginalLink, "GetOutboundMessageClicks: wrong OriginalLink for HTML click")

	textClick := clicks[1]
	s.Equal("Text", textClick.ClickLocation, "GetOutboundMessageClicks: wrong ClickLocation for Text click")
	s.Equal("http://example.com/another-link", textClick.OriginalLink, "GetOutboundMessageClicks: wrong OriginalLink for Text click")
	s.Equal("jane.doe@gmail.com", textClick.Recipient, "GetOutboundMessageClicks: wrong Recipient for Text click")
}

// Benchmark for GetOutboundMessagesClicks
func BenchmarkGetOutboundMessagesClicks(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
	  "TotalCount": 1,
	  "Clicks": [
		{
		  "RecordType": "Click",
		  "ClickLocation": "HTML",
		  "OriginalLink": "http://example.com/click-me",
		  "MessageID": "927e56d4-dc66-4c01-a0be-645b4b6f5fd7",
		  "MessageStream": "outbound",
		  "ReceivedAt": "2014-02-14T11:13:10.8054242-05:00",
		  "Tag": "Invitation",
		  "Recipient": "john.doe@example.com"
		}
	  ]
	}`

	mux.Get("/messages/outbound/clicks", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	options := map[string]interface{}{
		"tag":       "Invitation",
		"recipient": "john.doe@example.com",
		"platform":  "Desktop",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = client.GetOutboundMessagesClicks(context.Background(), 100, 0, options)
	}
}

// Benchmark for GetOutboundMessageClicks
func BenchmarkGetOutboundMessageClicks(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
	  "TotalCount": 1,
	  "Clicks": [
		{
		  "RecordType": "Click",
		  "ClickLocation": "HTML",
		  "OriginalLink": "http://example.com/click-me",
		  "MessageID": "927e56d4-dc66-4c01-a0be-645b4b6f5fd7",
		  "MessageStream": "outbound",
		  "ReceivedAt": "2014-02-14T11:13:10.8054242-05:00",
		  "Tag": "Invitation",
		  "Recipient": "john.doe@example.com"
		}
	  ]
	}`

	mux.Get("/messages/outbound/clicks/927e56d4-dc66-4c01-a0be-645b4b6f5fd7", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = client.GetOutboundMessageClicks(context.Background(), "927e56d4-dc66-4c01-a0be-645b4b6f5fd7", 50, 0)
	}
}
