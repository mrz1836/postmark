package postmark

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (s *PostmarkTestSuite) TestGetDeliveryStats() {
	responseJSON := `{
	  "InactiveMails": 192,
	  "Bounces": [
		{
		  "Name": "All",
		  "Count": 253
		},
		{
		  "Type": "HardBounce",
		  "Name": "Hard bounce",
		  "Count": 195
		},
		{
		  "Type": "Transient",
		  "Name": "Message delayed",
		  "Count": 10
		},
		{
		  "Type": "AutoResponder",
		  "Name": "Auto responder",
		  "Count": 14
		},
		{
		  "Type": "SpamNotification",
		  "Name": "Spam notification",
		  "Count": 3
		},
		{
		  "Type": "SoftBounce",
		  "Name": "Soft bounce",
		  "Count": 30
		},
		{
		  "Type": "SpamComplaint",
		  "Name": "Spam complaint",
		  "Count": 1
		}
	]}`

	s.mux.Get("/deliverystats", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetDeliveryStats(context.Background())
	s.Require().NoError(err, "GetDeliveryStats should not fail")
	s.Equal(int64(192), res.InactiveMails, "GetDeliveryStats should return correct inactive mail count")
}

func (s *PostmarkTestSuite) TestGetBounces() {
	responseJSON := `{
	  "TotalCount": 253,
	  "Bounces": [
		{
		  "RecordType": "Bounce",
		  "ID": 692560173,
		  "Type": "HardBounce",
		  "TypeCode": 1,
		  "Name": "Hard bounce",
		  "Tag": "Invitation",
		  "MessageID": "2c1b63fe-43f2-4db5-91b0-8bdfa44a9316",
		  "MessageStream": "outbound",
		  "Description": "The server was unable to deliver your message (ex: unknown user, mailbox not found).",
		  "Details": "action: failed\r\n",
		  "Email": "anything@blackhole.postmarkap.com",
		  "BouncedAt": "2014-01-15T16:09:19.6421112-05:00",
		  "DumpAvailable": false,
		  "Inactive": false,
		  "CanActivate": true,
		  "Subject": "SC API5 Test",
		  "Content": "Return-Path: <>\r\nReceived: ..."
		},
		{
		  "RecordType": "Bounce",
		  "ID": 676862817,
		  "Type": "HardBounce",
		  "TypeCode": 1,
		  "Name": "Hard bounce",
		  "Tag": "Invitation",
		  "MessageID": "623b2e90-82d0-4050-ae9e-2c3a734ba091",
		  "MessageStream": "outbound",
		  "Description": "The server was unable to deliver your message (ex: unknown user, mailbox not found).",
		  "Details": "smtp;554 delivery error: dd This user doesn't have a yahoo.com account (vicelcown@yahoo.com) [0] - mta1543.mail.ne1.yahoo.com",
		  "Email": "vicelcown@yahoo.com",
		  "BouncedAt": "2013-10-18T09:49:59.8253577-04:00",
		  "DumpAvailable": false,
		  "Inactive": true,
		  "CanActivate": true,
		  "Subject": "Production API Test",
		  "Content": "Return-Path: <>\r\nReceived: ..."
		}
		  ]
	}`

	s.mux.Get("/bounces", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	bounces, total, err := s.client.GetBounces(context.Background(), 100, 0, map[string]interface{}{
		"tag": "Invitation",
	})
	s.Require().NoError(err, "GetBounces should not fail")
	s.Equal(int64(253), total, "GetBounces should return correct total count")
	s.Require().Len(bounces, 2, "GetBounces should return 2 bounces")
	s.Equal("Bounce", bounces[0].RecordType, "GetBounces should return correct record type")
	s.Equal("outbound", bounces[0].MessageStream, "GetBounces should return correct message stream")
	s.Equal("Return-Path: <>\r\nReceived: ...", bounces[0].Content, "GetBounces should return correct content")
}

func (s *PostmarkTestSuite) TestGetBounce() {
	responseJSON := `{
	  "RecordType": "Bounce",
	  "ID": 692560173,
	  "Type": "HardBounce",
	  "TypeCode": 1,
	  "Name": "Hard bounce",
	  "Tag": "Invitation",
	  "MessageID": "2c1b63fe-43f2-4db5-91b0-8bdfa44a9316",
	  "MessageStream": "outbound",
	  "Description": "The server was unable to deliver your message (ex: unknown user, mailbox not found).",
	  "Details": "action: failed\r\n",
	  "Email": "anything@blackhole.postmarkap.com",
	  "BouncedAt": "2014-01-15T16:09:19.6421112-05:00",
	  "DumpAvailable": false,
	  "Inactive": false,
	  "CanActivate": true,
	  "Subject": "SC API5 Test",
	  "Content": "Return-Path: <>\r\nReceived: …"
	}`

	s.mux.Get("/bounces/692560173", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetBounce(context.Background(), 692560173)
	s.Require().NoError(err, "GetBounce should not fail")
	s.Equal(int64(692560173), res.ID, "GetBounce should return correct bounce ID")
	s.Equal("Bounce", res.RecordType, "GetBounce should return correct record type")
	s.Equal("outbound", res.MessageStream, "GetBounce should return correct message stream")
	s.Equal("Return-Path: <>\r\nReceived: …", res.Content, "GetBounce should return correct content")
}

func (s *PostmarkTestSuite) TestGetBounceDump() {
	responseJSON := `{
	  "Body": "..."
	}`

	s.mux.Get("/bounces/692560173/dump", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetBounceDump(context.Background(), 692560173)
	s.Require().NoError(err, "GetBounceDump should not fail")
	s.Equal("...", res, "GetBounceDump should return correct dump body")
}

func (s *PostmarkTestSuite) TestActivateBounce() {
	responseJSON := `{
		"Message": "OK",
		"Bounce": {
		  "RecordType": "Bounce",
		  "ID": 692560173,
		  "Type": "HardBounce",
		  "TypeCode": 1,
		  "Name": "Hard bounce",
		  "Tag": "Invitation",
		  "MessageID": "2c1b63fe-43f2-4db5-91b0-8bdfa44a9316",
		  "MessageStream": "outbound",
		  "Description": "The server was unable to deliver your message (ex: unknown user, mailbox not found).",
		  "Details": "action: failed\r\n",
		  "Email": "anything@blackhole.postmarkap.com",
		  "BouncedAt": "2014-01-15T16:09:19.6421112-05:00",
		  "DumpAvailable": false,
		  "Inactive": false,
		  "CanActivate": true,
		  "Subject": "SC API5 Test",
		  "Content": "Return-Path: <>\r\nReceived: …"
		}
	}`

	s.mux.Put("/bounces/692560173/activate", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, mess, err := s.client.ActivateBounce(context.Background(), 692560173)
	s.Require().NoError(err, "ActivateBounce should not fail")
	s.Equal(int64(692560173), res.ID, "ActivateBounce should return correct bounce ID")
	s.Equal("Bounce", res.RecordType, "ActivateBounce should return correct record type")
	s.Equal("outbound", res.MessageStream, "ActivateBounce should return correct message stream")
	s.Equal("Return-Path: <>\r\nReceived: …", res.Content, "ActivateBounce should return correct content")
	s.Equal("OK", mess, "ActivateBounce should return correct message")
}

func (s *PostmarkTestSuite) TestGetBouncedTags() {
	tests := []struct {
		name         string
		responseJSON string
		wantErr      bool
		expectedTags []string
	}{
		{
			name: "successful tags retrieval",
			responseJSON: `[
				"tag1",
				"tag2",
				"tag3"]`,
			wantErr:      false,
			expectedTags: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:         "empty tags list",
			responseJSON: `[]`,
			wantErr:      false,
			expectedTags: []string{},
		},
		{
			name:         "invalid JSON response",
			responseJSON: `invalid json`,
			wantErr:      true,
			expectedTags: []string{},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mux.Get("/bounces/tags", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(tt.responseJSON))
			})

			res, err := s.client.GetBouncedTags(context.Background())

			if tt.wantErr {
				s.Require().Error(err, "GetBouncedTags should fail")
			} else {
				s.Require().NoError(err, "GetBouncedTags should not fail")
				s.Equal(tt.expectedTags, res, "GetBouncedTags should return expected tags")
			}
		})
	}
}

// Benchmark for GetDeliveryStats
func BenchmarkGetDeliveryStats(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
	  "InactiveMails": 192,
	  "Bounces": [
		{
		  "Name": "All",
		  "Count": 253
		}
	  ]
	}`

	mux.Get("/deliverystats", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetDeliveryStats(context.Background())
	}
}

// Benchmark for GetBounces
func BenchmarkGetBounces(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
	  "TotalCount": 1,
	  "Bounces": [
		{
		  "RecordType": "Bounce",
		  "ID": 692560173,
		  "Type": "HardBounce",
		  "Tag": "Invitation",
		  "MessageID": "2c1b63fe-43f2-4db5-91b0-8bdfa44a9316",
		  "MessageStream": "outbound",
		  "Email": "test@example.com"
		}
	  ]
	}`

	mux.Get("/bounces", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	options := map[string]interface{}{
		"tag": "Invitation",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = client.GetBounces(context.Background(), 100, 0, options)
	}
}

// Benchmark for GetBounce
func BenchmarkGetBounce(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
	  "RecordType": "Bounce",
	  "ID": 692560173,
	  "Type": "HardBounce",
	  "Tag": "Invitation",
	  "MessageID": "2c1b63fe-43f2-4db5-91b0-8bdfa44a9316",
	  "MessageStream": "outbound",
	  "Email": "test@example.com"
	}`

	mux.Get("/bounces/692560173", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetBounce(context.Background(), 692560173)
	}
}

// Benchmark for GetBounceDump
func BenchmarkGetBounceDump(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
	  "Body": "test dump content"
	}`

	mux.Get("/bounces/692560173/dump", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetBounceDump(context.Background(), 692560173)
	}
}

// Benchmark for ActivateBounce
func BenchmarkActivateBounce(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
		"Message": "OK",
		"Bounce": {
		  "RecordType": "Bounce",
		  "ID": 692560173,
		  "Type": "HardBounce",
		  "MessageID": "2c1b63fe-43f2-4db5-91b0-8bdfa44a9316",
		  "MessageStream": "outbound",
		  "Email": "test@example.com"
		}
	}`

	mux.Put("/bounces/692560173/activate", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = client.ActivateBounce(context.Background(), 692560173)
	}
}

// Benchmark for GetBouncedTags
func BenchmarkGetBouncedTags(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `["tag1", "tag2", "tag3"]`

	mux.Get("/bounces/tags", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetBouncedTags(context.Background())
	}
}
