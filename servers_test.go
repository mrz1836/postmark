package postmark

import (
	"context"
	"encoding/json"
	"net/http"

	"goji.io/pat"
)

func (s *PostmarkTestSuite) TestGetServers() {
	responseJSON := `{
  "TotalCount": 2,
  "Servers": [
    {
      "ID": 1,
      "Name": "Production01",
      "ApiTokens": [
        "server token"
      ],
      "Color": "red",
      "SmtpApiActivated": true,
      "RawEmailEnabled": false,
      "DeliveryType": "Live",
      "ServerLink": "https://postmarkapp.com/servers/1/streams",
      "InboundAddress": "yourhash@inbound.postmarkapp.com",
      "InboundHookUrl": "http://inboundhook.example.com/inbound",
      "BounceHookUrl": "http://bouncehook.example.com/bounce",
      "OpenHookUrl": "http://openhook.example.com/open",
      "DeliveryHookUrl": "http://hooks.example.com/delivery",
      "PostFirstOpenOnly": true,
      "InboundDomain": "",
      "InboundHash": "yourhash",
      "InboundSpamThreshold": 5,
      "TrackOpens": false,
      "TrackLinks": "None",
      "IncludeBounceContentInHook": true,
      "ClickHookUrl": "http://hooks.example.com/click",
      "EnableSmtpApiErrorHooks": false
    },
    {
      "ID": 2,
      "Name": "Production02",
      "ApiTokens": [
        "server token"
      ],
      "Color": "green",
      "SmtpApiActivated": true,
      "RawEmailEnabled": false,
      "DeliveryType": "Sandbox",
      "ServerLink": "https://postmarkapp.com/servers/2/streams",
      "InboundAddress": "yourhash@inbound.postmarkapp.com",
      "InboundHookUrl": "",
      "BounceHookUrl": "",
      "OpenHookUrl": "",
      "DeliveryHookUrl": "http://hooks.example.com/delivery",
      "PostFirstOpenOnly": false,
      "InboundDomain": "",
      "InboundHash": "yourhash",
      "InboundSpamThreshold": 0,
      "TrackOpens": true,
      "TrackLinks": "HtmlAndText",
      "IncludeBounceContentInHook": false,
      "ClickHookUrl": "",
      "EnableSmtpApiErrorHooks": false
    }
  ]
}`

	s.mux.HandleFunc(pat.Get("/servers"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetServers(context.Background(), 100, 10, "")
	s.Require().NoError(err)

	s.NotEmpty(res.Servers, "GetServers: servers should not be empty")
	s.Equal(int(2), res.TotalCount, "GetServers: wrong total count")
}

func (s *PostmarkTestSuite) TestGetServer() {
	responseJSON := `{
	  "ID": 1,
	  "Name": "Staging Testing",
	  "ApiTokens": [
		"server token"
	  ],
	  "ServerLink": "https://postmarkapp.com/servers/1/overview",
	  "Color": "red",
	  "SmtpApiActivated": true,
	  "RawEmailEnabled": false,
	  "InboundAddress": "yourhash@inbound.postmarkapp.com",
	  "InboundHookUrl": "https://hooks.example.com/inbound",
	  "BounceHookUrl": "https://hooks.example.com/bounce",
	  "OpenHookUrl": "https://hooks.example.com/open",
	  "PostFirstOpenOnly": false,
	  "TrackOpens": false,
	  "InboundDomain": "",
	  "InboundHash": "yourhash",
	  "InboundSpamThreshold": 0
	}`

	s.mux.HandleFunc(pat.Get("/servers/:serverID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetServer(context.Background(), 1)
	s.Require().NoError(err)

	s.Equal("Staging Testing", res.Name, "GetServer: wrong name")
}

func (s *PostmarkTestSuite) TestCreateServer() {
	responseJSON := `{
  "ID": 1,
  "Name": "Staging Testing",
  "ApiTokens": [
    "server token"
  ],
  "Color": "red",
  "SmtpApiActivated": true,
  "RawEmailEnabled": false,
  "DeliveryType": "Live",
  "ServerLink": "https://postmarkapp.com/servers/1/streams",
  "InboundAddress": "yourhash@inbound.postmarkapp.com",
  "InboundHookUrl": "http://hooks.example.com/inbound",
  "PostFirstOpenOnly": false,
  "InboundDomain": "",
  "InboundHash": "yourhash",
  "InboundSpamThreshold": 5,
  "TrackOpens": false,
  "TrackLinks": "None",
  "IncludeBounceContentInHook": true,
  "EnableSmtpApiErrorHooks": false
}`

	s.mux.HandleFunc(pat.Post("/servers"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.CreateServer(context.Background(), ServerCreateRequest{
		Name:                       "Staging Testing",
		Color:                      "red",
		SMTPAPIActivated:           true,
		RawEmailEnabled:            false,
		InboundHookURL:             "http://hooks.example.com/inbound",
		PostFirstOpenOnly:          false,
		InboundDomain:              "",
		InboundSpamThreshold:       5,
		TrackOpens:                 false,
		TrackLinks:                 "None",
		IncludeBounceContentInHook: true,
		EnableSMTPAPIErrorHooks:    false,
	})
	s.Require().NoError(err)

	s.Equal("Staging Testing", res.Name, "CreateServer: wrong name")
}

func (s *PostmarkTestSuite) TestEditServer() {
	responseJSON := `{
	  "ID": 1,
	  "Name": "Production Testing",
	  "ApiTokens": [
		"Server Token"
	  ],
	  "ServerLink": "https://postmarkapp.com/servers/1/overview",
	  "Color": "blue",
	  "SmtpApiActivated": false,
	  "RawEmailEnabled": false,
	  "InboundAddress": "yourhash@inbound.postmarkapp.com",
	  "InboundHookUrl": "https://hooks.example.com/inbound",
	  "BounceHookUrl": "https://hooks.example.com/bounce",
	  "OpenHookUrl": "https://hooks.example.com/open",
	  "PostFirstOpenOnly": false,
	  "TrackOpens": false,
	  "InboundDomain": "",
	  "InboundHash": "yourhash",
	  "InboundSpamThreshold": 10
	}`

	s.mux.HandleFunc(pat.Put("/servers/:serverID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.EditServer(context.Background(), 1234, ServerEditRequest{
		Name: "Production Testing",
	})
	s.Require().NoError(err)

	s.Equal("Production Testing", res.Name, "EditServer: wrong name")
}

func (s *PostmarkTestSuite) TestDeleteServer() {
	responseJSON := `{
	  "ErrorCode": 0,
	  "Message": "Server 1234 removed."
	}`

	s.mux.HandleFunc(pat.Delete("/servers/:serverID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	// Success
	err := s.client.DeleteServer(context.Background(), 1234)
	s.Require().NoError(err)

	// Failure
	responseJSON = `{
	  "ErrorCode": 402,
	  "Message": "Invalid JSON"
	}`

	err = s.client.DeleteServer(context.Background(), 1234)
	s.Require().Error(err, "DeleteServer: should have failed")
}

func (s *PostmarkTestSuite) TestServerMarshalJSON() {
	s.Run("sets default values when empty", func() {
		server := Server{
			ID:   123,
			Name: "My Server",
			// TrackLinks and DeliveryType are empty
		}

		data, err := json.Marshal(server)
		s.Require().NoError(err, "unexpected error during marshal")

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		s.Require().NoError(err, "unexpected error during unmarshal")

		s.Equal("None", result["TrackLinks"], "expected TrackLinks to be 'None'")
		s.Equal("Live", result["DeliveryType"], "expected DeliveryType to be 'Live'")
	})

	s.Run("preserves existing values", func() {
		server := Server{
			ID:           456,
			Name:         "Another Server",
			TrackLinks:   "HtmlOnly",
			DeliveryType: "Sandbox",
		}

		data, err := json.Marshal(server)
		s.Require().NoError(err, "unexpected error during marshal")

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		s.Require().NoError(err, "unexpected error during unmarshal")

		s.Equal("HtmlOnly", result["TrackLinks"], "expected TrackLinks to be 'HtmlOnly'")
		s.Equal("Sandbox", result["DeliveryType"], "expected DeliveryType to be 'Sandbox'")
	})
}
