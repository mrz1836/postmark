package postmark

import (
	"context"
	"net/http"
)

func (s *PostmarkTestSuite) TestGetSuppressions() {
	responseJSON := `{
		"Suppressions":[
		  {
			"EmailAddress":"address@wildbit.com",
			"SuppressionReason":"ManualSuppression",
			"Origin": "Recipient",
			"CreatedAt":"2019-12-10T08:58:33-05:00"
		  },
		  {
			"EmailAddress":"bounce.address@wildbit.com",
			"SuppressionReason":"HardBounce",
			"Origin": "Recipient",
			"CreatedAt":"2019-12-11T08:58:33-05:00"
		  },
		  {
			"EmailAddress":"spam.complaint.address@wildbit.com",
			"SuppressionReason":"SpamComplaint",
			"Origin": "Recipient",
			"CreatedAt":"2019-12-12T08:58:33-05:00"
		  }
		]
	  }`

	s.mux.Get("/message-streams/:StreamID/suppressions/dump", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetSuppressions(context.Background(), "outbound", nil)
	s.Require().NoError(err)

	s.Len(res, 3, "GetSuppressions: wrong number of suppressions")
	s.Equal("address@wildbit.com", res[0].EmailAddress, "GetSuppressions: wrong suppression email address")

	responseJSON = `{
		"Suppressions":[
		  {
			"EmailAddress":"address@wildbit.com",
			"SuppressionReason":"ManualSuppression",
			"Origin": "Recipient",
			"CreatedAt":"2019-12-10T08:58:33-05:00"
		  }
		]
	  }`

	s.mux.Get("/message-streams/:StreamID/suppressions/dump", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err = s.client.GetSuppressions(context.Background(), "outbound", map[string]interface{}{
		"emailaddress":      "address@wildbit.com",
		"fromdate":          "2019-12-10",
		"todate":            "2019-12-11",
		"suppressionreason": HardBounceReason,
		"origin":            RecipientOrigin,
	})
	s.Require().NoError(err)

	s.Len(res, 1, "GetSuppressions: wrong number of suppressions")
	s.Equal("address@wildbit.com", res[0].EmailAddress, "GetSuppressions: wrong suppression email address")
}

func (s *PostmarkTestSuite) TestCreateSuppressions() {
	responseJSON := `{
		"Suppressions":[
		  {
			"EmailAddress":"good.address@wildbit.com",
			"Status":"Suppressed",
			"Message": null
		  },
		  {
			"EmailAddress":"spammy.address@wildbit.com",
			"Status":"Failed",
			"Message": "You do not have the required authority to change this suppression."
		  },
		  {
			"EmailAddress":"invalid-email-address",
			"Status":"Failed",
			"Message": "An invalid email address was provided."
		  }
		]
	  }`

	s.mux.Post("/message-streams/:StreamID/suppressions", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.CreateSuppressions(context.Background(), "outbound", []Suppression{})
	s.Require().NoError(err)

	s.Len(res, 3, "CreateSuppressions: wrong number of suppressions")
	s.Equal("good.address@wildbit.com", res[0].EmailAddress, "CreateSuppressions: wrong suppression email address")
}

func (s *PostmarkTestSuite) TestDeleteSuppressions() {
	responseJSON := `{
		"Suppressions":[
		  {
			"EmailAddress":"good.address@wildbit.com"
		  },
		  {
			"EmailAddress":"not.suppressed@wildbit.com"
		  },
		  {
			"EmailAddress":"spammy.address@wildbit.com"
		  }
		]
	  }`

	s.mux.Post("/message-streams/:StreamID/suppressions/delete", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.DeleteSuppressions(context.Background(), "outbound", []Suppression{})
	s.Require().NoError(err)

	s.Len(res, 3, "DeleteSuppressions: wrong number of suppressions")
	s.Equal("good.address@wildbit.com", res[0].EmailAddress, "DeleteSuppressions: wrong suppression email address")
}
