package postmark

import (
	"context"
	"net/http"

	"goji.io/pat"
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

	s.mux.HandleFunc(pat.Get("/deliverystats"), func(w http.ResponseWriter, _ *http.Request) {
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
		  "ID": 692560173,
		  "Type": "HardBounce",
		  "TypeCode": 1,
		  "Name": "Hard bounce",
		  "Tag": "Invitation",
		  "MessageID": "2c1b63fe-43f2-4db5-91b0-8bdfa44a9316",
		  "Description": "The server was unable to deliver your message (ex: unknown user, mailbox not found).",
		  "Details": "action: failed\r\n",
		  "Email": "anything@blackhole.postmarkap.com",
		  "BouncedAt": "2014-01-15T16:09:19.6421112-05:00",
		  "DumpAvailable": false,
		  "Inactive": false,
		  "CanActivate": true,
		  "Subject": "SC API5 Test"
		},
		{
		  "ID": 676862817,
		  "Type": "HardBounce",
		  "TypeCode": 1,
		  "Name": "Hard bounce",
		  "Tag": "Invitation",
		  "MessageID": "623b2e90-82d0-4050-ae9e-2c3a734ba091",
		  "Description": "The server was unable to deliver your message (ex: unknown user, mailbox not found).",
		  "Details": "smtp;554 delivery error: dd This user doesn't have a yahoo.com account (vicelcown@yahoo.com) [0] - mta1543.mail.ne1.yahoo.com",
		  "Email": "vicelcown@yahoo.com",
		  "BouncedAt": "2013-10-18T09:49:59.8253577-04:00",
		  "DumpAvailable": false,
		  "Inactive": true,
		  "CanActivate": true,
		  "Subject": "Production API Test"
		}
		  ]
	}`

	s.mux.HandleFunc(pat.Get("/bounces"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	_, total, err := s.client.GetBounces(context.Background(), 100, 0, map[string]interface{}{
		"tag": "Invitation",
	})
	s.Require().NoError(err, "GetBounces should not fail")
	s.Equal(int64(253), total, "GetBounces should return correct total count")
}

func (s *PostmarkTestSuite) TestGetBounce() {
	responseJSON := `{
	  "ID": 692560173,
	  "Type": "HardBounce",
	  "TypeCode": 1,
	  "Name": "Hard bounce",
	  "Tag": "Invitation",
	  "MessageID": "2c1b63fe-43f2-4db5-91b0-8bdfa44a9316",
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

	s.mux.HandleFunc(pat.Get("/bounces/692560173"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetBounce(context.Background(), 692560173)
	s.Require().NoError(err, "GetBounce should not fail")
	s.Equal(int64(692560173), res.ID, "GetBounce should return correct bounce ID")
}

func (s *PostmarkTestSuite) TestGetBounceDump() {
	responseJSON := `{
	  "Body": "..."
	}`

	s.mux.HandleFunc(pat.Get("/bounces/692560173/dump"), func(w http.ResponseWriter, _ *http.Request) {
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
		  "ID": 692560173,
		  "Type": "HardBounce",
		  "TypeCode": 1,
		  "Name": "Hard bounce",
		  "Tag": "Invitation",
		  "MessageID": "2c1b63fe-43f2-4db5-91b0-8bdfa44a9316",
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

	s.mux.HandleFunc(pat.Put("/bounces/692560173/activate"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, mess, err := s.client.ActivateBounce(context.Background(), 692560173)
	s.Require().NoError(err, "ActivateBounce should not fail")
	s.Equal(int64(692560173), res.ID, "ActivateBounce should return correct bounce ID")
	s.Equal("OK", mess, "ActivateBounce should return correct message")
}

func (s *PostmarkTestSuite) TestGetBouncedTags() {
	responseJSON := `[
		"tag1",
		"tag2",
		"tag3"]
	`

	s.mux.HandleFunc(pat.Get("/bounces/tags"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetBouncedTags(context.Background())
	s.Require().NoError(err, "GetBouncedTags should not fail")
	s.Len(res, 3, "GetBouncedTags should return 3 tags")
}
