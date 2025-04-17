package postmark

import "time"

// BaseEvent contains fields that are common across all webhook event types
type BaseEvent struct {
	RecordType    string                 `json:"RecordType"`
	MessageID     string                 `json:"MessageID"`
	MessageStream string                 `json:"MessageStream"`
	Metadata      map[string]interface{} `json:"Metadata,omitempty"`
	Tag           string                 `json:"Tag,omitempty"`
}

// DeliveryEvent represents a successful email delivery webhook event
type DeliveryEvent struct {
	BaseEvent
	ServerID    int       `json:"ServerID"`
	Recipient   string    `json:"Recipient"`
	DeliveredAt time.Time `json:"DeliveredAt"`
	Details     string    `json:"Details,omitempty"`
}

// OpenEvent represents an email open webhook event
type OpenEvent struct {
	BaseEvent
	FirstOpen   bool       `json:"FirstOpen"`
	Recipient   string     `json:"Recipient"`
	ReceivedAt  time.Time  `json:"ReceivedAt"`
	Platform    string     `json:"Platform"`
	ReadSeconds int        `json:"ReadSeconds"`
	UserAgent   string     `json:"UserAgent"`
	OS          OSInfo     `json:"OS"`
	Client      ClientInfo `json:"Client"`
	Geo         GeoInfo    `json:"Geo"`
}

// ClickEvent represents a link click webhook event
type ClickEvent struct {
	BaseEvent
	Recipient     string     `json:"Recipient"`
	ReceivedAt    time.Time  `json:"ReceivedAt"`
	Platform      string     `json:"Platform"`
	ClickLocation string     `json:"ClickLocation"`
	OriginalLink  string     `json:"OriginalLink"`
	UserAgent     string     `json:"UserAgent"`
	OS            OSInfo     `json:"OS"`
	Client        ClientInfo `json:"Client"`
	Geo           GeoInfo    `json:"Geo"`
}

// BounceEvent represents an email bounce webhook event
type BounceEvent struct {
	BaseEvent
	ID            int       `json:"ID"`
	Type          string    `json:"Type"`
	TypeCode      int       `json:"TypeCode"`
	Name          string    `json:"Name"`
	ServerID      int       `json:"ServerID"`
	Description   string    `json:"Description"`
	Details       string    `json:"Details,omitempty"`
	Email         string    `json:"Email"`
	From          string    `json:"From"`
	BouncedAt     time.Time `json:"BouncedAt"`
	DumpAvailable bool      `json:"DumpAvailable"`
	Inactive      bool      `json:"Inactive"`
	CanActivate   bool      `json:"CanActivate"`
	Subject       string    `json:"Subject"`
	Content       string    `json:"Content"`
}

// SpamComplaintEvent represents a spam complaint webhook event
type SpamComplaintEvent struct {
	BaseEvent
	ID            int       `json:"ID"`
	Type          string    `json:"Type"`
	TypeCode      int       `json:"TypeCode"`
	Name          string    `json:"Name"`
	ServerID      int       `json:"ServerID"`
	Description   string    `json:"Description"`
	Details       string    `json:"Details,omitempty"`
	Email         string    `json:"Email"`
	From          string    `json:"From"`
	BouncedAt     time.Time `json:"BouncedAt"`
	DumpAvailable bool      `json:"DumpAvailable"`
	Inactive      bool      `json:"Inactive"`
	CanActivate   bool      `json:"CanActivate"`
	Subject       string    `json:"Subject"`
	Content       string    `json:"Content"`
}

// SubscriptionChangeEvent represents a subscription change webhook event
type SubscriptionChangeEvent struct {
	BaseEvent
	ServerID          int       `json:"ServerID"`
	ChangedAt         time.Time `json:"ChangedAt"`
	Recipient         string    `json:"Recipient"`
	Origin            string    `json:"Origin"`
	SuppressSending   bool      `json:"SuppressSending"`
	SuppressionReason string    `json:"SuppressionReason,omitempty"`
}

// Common nested structures

// OSInfo contains operating system information
type OSInfo struct {
	Name    string `json:"Name"`
	Family  string `json:"Family"`
	Company string `json:"Company"`
}

// ClientInfo contains email client information
type ClientInfo struct {
	Name    string `json:"Name"`
	Family  string `json:"Family"`
	Company string `json:"Company"`
}

// GeoInfo contains geographical information
type GeoInfo struct {
	IP             string `json:"IP"`
	City           string `json:"City"`
	Country        string `json:"Country"`
	CountryISOCode string `json:"CountryISOCode"`
	Region         string `json:"Region"`
	RegionISOCode  string `json:"RegionISOCode"`
	Zip            string `json:"Zip"`
	Coords         string `json:"Coords"`
}
