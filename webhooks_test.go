package postmark

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

const (
	outbound = "outbound"
)

func (s *PostmarkTestSuite) TestGetWebhooks() {
	responseJSON := `{
		"Webhooks": [
			{
				"ID": 1234567,
				"Url": "http://www.example.com/webhook-test-tracking",
				"MessageStream": "outbound",
				"HttpAuth":{
					"Username": "user",
					"Password": "pass"
				},
				"HttpHeaders":[
					{
						"Name": "name",
						"Value": "value"
					}
				],
				"Triggers": {
					"Open":{
						"Enabled": true,
						"PostFirstOpenOnly": false
					},
					"Click":{
						"Enabled": true
					},
					"Delivery":{
						"Enabled": true
					},
					"Bounce":{
						"Enabled": false,
						"IncludeContent": false
					},
					"SpamComplaint":{
						"Enabled": false,
						"IncludeContent": false
					},
					"SubscriptionChange": {
						"Enabled": true
					}
				}
			},
			{
				"ID": 1234568,
				"Url": "http://www.example.com/webhook-test-bounce",
				"MessageStream": "outbound",
				"HttpAuth":{
					"Username": "user",
					"Password": "pass"
				},
				"HttpHeaders":[
					{
						"Name": "name",
						"Value": "value"
					}
				],
				"Triggers": {
					"Open":{
						"Enabled":false,
						"PostFirstOpenOnly":false
					},
					"Click":{
						"Enabled": false
					},
					"Delivery":{
						"Enabled": false
					},
					"Bounce":{
						"Enabled" :true,
						"IncludeContent": false
					},
					"SpamComplaint":{
						"Enabled": false,
						"IncludeContent": false
					},
					"SubscriptionChange": {
						"Enabled": false
					}
				}
			}
		]
	}`

	s.mux.Get("/webhooks", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.ListWebhooks(context.Background(), "")
	s.Require().NoError(err)

	s.Len(res, 2, "Webhook: wrong number of webhooks listed")
	s.Equal(int(1234567), res[0].ID, "Webhook: wrong first webhook ID")
	s.Equal(int(1234568), res[1].ID, "Webhook: wrong second webhook ID")
}

func (s *PostmarkTestSuite) TestListWebhooksError() {
	// Create a new mux for this specific test to avoid conflicts
	errorMux := NewTestRouter()
	errorServer := httptest.NewServer(errorMux)
	defer errorServer.Close()

	// Create a new client for this test
	errorClient := NewClient("server-token", "account-token")
	errorClient.BaseURL = errorServer.URL

	errorMux.Get("/webhooks", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"ErrorCode": 500, "Message": "Internal Server Error"}`))
	})

	res, err := errorClient.ListWebhooks(context.Background(), "")
	s.Require().Error(err, "ListWebhooks should fail")
	s.Nil(res, "ListWebhooks should return nil on error")
}

func (s *PostmarkTestSuite) TestGetWebhook() {
	responseJSON := `{
		"ID": 1234567,
		"Url": "http://www.example.com/webhook-test-tracking",
		"MessageStream": "outbound",
		"HttpAuth":{
			"Username": "user",
			"Password": "pass"
		},
		"HttpHeaders":[
			{
				"Name": "name",
				"Value": "value"
			}
		],
		"Triggers": {
			"Open":{
				"Enabled": true,
				"PostFirstOpenOnly": false
			},
			"Click":{
				"Enabled": true
			},
			"Delivery":{
				"Enabled": true
			},
			"Bounce":{
				"Enabled": false,
				"IncludeContent": false
			},
			"SpamComplaint":{
				"Enabled": false,
				"IncludeContent": false
			},
			"SubscriptionChange": {
				"Enabled": true
			}
		}
	}`

	s.mux.Get("/webhooks/:webhookID", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetWebhook(context.Background(), 1234567)
	s.Require().NoError(err)

	s.Equal(int(1234567), res.ID, "Webhook: wrong ID")
	s.Equal(outbound, res.MessageStream, "Webhook: wrong message stream")
	s.Equal("name", res.HTTPHeaders[0].Name, "Webhook: wrong HTTP Headers")
	s.True(res.Triggers.SubscriptionChange.Enabled, "Webhook: wrong Subscription Change trigger state")
}

func (s *PostmarkTestSuite) TestCreateWebhook() {
	webhook := Webhook{
		URL:           "http://www.example.com/webhook-test-tracking",
		MessageStream: outbound,
		HTTPAuth: &WebhookHTTPAuth{
			Username: "user",
			Password: "pass",
		},
		HTTPHeaders: []Header{
			{
				Name:  "name",
				Value: "value",
			},
		},
		Triggers: WebhookTrigger{
			Open: WebhookTriggerOpen{
				WebhookTriggerEnabled: WebhookTriggerEnabled{
					Enabled: true,
				},
				PostFirstOpenOnly: true,
			},
			Click: WebhookTriggerEnabled{
				Enabled: true,
			},
		},
	}

	s.mux.Post("/webhooks", func(w http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)

		var res Webhook
		err := decoder.Decode(&res)
		_ = req.Body.Close()

		s.NoError(err, "Failed to decode webhook request")

		s.Equal(outbound, res.MessageStream, "Webhook: wrong message stream")
		s.True(res.Triggers.Open.Enabled, "Webhook: wrong Open trigger state")

		res.ID = 12345

		resBytes, err := json.Marshal(res)
		s.NoError(err, "Failed to marshal webhook response")

		_, _ = w.Write(resBytes)
	})

	res, err := s.client.CreateWebhook(context.Background(), webhook)
	s.Require().NoError(err)

	s.Equal(int(12345), res.ID, "Webhook: wrong ID")
	s.Equal(outbound, res.MessageStream, "Webhook: wrong message stream")
	s.True(res.Triggers.Open.Enabled, "Webhook: wrong Open trigger state")
}

func (s *PostmarkTestSuite) TestEditWebhook() {
	webhook := Webhook{
		URL:           "http://www.example.com/webhook-test-tracking",
		MessageStream: outbound,
		HTTPAuth: &WebhookHTTPAuth{
			Username: "user",
			Password: "pass",
		},
		HTTPHeaders: []Header{
			{
				Name:  "name",
				Value: "value",
			},
		},
		Triggers: WebhookTrigger{
			Open: WebhookTriggerOpen{
				WebhookTriggerEnabled: WebhookTriggerEnabled{
					Enabled: true,
				},
				PostFirstOpenOnly: true,
			},
			Click: WebhookTriggerEnabled{
				Enabled: true,
			},
		},
	}

	s.mux.Put("/webhooks/:webhookID", func(w http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)

		var res Webhook
		err := decoder.Decode(&res)
		_ = req.Body.Close()

		s.NoError(err, "Failed to decode webhook request")

		s.Equal(outbound, res.MessageStream, "Webhook: wrong message stream")
		s.True(res.Triggers.Open.Enabled, "Webhook: wrong Open trigger state")

		res.ID = 12345

		resBytes, err := json.Marshal(res)
		s.NoError(err, "Failed to marshal webhook response")

		_, _ = w.Write(resBytes)
	})

	res, err := s.client.EditWebhook(context.Background(), 12345, webhook)
	s.Require().NoError(err)

	s.Equal(int(12345), res.ID, "Webhook: wrong ID")
	s.Equal(outbound, res.MessageStream, "Webhook: wrong message stream")
	s.True(res.Triggers.Open.Enabled, "Webhook: wrong Open trigger state")
}

func (s *PostmarkTestSuite) TestDeleteWebhook() {
	responseJSON := `{
	  "ErrorCode": 0,
	  "Message": "Webhook 1234 removed."
	}`

	s.mux.Delete("/webhooks/:webhookID", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	// Success
	err := s.client.DeleteWebhook(context.Background(), 1234)
	s.Require().NoError(err)

	// Failure
	responseJSON = `{
	  "ErrorCode": 402,
	  "Message": "Invalid JSON"
	}`

	err = s.client.DeleteWebhook(context.Background(), 1234)
	s.Require().Error(err, "DeleteWebhook: should have failed")
}
