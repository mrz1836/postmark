package postmark

import (
	"context"
	"net/http"
)

func (s *PostmarkTestSuite) TestGetCurrentServer() {
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
			"DeliveryHookUrl": "https://hooks.example.com/delivery",
			"InboundAddress": "yourhash@inbound.postmarkapp.com",
			"InboundHookUrl": "https://hooks.example.com/inbound",
			"BounceHookUrl": "https://hooks.example.com/bounce",
			"IncludeBounceContentInHook": true,
			"OpenHookUrl": "https://hooks.example.com/open",
			"PostFirstOpenOnly": false,
			"TrackOpens": false,
			"TrackLinks" : "None",
			"ClickHookUrl" : "https://hooks.example.com/click",
			"InboundDomain": "",
			"InboundHash": "yourhash",
			"InboundSpamThreshold": 0
	}`

	s.mux.Get("/server", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetCurrentServer(context.Background())
	s.Require().NoError(err)

	s.Equal("Staging Testing", res.Name, "GetCurrentServer: wrong name")
}

func (s *PostmarkTestSuite) TestEditCurrentServer() {
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
  "DeliveryHookUrl": "https://hooks.example.com/delivery",
  "InboundAddress": "yourhash@inbound.postmarkapp.com",
  "InboundHookUrl": "https://hooks.example.com/inbound",
  "BounceHookUrl": "https://hooks.example.com/bounce",
  "IncludeBounceContentInHook": true,
  "OpenHookUrl": "https://hooks.example.com/open",
  "PostFirstOpenOnly": false,
  "TrackOpens": false,
  "TrackLinks": "None",
  "ClickHookUrl": "https://hooks.example.com/click",
  "InboundDomain": "",
  "InboundHash": "yourhash",
  "InboundSpamThreshold": 10
}`
	s.mux.Put("/server", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.EditCurrentServer(context.Background(), Server{
		Name: "Production Testing",
	})
	s.Require().NoError(err)

	s.Equal("Production Testing", res.Name, "EditCurrentServer: wrong name")
}
