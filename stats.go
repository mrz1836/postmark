package postmark

import (
	"context"
)

// OutboundStats - a brief overview of statistics for all of your outbound email.
type OutboundStats struct {
	// Sent - Number of sent emails
	Sent int64
	// Bounced - Number of bounced emails
	Bounced int64
	// SMTPApiErrors - Number of SMTP errors
	SMTPApiErrors int64
	// BounceRate - Bounce rate percentage calculated by total sent.
	BounceRate float64
	// SpamComplaints - Number of spam complaints received
	SpamComplaints int64
	// SpamComplaintsRate - Spam complaints percentage calculated by total sent.
	SpamComplaintsRate float64
	// Opens - Number of opens
	Opens int64
	// UniqueOpens - Number of unique opens
	UniqueOpens int64
	// Tracked - Number of tracked emails sent
	Tracked int64
	// WithClientRecorded - Number of emails where the client was successfully tracked.
	WithClientRecorded int64
	// WithPlatformRecorded - Number of emails where platform was successfully tracked.
	WithPlatformRecorded int64
	// WithReadTimeRecorded - Number of emails where read time was successfully tracked.
	WithReadTimeRecorded int64
}

// GetOutboundStats - Gets a brief overview of statistics for all of your outbound email.
// Available options: http://developer.postmarkapp.com/developer-api-stats.html#overview
func (client *Client) GetOutboundStats(ctx context.Context, options map[string]interface{}) (OutboundStats, error) {
	res := OutboundStats{}
	err := client.get(ctx, buildURL("stats/outbound", options), &res)
	return res, err
}

// SendDay - send stats for a specific day
type SendDay struct {
	// Date - self-explanatory
	Date string
	// Sent - number of emails sent
	Sent int64
}

// SendCounts - send stats for a period
type SendCounts struct {
	// Days - List of objects that each represent sent counts by date.
	Days []SendDay
	// Sent - Indicates the number of total sent emails returned.
	Sent int64
}

// GetSentCounts - Gets a total count of emails you’ve sent out.
// Available options: http://developer.postmarkapp.com/developer-api-stats.html#sent-counts
func (client *Client) GetSentCounts(ctx context.Context, options map[string]interface{}) (SendCounts, error) {
	res := SendCounts{}
	err := client.get(ctx, buildURL("stats/outbound/sends", options), &res)
	return res, err
}

// BounceDay - bounce stats for a specific day
type BounceDay struct {
	// Date - self-explanatory
	Date string
	// HardBounce - number of hard bounces
	HardBounce int64
	// SoftBounce - number of soft bounces
	SoftBounce int64
	// SMTPApiError - number of SMTP errors
	SMTPApiError int64
	// Transient - number of transient bounces.
	Transient int64
}

// BounceCounts - bounce stats for a period
type BounceCounts struct {
	// Days - List of objects that each represent sent counts by date.
	Days []BounceDay
	// HardBounce - total number of hard bounces
	HardBounce int64
	// SoftBounce - total number of soft bounces
	SoftBounce int64
	// SMTPApiError - total number of SMTP errors
	SMTPApiError int64
	// Transient - total number of transient bounces.
	Transient int64
}

// GetBounceCounts - Gets total counts of emails you’ve sent out that have been returned as bounced.
// Available options: http://developer.postmarkapp.com/developer-api-stats.html#bounce-counts
func (client *Client) GetBounceCounts(ctx context.Context, options map[string]interface{}) (BounceCounts, error) {
	res := BounceCounts{}
	err := client.get(ctx, buildURL("stats/outbound/bounces", options), &res)
	return res, err
}

// SpamDay - spam complaints for a specific day
type SpamDay struct {
	// Date - self-explanatory
	Date string
	// SpamComplaint - number of spam complaints received
	SpamComplaint int64
}

// SpamCounts - spam complaints for a period
type SpamCounts struct {
	// Days - List of objects that each represent spam complaint counts by date.
	Days []SpamDay
	// SpamComplaint - Indicates total number of spam complaints.
	SpamComplaint int64
}

// GetSpamCounts - Gets a total count of recipients who have marked your email as spam.
// Days that did not produce statistics won’t appear in the JSON response.
// Available options: http://developer.postmarkapp.com/developer-api-stats.html#spam-complaints
func (client *Client) GetSpamCounts(ctx context.Context, options map[string]interface{}) (SpamCounts, error) {
	res := SpamCounts{}
	err := client.get(ctx, buildURL("stats/outbound/spam", options), &res)
	return res, err
}

// TrackedDay - tracked emails sent on a specific day
type TrackedDay struct {
	// Date - self-explanatory
	Date string
	// Tracked - number of emails tracked sent
	Tracked int64
}

// TrackedCounts - tracked emails sent for a period
type TrackedCounts struct {
	// Days - List of objects that each represent tracked email counts by date.
	Days []TrackedDay
	// Tracked - Indicates total number of tracked emails sent.
	Tracked int64
}

// GetTrackedCounts - Gets a total count of emails you’ve sent with open tracking enabled.
// Available options: http://developer.postmarkapp.com/developer-api-stats.html#email-tracked-count
func (client *Client) GetTrackedCounts(ctx context.Context, options map[string]interface{}) (TrackedCounts, error) {
	res := TrackedCounts{}
	err := client.get(ctx, buildURL("stats/outbound/tracked", options), &res)
	return res, err
}

// OpenedDay - opened outbound emails sent on a specific day
type OpenedDay struct {
	// Date - self-explanatory
	Date string
	// Opens - Indicates total number of opened emails. This total includes recipients who opened your email multiple times.
	Opens int64
	// Unique - Indicates total number of uniquely opened emails.
	Unique int64
}

// OpenCounts - opened outbound emails for a period
type OpenCounts struct {
	// Days - List of objects that each represent opens by date.
	Days []OpenedDay
	// Opens - Indicates total number of opened emails. This total includes recipients who opened your email multiple times.
	Opens int64
	// Unique int64 - Indicates total number of uniquely opened emails.
	Unique int64
}

// GetOpenCounts - Gets total counts of recipients who opened your emails. This is only recorded when open tracking is enabled for that email.
// Available options: http://developer.postmarkapp.com/developer-api-stats.html#email-opens-count
func (client *Client) GetOpenCounts(ctx context.Context, options map[string]interface{}) (OpenCounts, error) {
	res := OpenCounts{}
	err := client.get(ctx, buildURL("stats/outbound/opens", options), &res)
	return res, err
}

// PlatformCounts contains day-to-day usages, along with totals of email usages by platform
type PlatformCounts struct {
	// Days - List of objects that each represent email platform usages by date
	Days []PlatformDay
	// Desktop - The total number of email platform usages by Desktop
	Desktop int64

	// Mobile - The total number of email platform usages by Mobile
	Mobile int64

	// Unknown - The total number of email platform usages by others
	Unknown int64

	// WebMail - The total number of email platform usages by WebMail
	WebMail int64
}

// PlatformDay contains the totals of email usages by platform for a specific date
type PlatformDay struct {
	// Date - the date in question
	Date string

	// Desktop - The total number of email platform usages by Desktop for this date
	Desktop int64

	// Mobile - The total number of email platform usages by Mobile for this date
	Mobile int64

	// Unknown - The total number of email platform usages by others for this date
	Unknown int64

	// WebMail - The total number of email platform usages by WebMail for this date
	WebMail int64
}

// GetPlatformCounts gets the email platform usage
func (client *Client) GetPlatformCounts(ctx context.Context, options map[string]interface{}) (PlatformCounts, error) {
	res := PlatformCounts{}
	err := client.get(ctx, buildURL("stats/outbound/platform", options), &res)
	return res, err
}

// ClickDay - click stats for a specific day
type ClickDay struct {
	// Date - self-explanatory
	Date string
	// Clicks - number of total clicks
	Clicks int64
	// Unique - number of unique clicks
	Unique int64
}

// ClickCounts - click stats for a period
type ClickCounts struct {
	// Days - List of objects that each represent click counts by date.
	Days []ClickDay
	// Clicks - Indicates total number of clicks. This total includes recipients who clicked your links multiple times.
	Clicks int64
	// Unique - Indicates total number of uniquely clicked links.
	Unique int64
}

// GetClickCounts - Gets total counts of recipients who clicked links in your emails. This is only recorded when link tracking is enabled for that email.
// Available options: http://developer.postmarkapp.com/developer-api-stats.html#click-counts
func (client *Client) GetClickCounts(ctx context.Context, options map[string]interface{}) (ClickCounts, error) {
	res := ClickCounts{}
	err := client.get(ctx, buildURL("stats/outbound/clicks", options), &res)
	return res, err
}

// BrowserFamilyDay - browser family usage stats for a specific day
type BrowserFamilyDay struct {
	// Date - self-explanatory
	Date string
	// Chrome - number of clicks from Chrome browser
	Chrome int64
	// Safari - number of clicks from Safari browser
	Safari int64
	// Firefox - number of clicks from Firefox browser
	Firefox int64
	// InternetExplorer - number of clicks from Internet Explorer browser
	InternetExplorer int64
	// Opera - number of clicks from Opera browser
	Opera int64
	// Unknown - number of clicks from unknown browsers
	Unknown int64
}

// BrowserFamilyCounts - browser family usage stats for a period
type BrowserFamilyCounts struct {
	// Days - List of objects that each represent browser family usage by date.
	Days []BrowserFamilyDay
	// Chrome - total number of clicks from Chrome browser
	Chrome int64
	// Safari - total number of clicks from Safari browser
	Safari int64
	// Firefox - total number of clicks from Firefox browser
	Firefox int64
	// InternetExplorer - total number of clicks from Internet Explorer browser
	InternetExplorer int64
	// Opera - total number of clicks from Opera browser
	Opera int64
	// Unknown - total number of clicks from unknown browsers
	Unknown int64
}

// GetBrowserFamilyCounts - Gets total counts of clicks by browser family. This is only recorded when link tracking is enabled for that email.
// Available options: http://developer.postmarkapp.com/developer-api-stats.html#browser-usage
func (client *Client) GetBrowserFamilyCounts(ctx context.Context, options map[string]interface{}) (BrowserFamilyCounts, error) {
	res := BrowserFamilyCounts{}
	err := client.get(ctx, buildURL("stats/outbound/clicks/browserfamilies", options), &res)
	return res, err
}

// ClickLocationDay - click location stats for a specific day
type ClickLocationDay struct {
	// Date - self-explanatory
	Date string
	// HTML - number of clicks from HTML part of the email
	HTML int64
	// Text - number of clicks from text part of the email
	Text int64
}

// ClickLocationCounts - click location stats for a period
type ClickLocationCounts struct {
	// Days - List of objects that each represent click location counts by date.
	Days []ClickLocationDay
	// HTML - total number of clicks from HTML part of the email
	HTML int64
	// Text - total number of clicks from text part of the email
	Text int64
}

// GetClickLocationCounts - Gets total counts of clicks by email format. This is only recorded when link tracking is enabled for that email.
// Available options: http://developer.postmarkapp.com/developer-api-stats.html#click-location
func (client *Client) GetClickLocationCounts(ctx context.Context, options map[string]interface{}) (ClickLocationCounts, error) {
	res := ClickLocationCounts{}
	err := client.get(ctx, buildURL("stats/outbound/clicks/location", options), &res)
	return res, err
}

// ClickPlatformDay - click platform usage stats for a specific day
type ClickPlatformDay struct {
	// Date - self-explanatory
	Date string
	// Desktop - number of clicks from Desktop platforms
	Desktop int64
	// Mobile - number of clicks from Mobile platforms
	Mobile int64
	// Unknown - number of clicks from unknown platforms
	Unknown int64
	// WebMail - number of clicks from WebMail platforms
	WebMail int64
}

// ClickPlatformCounts - click platform usage stats for a period
type ClickPlatformCounts struct {
	// Days - List of objects that each represent click platform usage by date.
	Days []ClickPlatformDay
	// Desktop - total number of clicks from Desktop platforms
	Desktop int64
	// Mobile - total number of clicks from Mobile platforms
	Mobile int64
	// Unknown - total number of clicks from unknown platforms
	Unknown int64
	// WebMail - total number of clicks from WebMail platforms
	WebMail int64
}

// GetClickPlatformCounts - Gets total counts of clicks by platform. This is only recorded when link tracking is enabled for that email.
// Available options: http://developer.postmarkapp.com/developer-api-stats.html#browser-platform-usage
func (client *Client) GetClickPlatformCounts(ctx context.Context, options map[string]interface{}) (ClickPlatformCounts, error) {
	res := ClickPlatformCounts{}
	err := client.get(ctx, buildURL("stats/outbound/clicks/platforms", options), &res)
	return res, err
}

// EmailClientDay - email client usage stats for a specific day
type EmailClientDay struct {
	// Date - self-explanatory
	Date string
	// Outlook - number of opens from Outlook email client
	Outlook int64
	// Gmail - number of opens from Gmail email client
	Gmail int64
	// AppleMail - number of opens from Apple Mail email client
	AppleMail int64
	// Thunderbird - number of opens from Thunderbird email client
	Thunderbird int64
	// Yahoo - number of opens from Yahoo email client
	Yahoo int64
	// Unknown - number of opens from unknown email clients
	Unknown int64
}

// EmailClientCounts - email client usage stats for a period
type EmailClientCounts struct {
	// Days - List of objects that each represent email client usage by date.
	Days []EmailClientDay
	// Outlook - total number of opens from Outlook email client
	Outlook int64
	// Gmail - total number of opens from Gmail email client
	Gmail int64
	// AppleMail - total number of opens from Apple Mail email client
	AppleMail int64
	// Thunderbird - total number of opens from Thunderbird email client
	Thunderbird int64
	// Yahoo - total number of opens from Yahoo email client
	Yahoo int64
	// Unknown - total number of opens from unknown email clients
	Unknown int64
}

// GetEmailClientCounts - Gets total counts of opens by email client. This is only recorded when open tracking is enabled for that email.
// Available options: http://developer.postmarkapp.com/developer-api-stats.html#email-client-usage
func (client *Client) GetEmailClientCounts(ctx context.Context, options map[string]interface{}) (EmailClientCounts, error) {
	res := EmailClientCounts{}
	err := client.get(ctx, buildURL("stats/outbound/opens/emailclients", options), &res)
	return res, err
}
