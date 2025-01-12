package postmark

import (
	"context"
	"net/http"
	"testing"

	"goji.io/pat"
)

func TestGetServers(t *testing.T) {
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

	tMux.HandleFunc(pat.Get("/servers"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := client.GetServers(context.Background(), 100, 10, "")
	if err != nil {
		t.Fatalf("GetServers: %s", err.Error())
	}

	if len(res.Servers) == 0 {
		t.Fatalf("GetServers: unmarshaled to empty")
	}

	if res.TotalCount != 2 {
		t.Fatalf("GetServers: unmarshaled to empty")
	}
}

func TestGetServer(t *testing.T) {
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

	tMux.HandleFunc(pat.Get("/servers/:serverID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := client.GetServer(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetServer: %s", err.Error())
	}

	if res.Name != "Staging Testing" {
		t.Fatalf("GetServer: wrong name!: %s", res.Name)
	}
}

func TestCreateServer(t *testing.T) {
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

	tMux.HandleFunc(pat.Post("/servers"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := client.CreateServer(context.Background(), ServerCreateRequest{
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
	if err != nil {
		t.Fatalf("CreateServer: %s", err.Error())
	}

	if res.Name != "Staging Testing" {
		t.Fatalf("CreateServer: wrong name!")
	}
}

func TestEditServer(t *testing.T) {
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

	tMux.HandleFunc(pat.Put("/servers/:serverID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := client.EditServer(context.Background(), 1234, ServerEditRequest{
		Name: "Production Testing",
	})
	if err != nil {
		t.Fatalf("EditServer: %s", err.Error())
	}

	if res.Name != "Production Testing" {
		t.Fatalf("EditServer: wrong name!: %s", res.Name)
	}
}

func TestDeleteServer(t *testing.T) {
	responseJSON := `{
	  "ErrorCode": 0,
	  "Message": "Server 1234 removed."
	}`

	tMux.HandleFunc(pat.Delete("/servers/:serverID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	// Success
	err := client.DeleteServer(context.Background(), 1234)
	if err != nil {
		t.Fatalf("DeleteServer: %s", err.Error())
	}

	// Failure
	responseJSON = `{
	  "ErrorCode": 402,
	  "Message": "Invalid JSON"
	}`

	err = client.DeleteServer(context.Background(), 1234)
	if err == nil {
		t.Fatalf("DeleteServer: should have failed")
	}
}
