package postmark

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	transactionalDev = "transactional-dev"
)

func (s *PostmarkTestSuite) TestListMessageStreams() {
	responseJSON := `{
		"MessageStreams": [			{
				"ID": "outbound",
				"ServerID": 123457,
				"Name": "Transactional Stream",
				"Description": "This is my stream to send transactional messages",
				"MessageStreamType": "Transactional",
				"CreatedAt": "2020-07-01T00:00:00-04:00",
				"UpdatedAt": "2020-07-05T00:00:00-04:00",
				"ArchivedAt": null,
				"ExpectedPurgeDate": null,
				"SubscriptionManagementConfiguration": {
					"UnsubscribeHandlingType": "none"
				}
			},
			{
				"ID": "inbound",
				"ServerID": 123457,
				"Name": "Inbound Stream",
				"Description": "Stream used for receiving inbound messages",
				"MessageStreamType": "Inbound",
				"CreatedAt": "2020-07-01T00:00:00-04:00",
				"UpdatedAt": null,
				"ArchivedAt": null,
				"ExpectedPurgeDate": null,
				"SubscriptionManagementConfiguration": {
					"UnsubscribeHandlingType": "none"
				}
			},
			{
				"ID": "transactional-dev",
				"ServerID": 123457,
				"Name": "My Dev Transactional Stream",
				"Description": "This is my second transactional stream",
				"MessageStreamType": "Transactional",
				"CreatedAt": "2020-07-02T00:00:00-04:00",
				"UpdatedAt": "2020-07-04T00:00:00-04:00",
				"ArchivedAt": null,
				"ExpectedPurgeDate": null,
				"SubscriptionManagementConfiguration": {
					"UnsubscribeHandlingType": "none"
				}
			}
		],
		"TotalCount": 3
	}`

	s.mux.Get("/message-streams", func(w http.ResponseWriter, req *http.Request) {
		s.Equal("false", req.URL.Query().Get("IncludeArchivedStreams"), "MessageStreams: wrong IncludeArchivedStreams value")
		s.Equal("All", req.URL.Query().Get("MessageStreamType"), "MessageStreams: wrong messageStreamType value")
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.ListMessageStreams(context.Background(), "All", false)
	s.Require().NoError(err)
	s.Len(res, 3, "MessageStreams: wrong number of message streams")

	// For each message stream, check the ServerID
	for _, ms := range res {
		s.Equal(int(123457), ms.ServerID, "MessageStreams: wrong ServerID")
		s.Nil(ms.ArchivedAt, "MessageStreams: ArchivedAt should be nil")
	}

	s.Equal("outbound", res[0].ID, "MessageStreams: wrong ID for first stream")
	s.Equal("inbound", res[1].ID, "MessageStreams: wrong ID for second stream")
	s.Equal(transactionalDev, res[2].ID, "MessageStreams: wrong ID for third stream")
}

func (s *PostmarkTestSuite) TestListMessageStreamsError() {
	// Create a new mux for this specific test to avoid conflicts
	errorMux := NewTestRouter()
	errorServer := httptest.NewServer(errorMux)
	defer errorServer.Close()

	// Create a new client for this test
	errorClient := NewClient("server-token", "account-token")
	errorClient.BaseURL = errorServer.URL

	errorMux.Get("/message-streams", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"ErrorCode": 500, "Message": "Internal Server Error"}`))
	})

	res, err := errorClient.ListMessageStreams(context.Background(), "All", false)
	s.Require().Error(err, "ListMessageStreams should fail")
	s.Nil(res, "ListMessageStreams should return nil on error")
}

func (s *PostmarkTestSuite) TestListMessageStreamsSpecificTypes() {
	// Test each specific message stream type to ensure switch cases are covered
	types := []string{"Inbound", "Transactional", "Broadcasts"}

	for _, msgType := range types {
		s.Run(msgType, func() {
			// Create separate mux/server for each subtest to avoid conflicts
			testMux := NewTestRouter()
			testServer := httptest.NewServer(testMux)
			defer testServer.Close()

			testClient := NewClient("server-token", "account-token")
			testClient.BaseURL = testServer.URL

			responseJSON := `{
				"MessageStreams": [],
				"TotalCount": 0
			}`

			testMux.Get("/message-streams", func(w http.ResponseWriter, req *http.Request) {
				s.Equal(msgType, req.URL.Query().Get("MessageStreamType"))
				_, _ = w.Write([]byte(responseJSON))
			})

			res, err := testClient.ListMessageStreams(context.Background(), msgType, false)
			s.Require().NoError(err)
			s.NotNil(res)
		})
	}
}

func (s *PostmarkTestSuite) TestGetUnknownMessageStream() {
	responseJSON := `{"ErrorCode":1226,"Message":"The message stream for the provided 'ID' was not found."}`

	s.mux.Get("/message-streams/unknown", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetMessageStream(context.Background(), "unknown")
	s.Require().Error(err, "MessageStream: expected error")
	s.Equal("The message stream for the provided 'ID' was not found.", err.Error(), "MessageStream: wrong error message")

	var zero MessageStream
	s.Equal(zero, res, "MessageStream: expected empty response")
}

func (s *PostmarkTestSuite) TestGetMessageStream() {
	responseJSON := `{
		"ID": "broadcasts",
		"ServerID": 123456,
		"Name": "Broadcast Stream",
		"Description": "This is my stream to send broadcast messages",
		"MessageStreamType": "Broadcasts",
		"CreatedAt": "2020-07-01T00:00:00-04:00",
		"UpdatedAt": "2020-07-01T00:00:00-04:00",
		"ArchivedAt": null,
		"ExpectedPurgeDate": null,
		"SubscriptionManagementConfiguration": {
			"UnsubscribeHandlingType": "Postmark"
		}
	}`

	s.mux.Get("/message-streams/broadcasts", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetMessageStream(context.Background(), "broadcasts")
	s.Require().NoError(err)

	s.Equal("broadcasts", res.ID, "MessageStream: wrong ID")
	s.Equal("Broadcast Stream", res.Name, "MessageStream: wrong Name")
	s.Require().NotNil(res.Description, "MessageStream: Description should not be nil")
	s.Equal("This is my stream to send broadcast messages", *res.Description, "MessageStream: wrong Description")
}

func (s *PostmarkTestSuite) TestEditMessageStream() {
	responseJSON := `{
		"ID": "transactional-dev",
		"ServerID": 123457,
		"Name": "Updated Dev Stream",
		"Description": "Updating my dev transactional stream",
		"MessageStreamType": "Transactional",
		"CreatedAt": "2020-07-02T00:00:00-04:00",
		"UpdatedAt": "2020-07-03T00:00:00-04:00",
		"ArchivedAt": null,
		"ExpectedPurgeDate": null,
		"SubscriptionManagementConfiguration": {
			"UnsubscribeHandlingType": "none"
		}
	}`

	editReq := EditMessageStreamRequest{
		Name: "Updated Dev Stream",
		SubscriptionManagementConfiguration: MessageStreamSubscriptionManagementConfiguration{
			UnsubscribeHandlingType: "none",
		},
	}

	s.mux.Patch("/message-streams/transactional-dev", func(w http.ResponseWriter, req *http.Request) {
		var body EditMessageStreamRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		s.NoError(err, "Failed to read request body")

		s.Nil(body.Description, "EditMessageStream: Description should be nil")
		s.Equal(editReq.Name, body.Name, "EditMessageStream: wrong Name")
		s.Equal(editReq.SubscriptionManagementConfiguration.UnsubscribeHandlingType, body.SubscriptionManagementConfiguration.UnsubscribeHandlingType, "EditMessageStream: wrong UnsubscribeHandlingType")

		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.EditMessageStream(context.Background(), transactionalDev, editReq)
	s.Require().NoError(err)

	s.Equal(transactionalDev, res.ID, "MessageStream: wrong ID")
	s.Equal(int(123457), res.ServerID, "MessageStream: wrong ServerID")
	s.Require().NotNil(res.Description, "MessageStream: Description should not be nil")
	s.Equal("Updating my dev transactional stream", *res.Description, "MessageStream: wrong Description")
}

func (s *PostmarkTestSuite) TestCreateMessageStream() {
	responseJSON := `{
		"ID": "transactional-dev",
		"ServerID": 123457,
		"Name": "My Dev Transactional Stream",
		"Description": "This is my second transactional stream",
		"MessageStreamType": "Transactional",
		"CreatedAt": "2020-07-02T00:00:00-04:00",
		"UpdatedAt": "2020-07-02T00:00:00-04:00",
		"ArchivedAt": "2020-07-02T00:00:00-04:00",
		"SubscriptionManagementConfiguration": {
			"UnsubscribeHandlingType": "None"
		}
	}`

	desc := "This is my second transactional stream"
	createReq := CreateMessageStreamRequest{
		ID:                transactionalDev,
		Name:              "My Dev Transactional Stream",
		Description:       &desc,
		MessageStreamType: "Transactional",
		SubscriptionManagementConfiguration: MessageStreamSubscriptionManagementConfiguration{
			UnsubscribeHandlingType: "None",
		},
	}

	s.mux.Post("/message-streams", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.CreateMessageStream(context.Background(), createReq)
	s.Require().NoError(err)

	s.Equal(transactionalDev, res.ID, "MessageStream: wrong ID")
	s.Equal(int(123457), res.ServerID, "MessageStream: wrong ServerID")
	s.Equal(MessageStreamType("Transactional"), res.MessageStreamType, "MessageStream: wrong MessageStreamType")
}

func (s *PostmarkTestSuite) TestArchiveMessageStream() {
	responseJSON := `{
		"ID": "transactional-dev",
		"ServerID": 123457,
		"ExpectedPurgeDate": "2020-08-30T12:30:00.00-04:00"
	}`

	s.mux.Post("/message-streams/transactional-dev/archive", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.ArchiveMessageStream(context.Background(), transactionalDev)
	s.Require().NoError(err)

	s.Equal(transactionalDev, res.ID, "MessageStream: wrong ID")
	s.Equal(int(123457), res.ServerID, "MessageStream: wrong ServerID")
	s.Equal("2020-08-30T12:30:00.00-04:00", res.ExpectedPurgeDate, "MessageStream: wrong ExpectedPurgeDate")
}

func (s *PostmarkTestSuite) TestUnarchiveMessageStream() {
	responseJSON := `{
		"ID": "transactional-dev",
		"ServerID": 123457,
		"Name": "Updated Dev Stream",
		"Description": "Updating my dev transactional stream",
		"MessageStreamType": "Transactional",
		"CreatedAt": "2020-07-02T00:00:00-04:00",
		"UpdatedAt": "2020-07-04T00:00:00-04:00",
		"ArchivedAt": null,
		"SubscriptionManagementConfiguration": {
			"UnsubscribeHandlingType": "none"
		}
	}`

	s.mux.Post("/message-streams/transactional-dev/unarchive", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.UnarchiveMessageStream(context.Background(), transactionalDev)
	s.Require().NoError(err)

	s.Equal(transactionalDev, res.ID, "MessageStream: wrong ID")
	s.Equal(int(123457), res.ServerID, "MessageStream: wrong ServerID")
}

// Benchmarks for Message Streams API

func BenchmarkListMessageStreams(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
		"MessageStreams": [
			{
				"ID": "outbound",
				"ServerID": 123457,
				"Name": "Transactional Stream",
				"MessageStreamType": "Transactional",
				"SubscriptionManagementConfiguration": {
					"UnsubscribeHandlingType": "none"
				}
			}
		],
		"TotalCount": 1
	}`

	mux.Get("/message-streams", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.ListMessageStreams(context.Background(), "All", false)
	}
}

func BenchmarkGetMessageStream(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
		"ID": "transactional-dev",
		"ServerID": 123456,
		"Name": "Dev Stream",
		"MessageStreamType": "Transactional",
		"SubscriptionManagementConfiguration": {
			"UnsubscribeHandlingType": "none"
		}
	}`

	mux.Get("/message-streams/transactional-dev", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetMessageStream(context.Background(), "transactional-dev")
	}
}

func BenchmarkEditMessageStream(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
		"ID": "transactional-dev",
		"ServerID": 123457,
		"Name": "Benchmark Stream",
		"MessageStreamType": "Transactional",
		"SubscriptionManagementConfiguration": {
			"UnsubscribeHandlingType": "none"
		}
	}`

	mux.Patch("/message-streams/transactional-dev", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	editReq := EditMessageStreamRequest{
		Name: "Benchmark Stream",
		SubscriptionManagementConfiguration: MessageStreamSubscriptionManagementConfiguration{
			UnsubscribeHandlingType: NoneUnsubscribeHandlingType,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.EditMessageStream(context.Background(), "transactional-dev", editReq)
	}
}

func BenchmarkCreateMessageStream(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
		"ID": "benchmark-stream",
		"ServerID": 123457,
		"Name": "Benchmark Stream",
		"MessageStreamType": "Transactional",
		"SubscriptionManagementConfiguration": {
			"UnsubscribeHandlingType": "none"
		}
	}`

	mux.Post("/message-streams", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	createReq := CreateMessageStreamRequest{
		ID:                "benchmark-stream",
		Name:              "Benchmark Stream",
		MessageStreamType: TransactionalMessageStreamType,
		SubscriptionManagementConfiguration: MessageStreamSubscriptionManagementConfiguration{
			UnsubscribeHandlingType: NoneUnsubscribeHandlingType,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateMessageStream(context.Background(), createReq)
	}
}

func BenchmarkArchiveMessageStream(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
		"ID": "benchmark-stream",
		"ServerID": 123457,
		"ExpectedPurgeDate": "2020-08-30T12:30:00.00-04:00"
	}`

	mux.Post("/message-streams/benchmark-stream/archive", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.ArchiveMessageStream(context.Background(), "benchmark-stream")
	}
}

func BenchmarkUnarchiveMessageStream(b *testing.B) {
	mux := NewTestRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("server-token", "account-token")
	client.BaseURL = server.URL

	responseJSON := `{
		"ID": "benchmark-stream",
		"ServerID": 123457,
		"Name": "Benchmark Stream",
		"MessageStreamType": "Transactional",
		"SubscriptionManagementConfiguration": {
			"UnsubscribeHandlingType": "none"
		}
	}`

	mux.Post("/message-streams/benchmark-stream/unarchive", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.UnarchiveMessageStream(context.Background(), "benchmark-stream")
	}
}
