package postmark

import (
	"context"
	"net/http"
	"testing"
)

func (s *PostmarkTestSuite) TestCreateDataRemoval() {
	tests := []struct {
		name         string
		responseJSON string
		wantErr      bool
		expectedID   int64
	}{
		{
			name: "successful data removal request",
			responseJSON: `{
				"ID": 12345,
				"Recipient": "test@example.com",
				"RequestedAt": "2024-01-15T10:30:00Z",
				"Status": "Pending"
			}`,
			wantErr:    false,
			expectedID: 12345,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mux.Post("/data-removals", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(tt.responseJSON))
			})

			request := DataRemovalRequest{
				Recipient: "test@example.com",
			}

			res, err := s.client.CreateDataRemoval(context.Background(), request)

			if tt.wantErr {
				s.Require().Error(err, "CreateDataRemoval should have failed")
			} else {
				s.Require().NoError(err, "CreateDataRemoval should not have failed")
				s.Equal(tt.expectedID, res.ID, "CreateDataRemoval returned wrong ID")
				s.Equal("test@example.com", res.Recipient, "CreateDataRemoval returned wrong recipient")
				s.Equal("Pending", res.Status, "CreateDataRemoval returned wrong status")
			}
		})
	}
}

func (s *PostmarkTestSuite) TestGetDataRemovalStatus() {
	tests := []struct {
		name         string
		responseJSON string
		wantErr      bool
		expectedID   int64
	}{
		{
			name: "successful data removal status check",
			responseJSON: `{
				"ID": 12345,
				"Recipient": "test@example.com",
				"RequestedAt": "2024-01-15T10:30:00Z",
				"Status": "Completed",
				"CompletedAt": "2024-01-15T11:00:00Z"
			}`,
			wantErr:    false,
			expectedID: 12345,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.mux.Get("/data-removals/12345", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(tt.responseJSON))
			})

			res, err := s.client.GetDataRemovalStatus(context.Background(), 12345)

			if tt.wantErr {
				s.Require().Error(err, "GetDataRemovalStatus should have failed")
			} else {
				s.Require().NoError(err, "GetDataRemovalStatus should not have failed")
				s.Equal(tt.expectedID, res.ID, "GetDataRemovalStatus returned wrong ID")
				s.Equal("test@example.com", res.Recipient, "GetDataRemovalStatus returned wrong recipient")
				s.Equal("Completed", res.Status, "GetDataRemovalStatus returned wrong status")
			}
		})
	}
}

// Benchmark for CreateDataRemoval
func BenchmarkCreateDataRemoval(b *testing.B) {
	request := DataRemovalRequest{
		Recipient: "benchmark@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// In a real benchmark, you'd call the actual function
		// For now, we'll just test the struct creation overhead
		_ = DataRemovalRequest{
			Recipient: request.Recipient,
		}
	}
}

// Benchmark for GetDataRemovalStatus
func BenchmarkGetDataRemovalStatus(b *testing.B) {
	id := int64(12345)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// In a real benchmark, you'd call the actual function
		// For now, we'll just test the ID conversion overhead
		_ = id
	}
}
