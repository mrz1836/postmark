package postmark

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeliveryEvent_Marshal(t *testing.T) {
	event := DeliveryEvent{
		BaseEvent: BaseEvent{
			RecordType:    "Delivery",
			MessageID:     "883953f4-6105-42a2-a16a-77a8eac79483",
			MessageStream: "outbound",
			Tag:           "welcome-email",
			Metadata: map[string]interface{}{
				"example": "value",
			},
		},
		ServerID:    12345,
		Recipient:   "john@example.com",
		DeliveredAt: time.Date(2014, 4, 1, 13, 42, 10, 0, time.UTC),
		Details:     "Test delivery details",
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)
	assert.Contains(t, string(data), "Delivery")
	assert.Contains(t, string(data), "john@example.com")
}

func TestDeliveryEvent_Unmarshal(t *testing.T) {
	jsonData := `{
		"RecordType": "Delivery",
		"MessageID": "883953f4-6105-42a2-a16a-77a8eac79483",
		"MessageStream": "outbound",
		"ServerID": 12345,
		"Recipient": "john@example.com",
		"DeliveredAt": "2014-04-01T13:42:10Z",
		"Details": "Test delivery details",
		"Tag": "welcome-email",
		"Metadata": {
			"example": "value"
		}
	}`

	var event DeliveryEvent
	err := json.Unmarshal([]byte(jsonData), &event)
	require.NoError(t, err)

	assert.Equal(t, "Delivery", event.RecordType)
	assert.Equal(t, "883953f4-6105-42a2-a16a-77a8eac79483", event.MessageID)
	assert.Equal(t, "outbound", event.MessageStream)
	assert.Equal(t, 12345, event.ServerID)
	assert.Equal(t, "john@example.com", event.Recipient)
	assert.Equal(t, "welcome-email", event.Tag)
	assert.NotNil(t, event.Metadata)
}

func TestOpenEvent_Marshal(t *testing.T) {
	event := OpenEvent{
		BaseEvent: BaseEvent{
			RecordType:    "Open",
			MessageID:     "883953f4-6105-42a2-a16a-77a8eac79483",
			MessageStream: "outbound",
		},
		FirstOpen:   true,
		Recipient:   "john@example.com",
		ReceivedAt:  time.Date(2014, 4, 1, 13, 42, 10, 0, time.UTC),
		Platform:    "Desktop",
		ReadSeconds: 5,
		UserAgent:   "Mozilla/5.0",
		OS: OSInfo{
			Name:    "macOS",
			Family:  "OS X",
			Company: "Apple Inc.",
		},
		Client: ClientInfo{
			Name:    "Chrome",
			Family:  "Chrome",
			Company: "Google Inc.",
		},
		Geo: GeoInfo{
			IP:             "192.168.1.1",
			City:           "San Francisco",
			Country:        "United States",
			CountryISOCode: "US",
			Region:         "California",
			RegionISOCode:  "CA",
			Zip:            "94102",
			Coords:         "37.7749,-122.4194",
		},
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)
	assert.Contains(t, string(data), "Open")
	assert.Contains(t, string(data), "john@example.com")
	assert.Contains(t, string(data), "macOS")
}

func TestOpenEvent_Unmarshal(t *testing.T) {
	jsonData := `{
		"RecordType": "Open",
		"MessageID": "883953f4-6105-42a2-a16a-77a8eac79483",
		"MessageStream": "outbound",
		"FirstOpen": true,
		"Recipient": "john@example.com",
		"ReceivedAt": "2014-04-01T13:42:10Z",
		"Platform": "Desktop",
		"ReadSeconds": 5,
		"UserAgent": "Mozilla/5.0",
		"OS": {
			"Name": "macOS",
			"Family": "OS X",
			"Company": "Apple Inc."
		},
		"Client": {
			"Name": "Chrome",
			"Family": "Chrome",
			"Company": "Google Inc."
		},
		"Geo": {
			"IP": "192.168.1.1",
			"City": "San Francisco",
			"Country": "United States",
			"CountryISOCode": "US",
			"Region": "California",
			"RegionISOCode": "CA",
			"Zip": "94102",
			"Coords": "37.7749,-122.4194"
		}
	}`

	var event OpenEvent
	err := json.Unmarshal([]byte(jsonData), &event)
	require.NoError(t, err)

	assert.Equal(t, "Open", event.RecordType)
	assert.True(t, event.FirstOpen)
	assert.Equal(t, "john@example.com", event.Recipient)
	assert.Equal(t, "macOS", event.OS.Name)
	assert.Equal(t, "Chrome", event.Client.Name)
	assert.Equal(t, "San Francisco", event.Geo.City)
}

func TestClickEvent_Marshal(t *testing.T) {
	event := ClickEvent{
		BaseEvent: BaseEvent{
			RecordType:    "Click",
			MessageID:     "883953f4-6105-42a2-a16a-77a8eac79483",
			MessageStream: "outbound",
		},
		Recipient:     "john@example.com",
		ReceivedAt:    time.Date(2014, 4, 1, 13, 42, 10, 0, time.UTC),
		Platform:      "Desktop",
		ClickLocation: "HTML body",
		OriginalLink:  "https://example.com",
		UserAgent:     "Mozilla/5.0",
		OS: OSInfo{
			Name: "Windows",
		},
		Client: ClientInfo{
			Name: "Chrome",
		},
		Geo: GeoInfo{
			IP:      "192.168.1.1",
			Country: "United States",
		},
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)
	assert.Contains(t, string(data), "Click")
	assert.Contains(t, string(data), "https://example.com")
}

func TestClickEvent_Unmarshal(t *testing.T) {
	jsonData := `{
		"RecordType": "Click",
		"MessageID": "883953f4-6105-42a2-a16a-77a8eac79483",
		"MessageStream": "outbound",
		"Recipient": "john@example.com",
		"ReceivedAt": "2014-04-01T13:42:10Z",
		"Platform": "Desktop",
		"ClickLocation": "HTML body",
		"OriginalLink": "https://example.com",
		"UserAgent": "Mozilla/5.0",
		"OS": {
			"Name": "Windows"
		},
		"Client": {
			"Name": "Chrome"
		},
		"Geo": {
			"IP": "192.168.1.1",
			"Country": "United States"
		}
	}`

	var event ClickEvent
	err := json.Unmarshal([]byte(jsonData), &event)
	require.NoError(t, err)

	assert.Equal(t, "Click", event.RecordType)
	assert.Equal(t, "john@example.com", event.Recipient)
	assert.Equal(t, "https://example.com", event.OriginalLink)
	assert.Equal(t, "HTML body", event.ClickLocation)
}

func TestBounceEvent_Marshal(t *testing.T) {
	event := BounceEvent{
		BaseEvent: BaseEvent{
			RecordType:    "Bounce",
			MessageID:     "883953f4-6105-42a2-a16a-77a8eac79483",
			MessageStream: "outbound",
		},
		ID:            42,
		Type:          "HardBounce",
		TypeCode:      1,
		Name:          "Hard bounce",
		ServerID:      12345,
		Description:   "The server was unable to deliver your message",
		Details:       "Test bounce details",
		Email:         "john@example.com",
		From:          "sender@example.com",
		BouncedAt:     time.Date(2014, 4, 1, 13, 42, 10, 0, time.UTC),
		DumpAvailable: true,
		Inactive:      true,
		CanActivate:   false,
		Subject:       "Test Subject",
		Content:       "Test content",
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)
	assert.Contains(t, string(data), "Bounce")
	assert.Contains(t, string(data), "HardBounce")
}

func TestBounceEvent_Unmarshal(t *testing.T) {
	jsonData := `{
		"RecordType": "Bounce",
		"MessageID": "883953f4-6105-42a2-a16a-77a8eac79483",
		"MessageStream": "outbound",
		"ID": 42,
		"Type": "HardBounce",
		"TypeCode": 1,
		"Name": "Hard bounce",
		"ServerID": 12345,
		"Description": "The server was unable to deliver your message",
		"Details": "Test bounce details",
		"Email": "john@example.com",
		"From": "sender@example.com",
		"BouncedAt": "2014-04-01T13:42:10Z",
		"DumpAvailable": true,
		"Inactive": true,
		"CanActivate": false,
		"Subject": "Test Subject",
		"Content": "Test content"
	}`

	var event BounceEvent
	err := json.Unmarshal([]byte(jsonData), &event)
	require.NoError(t, err)

	assert.Equal(t, "Bounce", event.RecordType)
	assert.Equal(t, 42, event.ID)
	assert.Equal(t, "HardBounce", event.Type)
	assert.Equal(t, "john@example.com", event.Email)
	assert.True(t, event.DumpAvailable)
}

func TestSpamComplaintEvent_Marshal(t *testing.T) {
	event := SpamComplaintEvent{
		BaseEvent: BaseEvent{
			RecordType:    "SpamComplaint",
			MessageID:     "883953f4-6105-42a2-a16a-77a8eac79483",
			MessageStream: "outbound",
		},
		ID:            42,
		Type:          "SpamComplaint",
		TypeCode:      512,
		Name:          "Spam complaint",
		ServerID:      12345,
		Description:   "The subscriber reported this message as spam",
		Email:         "john@example.com",
		From:          "sender@example.com",
		BouncedAt:     time.Date(2014, 4, 1, 13, 42, 10, 0, time.UTC),
		DumpAvailable: false,
		Inactive:      true,
		CanActivate:   true,
		Subject:       "Test Subject",
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)
	assert.Contains(t, string(data), "SpamComplaint")
}

func TestSpamComplaintEvent_Unmarshal(t *testing.T) {
	jsonData := `{
		"RecordType": "SpamComplaint",
		"MessageID": "883953f4-6105-42a2-a16a-77a8eac79483",
		"MessageStream": "outbound",
		"ID": 42,
		"Type": "SpamComplaint",
		"TypeCode": 512,
		"Name": "Spam complaint",
		"ServerID": 12345,
		"Description": "The subscriber reported this message as spam",
		"Email": "john@example.com",
		"From": "sender@example.com",
		"BouncedAt": "2014-04-01T13:42:10Z",
		"DumpAvailable": false,
		"Inactive": true,
		"CanActivate": true,
		"Subject": "Test Subject"
	}`

	var event SpamComplaintEvent
	err := json.Unmarshal([]byte(jsonData), &event)
	require.NoError(t, err)

	assert.Equal(t, "SpamComplaint", event.RecordType)
	assert.Equal(t, 42, event.ID)
	assert.Equal(t, "john@example.com", event.Email)
	assert.True(t, event.Inactive)
	assert.True(t, event.CanActivate)
}

func TestSubscriptionChangeEvent_Marshal(t *testing.T) {
	event := SubscriptionChangeEvent{
		BaseEvent: BaseEvent{
			RecordType:    "SubscriptionChange",
			MessageID:     "883953f4-6105-42a2-a16a-77a8eac79483",
			MessageStream: "outbound",
		},
		ServerID:          12345,
		ChangedAt:         time.Date(2014, 4, 1, 13, 42, 10, 0, time.UTC),
		Recipient:         "john@example.com",
		Origin:            "Recipient",
		SuppressSending:   true,
		SuppressionReason: "ManualSuppression",
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)
	assert.Contains(t, string(data), "SubscriptionChange")
	assert.Contains(t, string(data), "ManualSuppression")
}

func TestSubscriptionChangeEvent_Unmarshal(t *testing.T) {
	jsonData := `{
		"RecordType": "SubscriptionChange",
		"MessageID": "883953f4-6105-42a2-a16a-77a8eac79483",
		"MessageStream": "outbound",
		"ServerID": 12345,
		"ChangedAt": "2014-04-01T13:42:10Z",
		"Recipient": "john@example.com",
		"Origin": "Recipient",
		"SuppressSending": true,
		"SuppressionReason": "ManualSuppression"
	}`

	var event SubscriptionChangeEvent
	err := json.Unmarshal([]byte(jsonData), &event)
	require.NoError(t, err)

	assert.Equal(t, "SubscriptionChange", event.RecordType)
	assert.Equal(t, 12345, event.ServerID)
	assert.Equal(t, "john@example.com", event.Recipient)
	assert.True(t, event.SuppressSending)
	assert.Equal(t, "ManualSuppression", event.SuppressionReason)
}

func TestBaseEvent_WithMetadata(t *testing.T) {
	event := BaseEvent{
		RecordType:    "Delivery",
		MessageID:     "test-id",
		MessageStream: "outbound",
		Tag:           "test-tag",
		Metadata: map[string]interface{}{
			"user_id":     "12345",
			"campaign_id": "summer_2024",
			"version":     2,
		},
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded BaseEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "test-id", decoded.MessageID)
	assert.Equal(t, "test-tag", decoded.Tag)
	assert.NotNil(t, decoded.Metadata)
	assert.Equal(t, "12345", decoded.Metadata["user_id"])
}

func TestOSInfo_Marshal(t *testing.T) {
	os := OSInfo{
		Name:    "macOS",
		Family:  "OS X",
		Company: "Apple Inc.",
	}

	data, err := json.Marshal(os)
	require.NoError(t, err)
	assert.Contains(t, string(data), "macOS")
	assert.Contains(t, string(data), "Apple Inc.")
}

func TestClientInfo_Marshal(t *testing.T) {
	client := ClientInfo{
		Name:    "Chrome",
		Family:  "Chrome",
		Company: "Google Inc.",
	}

	data, err := json.Marshal(client)
	require.NoError(t, err)
	assert.Contains(t, string(data), "Chrome")
	assert.Contains(t, string(data), "Google Inc.")
}

func TestGeoInfo_Marshal(t *testing.T) {
	geo := GeoInfo{
		IP:             "192.168.1.1",
		City:           "San Francisco",
		Country:        "United States",
		CountryISOCode: "US",
		Region:         "California",
		RegionISOCode:  "CA",
		Zip:            "94102",
		Coords:         "37.7749,-122.4194",
	}

	data, err := json.Marshal(geo)
	require.NoError(t, err)
	assert.Contains(t, string(data), "San Francisco")
	assert.Contains(t, string(data), "US")
	assert.Contains(t, string(data), "37.7749,-122.4194")
}
