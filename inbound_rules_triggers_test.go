package postmark

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func (s *PostmarkTestSuite) TestGetInboundRuleTriggers() {
	responseJSON := `{
		"TotalCount": 2,
		"InboundRules": [
			{
				"ID": 123456,
				"Rule": "spam@example.com"
			},
			{
				"ID": 123457,
				"Rule": "*.spammer.com"
			}
		]
	}`

	s.mux.Get("/triggers/inboundrules", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	triggers, totalCount, err := s.client.GetInboundRuleTriggers(context.Background(), 50, 0)
	s.Require().NoError(err, "GetInboundRuleTriggers should not fail")
	s.Equal(int64(2), totalCount, "TotalCount should match")
	s.Len(triggers, 2, "Should return 2 triggers")
	s.Equal(int64(123456), triggers[0].ID)
	s.Equal("spam@example.com", triggers[0].Rule)
	s.Equal(int64(123457), triggers[1].ID)
	s.Equal("*.spammer.com", triggers[1].Rule)
}

func (s *PostmarkTestSuite) TestCreateInboundRuleTrigger() {
	responseJSON := `{
		"ID": 123456,
		"Rule": "spam@example.com"
	}`

	s.mux.Post("/triggers/inboundrules", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	trigger, err := s.client.CreateInboundRuleTrigger(context.Background(), "spam@example.com")
	s.Require().NoError(err, "CreateInboundRuleTrigger should not fail")
	s.Equal(int64(123456), trigger.ID)
	s.Equal("spam@example.com", trigger.Rule)
}

func (s *PostmarkTestSuite) TestDeleteInboundRuleTrigger() {
	responseJSON := `{
		"ErrorCode": 0,
		"Message": "Trigger removed"
	}`

	s.mux.Delete("/triggers/inboundrules/123456", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	err := s.client.DeleteInboundRuleTrigger(context.Background(), 123456)
	s.Require().NoError(err, "DeleteInboundRuleTrigger should not fail")
}

func (s *PostmarkTestSuite) TestDeleteInboundRuleTriggerNotFound() {
	responseJSON := `{
		"ErrorCode": 1,
		"Message": "Trigger not found"
	}`

	s.mux.Delete("/triggers/inboundrules/999999", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	err := s.client.DeleteInboundRuleTrigger(context.Background(), 999999)
	s.Require().Error(err, "DeleteInboundRuleTrigger should fail for non-existent trigger")
	var apiErr APIError
	s.True(errors.As(err, &apiErr), "Error should be APIError")
	s.Equal(int64(1), apiErr.ErrorCode)
	s.Equal("Trigger not found", apiErr.Message)
}

// Benchmarks

func BenchmarkGetInboundRuleTriggers(b *testing.B) {
	ctx := context.Background()
	count := int64(50)
	offset := int64(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = count
		_ = offset
	}
}

func BenchmarkCreateInboundRuleTrigger(b *testing.B) {
	ctx := context.Background()
	rule := "benchmark@example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = rule
	}
}

func BenchmarkDeleteInboundRuleTrigger(b *testing.B) {
	ctx := context.Background()
	triggerID := int64(123456)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx
		_ = triggerID
	}
}
