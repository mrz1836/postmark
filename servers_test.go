package postmark

import (
	"context"
	"net/http"
	"testing"

	"goji.io/pat"
)

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

	res, err := client.GetServer(context.Background(), "1")
	if err != nil {
		t.Fatalf("GetServer: %s", err.Error())
	}

	if res.Name != "Staging Testing" {
		t.Fatalf("GetServer: wrong name!: %s", res.Name)
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

	res, err := client.EditServer(context.Background(), "1234", ServerEditRequest{
		Name: "Production Testing",
	})
	if err != nil {
		t.Fatalf("EditServer: %s", err.Error())
	}

	if res.Name != "Production Testing" {
		t.Fatalf("EditServer: wrong name!: %s", res.Name)
	}
}
