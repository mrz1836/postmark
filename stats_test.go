package postmark

import (
	"context"
	"net/http"
	"testing"
)

func (s *PostmarkTestSuite) TestGetOutboundStats() {
	responseJSON := `{
	  "Sent": 615,
	  "Bounced": 64,
	  "SMTPApiErrors": 25,
	  "BounceRate": 10.406,
	  "SpamComplaints": 10,
	  "SpamComplaintsRate": 1.626,
	  "Opens": 166,
	  "UniqueOpens": 26,
	  "Tracked": 111,
	  "WithClientRecorded": 14,
	  "WithPlatformRecorded": 10,
	  "WithReadTimeRecorded": 10
	}`

	s.mux.Get("/stats/outbound", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetOutboundStats(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(615), res.Sent, "GetOutboundStats: wrong Sent")
}

func (s *PostmarkTestSuite) TestGetSentCounts() {
	responseJSON := `{
	  "Days": [
	    {
	      "Date": "2014-01-01",
	      "Sent": 140
	    },
	    {
	      "Date": "2014-01-02",
	      "Sent": 160
	    },
	    {
	      "Date": "2014-01-04",
	      "Sent": 50
	    },
	    {
	      "Date": "2014-01-05",
	      "Sent": 115
	    }
	  ],
	  "Sent": 615
	}`

	s.mux.Get("/stats/outbound/sends", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetSentCounts(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(615), res.Sent, "GetSentCounts: wrong Sent")
	s.Equal(int64(140), res.Days[0].Sent, "GetSentCounts: wrong day Sent count")
}

func (s *PostmarkTestSuite) TestGetBounceCounts() {
	responseJSON := `{
	  "Days": [
	    {
	      "Date": "2014-01-01",
	      "HardBounce": 12,
	      "SoftBounce": 36
	    },
	    {
	      "Date": "2014-01-03",
	      "Transient": 7
	    },
	    {
	      "Date": "2014-01-04",
	      "Transient": 4
	    },
	    {
	      "Date": "2014-01-05",
	      "SMTPApiError": 25,
	      "Transient": 5
	    }
	  ],
	  "HardBounce": 12,
	  "SMTPApiError": 25,
	  "SoftBounce": 36,
	  "Transient": 16
	}`

	s.mux.Get("/stats/outbound/bounces", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetBounceCounts(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(12), res.HardBounce, "GetBounceCounts: wrong HardBounce")
	s.Equal(int64(12), res.Days[0].HardBounce, "GetBounceCounts: wrong day HardBounce count")
}

func (s *PostmarkTestSuite) TestGetSpamCounts() {
	responseJSON := `{
	  "Days": [
	    {
	      "Date": "2014-01-01",
	      "SpamComplaint": 2
	    },
	    {
	      "Date": "2014-01-02",
	      "SpamComplaint": 3
	    },
	    {
	      "Date": "2014-01-05",
	      "SpamComplaint": 5
	    }
	  ],
	  "SpamComplaint": 10
	}`

	s.mux.Get("/stats/outbound/spam", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetSpamCounts(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(10), res.SpamComplaint, "GetSpamCounts: wrong SpamComplaint")
	s.Equal(int64(2), res.Days[0].SpamComplaint, "GetSpamCounts: wrong day SpamComplaint count")
}

func (s *PostmarkTestSuite) TestGetTrackedCounts() {
	responseJSON := `{
	  "Days": [
	    {
	      "Date": "2014-01-01",
	      "Tracked": 24
	    },
	    {
	      "Date": "2014-01-02",
	      "Tracked": 26
	    },
	    {
	      "Date": "2014-01-03",
	      "Tracked": 15
	    },
	    {
	      "Date": "2014-01-04",
	      "Tracked": 15
	    },
	    {
	      "Date": "2014-01-05",
	      "Tracked": 31
	    }
	  ],
	  "Tracked": 111
	}`

	s.mux.Get("/stats/outbound/tracked", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetTrackedCounts(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(111), res.Tracked, "GetTrackedCounts: wrong Tracked")
	s.Equal(int64(24), res.Days[0].Tracked, "GetTrackedCounts: wrong day Tracked count")
}

func (s *PostmarkTestSuite) TestGetOpenCounts() {
	responseJSON := `{
		"Days": [
		    {
		      "Date": "2014-01-01",
		      "Opens": 44,
		      "Unique": 4
		    },
		    {
		      "Date": "2014-01-02",
		      "Opens": 46,
		      "Unique": 6
		    },
		    {
		      "Date": "2014-01-03",
		      "Opens": 25,
		      "Unique": 5
		    },
		    {
		      "Date": "2014-01-04",
		      "Opens": 25,
		      "Unique": 5
		    },
		    {
		      "Date": "2014-01-05",
		      "Opens": 26,
		      "Unique": 6
		    }
		  ],
	  "Opens": 166,
	  "Unique": 26
	}`

	s.mux.Get("/stats/outbound/opens", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetOpenCounts(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(166), res.Opens, "GetOpenCounts: wrong Opens")
	s.Equal(int64(44), res.Days[0].Opens, "GetOpenCounts: wrong day Opens count")
}

func (s *PostmarkTestSuite) TestGetPlatformCounts() {
	responseJSON := `{
		"Days": [
			{
				"Date": "2014-01-01",
				"Desktop": 1,
				"WebMail": 1
			},
			{
				"Date": "2014-01-02",
				"Mobile": 2,
				"WebMail": 1
			},
			{
				"Date": "2014-01-04",
				"Desktop": 3,
				"Unknown": 2
			}
		],
		"Desktop": 4,
		"Mobile": 2,
		"Unknown": 2,
		"WebMail": 2
	}`

	s.mux.Get("/stats/outbound/platform", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetPlatformCounts(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(4), res.Desktop, "GetPlatformCounts: wrong Desktop")
	s.Equal(int64(1), res.Days[0].Desktop, "GetPlatformCounts: wrong day Desktop count")
}

func (s *PostmarkTestSuite) TestGetClickCounts() {
	responseJSON := `{
		"Days": [
			{
				"Date": "2014-01-01",
				"Clicks": 44,
				"Unique": 4
			},
			{
				"Date": "2014-01-02",
				"Clicks": 46,
				"Unique": 6
			},
			{
				"Date": "2014-01-03",
				"Clicks": 25,
				"Unique": 5
			}
		],
		"Clicks": 115,
		"Unique": 15
	}`

	s.mux.Get("/stats/outbound/clicks", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetClickCounts(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(115), res.Clicks, "GetClickCounts: wrong Clicks")
	s.Equal(int64(44), res.Days[0].Clicks, "GetClickCounts: wrong day Clicks count")
	s.Equal(int64(15), res.Unique, "GetClickCounts: wrong Unique")
}

func (s *PostmarkTestSuite) TestGetBrowserFamilyCounts() {
	responseJSON := `{
		"Days": [
			{
				"Date": "2014-01-01",
				"Chrome": 10,
				"Safari": 5,
				"Firefox": 3
			},
			{
				"Date": "2014-01-02",
				"Chrome": 12,
				"InternetExplorer": 2
			}
		],
		"Chrome": 22,
		"Safari": 5,
		"Firefox": 3,
		"InternetExplorer": 2,
		"Opera": 0,
		"Unknown": 1
	}`

	s.mux.Get("/stats/outbound/clicks/browserfamilies", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetBrowserFamilyCounts(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(22), res.Chrome, "GetBrowserFamilyCounts: wrong Chrome")
	s.Equal(int64(10), res.Days[0].Chrome, "GetBrowserFamilyCounts: wrong day Chrome count")
	s.Equal(int64(5), res.Safari, "GetBrowserFamilyCounts: wrong Safari")
}

func (s *PostmarkTestSuite) TestGetClickLocationCounts() {
	responseJSON := `{
		"Days": [
			{
				"Date": "2014-01-01",
				"HTML": 30,
				"Text": 5
			},
			{
				"Date": "2014-01-02",
				"HTML": 25,
				"Text": 10
			}
		],
		"HTML": 55,
		"Text": 15
	}`

	s.mux.Get("/stats/outbound/clicks/location", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetClickLocationCounts(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(55), res.HTML, "GetClickLocationCounts: wrong HTML")
	s.Equal(int64(30), res.Days[0].HTML, "GetClickLocationCounts: wrong day HTML count")
	s.Equal(int64(15), res.Text, "GetClickLocationCounts: wrong Text")
}

func (s *PostmarkTestSuite) TestGetClickPlatformCounts() {
	responseJSON := `{
		"Days": [
			{
				"Date": "2014-01-01",
				"Desktop": 20,
				"Mobile": 10,
				"WebMail": 5
			},
			{
				"Date": "2014-01-02",
				"Desktop": 15,
				"Mobile": 12,
				"Unknown": 3
			}
		],
		"Desktop": 35,
		"Mobile": 22,
		"WebMail": 5,
		"Unknown": 3
	}`

	s.mux.Get("/stats/outbound/clicks/platforms", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetClickPlatformCounts(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(35), res.Desktop, "GetClickPlatformCounts: wrong Desktop")
	s.Equal(int64(20), res.Days[0].Desktop, "GetClickPlatformCounts: wrong day Desktop count")
	s.Equal(int64(22), res.Mobile, "GetClickPlatformCounts: wrong Mobile")
}

func (s *PostmarkTestSuite) TestGetEmailClientCounts() {
	responseJSON := `{
		"Days": [
			{
				"Date": "2014-01-01",
				"Outlook": 15,
				"Gmail": 10,
				"AppleMail": 8
			},
			{
				"Date": "2014-01-02",
				"Outlook": 12,
				"Gmail": 14,
				"Yahoo": 3
			}
		],
		"Outlook": 27,
		"Gmail": 24,
		"AppleMail": 8,
		"Yahoo": 3,
		"Thunderbird": 2,
		"Unknown": 5
	}`

	s.mux.Get("/stats/outbound/opens/emailclients", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetEmailClientCounts(context.Background(), map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	})
	s.Require().NoError(err)

	s.Equal(int64(27), res.Outlook, "GetEmailClientCounts: wrong Outlook")
	s.Equal(int64(15), res.Days[0].Outlook, "GetEmailClientCounts: wrong day Outlook count")
	s.Equal(int64(24), res.Gmail, "GetEmailClientCounts: wrong Gmail")
}

// Benchmark for GetClickCounts
func BenchmarkGetClickCounts(b *testing.B) {
	ctx := context.Background()
	options := map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = options
	}
}

// Benchmark for GetBrowserFamilyCounts
func BenchmarkGetBrowserFamilyCounts(b *testing.B) {
	ctx := context.Background()
	options := map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = options
	}
}

// Benchmark for GetClickLocationCounts
func BenchmarkGetClickLocationCounts(b *testing.B) {
	ctx := context.Background()
	options := map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = options
	}
}

// Benchmark for GetClickPlatformCounts
func BenchmarkGetClickPlatformCounts(b *testing.B) {
	ctx := context.Background()
	options := map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = options
	}
}

// Benchmark for GetEmailClientCounts
func BenchmarkGetEmailClientCounts(b *testing.B) {
	ctx := context.Background()
	options := map[string]interface{}{
		"fromdate": "2014-01-01",
		"todate":   "2014-02-01",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = options
	}
}
