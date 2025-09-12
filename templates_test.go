package postmark

import (
	"context"
	"net/http"

	"goji.io/pat"
)

func (s *PostmarkTestSuite) TestGetTemplate() {
	responseJSON := `{
		"Name": "Onboarding Email",
		"TemplateId": 1234,
		"Subject": "Hi there, {{Name}}",
		"HtmlBody": "Hello dear Postmark user. {{Name}}",
		"TextBody": "{{Name}} is a {{Occupation}}",
		"AssociatedServerId": 1,
		"Active": false
	}`

	s.mux.HandleFunc(pat.Get("/templates/:templateID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.GetTemplate(context.Background(), "1234")
	s.Require().NoError(err)

	s.Equal("Onboarding Email", res.Name, "Template: wrong name")
}

func (s *PostmarkTestSuite) TestGetTemplates() {
	responseJSON := `{
		"TotalCount": 2,
		"Templates": [
		  {
			"Active": true,
			"TemplateId": 1234,
			"Name": "Account Activation Email"
		  },
		  {
			"Active": true,
			"TemplateId": 5678,
			"Name": "Password Recovery Email"
		  }
		]
	}`

	s.mux.HandleFunc(pat.Get("/templates"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, count, err := s.client.GetTemplates(context.Background(), 100, 10)
	s.Require().NoError(err)

	s.NotEmpty(res, "GetTemplates: result should not be empty")
	s.Equal(int64(2), count, "GetTemplates: wrong count")
}

func (s *PostmarkTestSuite) TestCreateTemplate() {
	responseJSON := `{
		"TemplateId": 1234,
		"Name": "Onboarding Email",
		"Active": true
	}`

	s.mux.HandleFunc(pat.Post("/templates"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.CreateTemplate(context.Background(), Template{
		Name:     "Onboarding Email",
		Subject:  "Hello from {{company.name}}!",
		TextBody: "Hello, {{name}}!",
		HTMLBody: "<html><body>Hello, {{name}}!</body></html>",
	})
	s.Require().NoError(err)

	s.Equal("Onboarding Email", res.Name, "CreateTemplate: wrong name")
}

func (s *PostmarkTestSuite) TestEditTemplate() {
	responseJSON := `{
		"TemplateId": 1234,
		  "Name": "Onboarding Emailzzzzz",
		  "Active": true
	}`

	s.mux.HandleFunc(pat.Put("/templates/:templateID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.EditTemplate(context.Background(), "1234", Template{
		Name:     "Onboarding Emailzzzzz",
		Subject:  "Hello from {{company.name}}!",
		TextBody: "Hello, {{name}}!",
		HTMLBody: "<html><body>Hello, {{name}}!</body></html>",
	})
	s.Require().NoError(err)

	s.Equal("Onboarding Emailzzzzz", res.Name, "EditTemplate: wrong name")
}

func (s *PostmarkTestSuite) TestDeleteTemplate() {
	responseJSON := `{
	  "ErrorCode": 0,
	  "Message": "Template 1234 removed."
	}`

	s.mux.HandleFunc(pat.Delete("/templates/:templateID"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	// Success
	err := s.client.DeleteTemplate(context.Background(), "1234")
	s.Require().NoError(err)

	// Failure
	responseJSON = `{
	  "ErrorCode": 402,
	  "Message": "Invalid JSON"
	}`

	err = s.client.DeleteTemplate(context.Background(), "1234")
	s.Require().Error(err, "DeleteTemplate should have failed")
}

func (s *PostmarkTestSuite) TestValidateTemplate() {
	responseJSON := `{
		"AllContentIsValid": true,
		"HtmlBody": {
			"ContentIsValid": true,
			"ValidationErrors": [],
			"RenderedContent": "address_Value name_Value "
		},
		"TextBody": {
			"ContentIsValid": true,
			"ValidationErrors": [{
				"Message" : "The syntax for this template is invalid.",
				"Line" : 1,
				"CharacterPosition" : 1
			}],
			"RenderedContent": "phone_Value name_Value "
		},
		"Subject": {
			"ContentIsValid": true,
			"ValidationErrors": [],
			"RenderedContent": "name_Value subjectHeadline_Value"
		},
		"SuggestedTemplateModel": {
			"userName": "bobby joe",
			"company": {
			"address": "address_Value",
			"phone": "phone_Value",
			"name": "name_Value"
			},
			"person": [
			{
				"name": "name_Value"
			}
			],
			"subjectHeadline": "subjectHeadline_Value"
		}
	}`

	s.mux.HandleFunc(pat.Post("/templates/validate"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.ValidateTemplate(context.Background(), ValidateTemplateBody{
		Subject:  "{{#company}}{{name}}{{/company}} {{subjectHeadline}}",
		TextBody: "{{#company}}{{address}}{{/company}}{{#each person}} {{name}} {{/each}}",
		HTMLBody: "{{#company}}{{phone}}{{/company}}{{#each person}} {{name}} {{/each}}",
		TestRenderModel: map[string]interface{}{
			"userName": "bobby joe",
		},
		InlineCSSForHTMLTestRender: false,
	})
	s.Require().NoError(err)

	s.True(res.AllContentIsValid, "ValidateTemplate: AllContentIsValid should be true")
}

func getTestTemplatedEmail() TemplatedEmail {
	return TemplatedEmail{
		TemplateID: 1234,
		TemplateModel: map[string]interface{}{
			"user_name": "John Smith",
			"company": map[string]interface{}{
				"name": "ACME",
			},
		},
		InlineCSS: true,
		From:      "sender@example.com",
		To:        "receiver@example.com",
		Cc:        "copied@example.com",
		Bcc:       "blank-copied@example.com",
		Tag:       "Invitation",
		ReplyTo:   "reply@example.com",
		Headers: []Header{
			{
				Name:  "CUSTOM-HEADER",
				Value: "value",
			},
		},
		TrackOpens: true,
		TrackLinks: "HtmlAndText",
		Attachments: []Attachment{
			{
				Name:        "readme.txt",
				Content:     "dGVzdCBjb250ZW50",
				ContentType: "text/plain",
			},
			{
				Name:        "report.pdf",
				Content:     "dGVzdCBjb250ZW50",
				ContentType: "application/octet-stream",
			},
		},
	}
}

func (s *PostmarkTestSuite) TestSendTemplatedEmail() {
	responseJSON := `{
		"To": "receiver@example.com",
		"SubmittedAt": "2014-02-17T07:25:01.4178645-05:00",
		"MessageID": "0a129aee-e1cd-480d-b08d-4f48548ff48d",
		"ErrorCode": 0,
		"Message": "OK"
	}`

	s.mux.HandleFunc(pat.Post("/email/withTemplate"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	res, err := s.client.SendTemplatedEmail(context.Background(), getTestTemplatedEmail())
	s.Require().NoError(err)

	s.Equal("0a129aee-e1cd-480d-b08d-4f48548ff48d", res.MessageID, "SendTemplatedEmail: incorrect message ID")
}

func (s *PostmarkTestSuite) TestSendTemplatedBatch() {
	responseJSON := `[
	  {
		"To": "receiver@example.com",
		"SubmittedAt": "2014-02-17T07:25:01.4178645-05:00",
		"MessageID": "0a129aee-e1cd-480d-b08d-4f48548ff48d",
		"ErrorCode": 0,
		"Message": "OK"
	},{
		"To": "receiver@example.com",
		"SubmittedAt": "2014-02-17T07:25:01.4178645-05:00",
		"MessageID": "0a129aee-e1cd-480d-b08d-4f48548ff48d",
		"ErrorCode": 0,
		"Message": "OK"
	}
	]`

	s.mux.HandleFunc(pat.Post("/email/batchWithTemplates"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	testTemplatedEmail := getTestTemplatedEmail()
	res, err := s.client.SendTemplatedEmailBatch(context.Background(), []TemplatedEmail{testTemplatedEmail, testTemplatedEmail})
	s.Require().NoError(err)

	s.Len(res, 2, "SendTemplatedBatch: wrong response array size")
}
