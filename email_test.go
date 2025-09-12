package postmark

import (
	"context"
	"net/http"
	"testing"

	"goji.io/pat"
)

func getTestEmail() Email {
	return Email{
		From:     "sender@example.com",
		To:       "receiver@example.com",
		Cc:       "copied@example.com",
		Bcc:      "blank-copied@example.com",
		Subject:  "Test",
		Tag:      "Invitation",
		HTMLBody: "<b>Hello</b>",
		TextBody: "Hello",
		ReplyTo:  "reply@example.com",
		Headers: []Header{
			{
				Name:  "CUSTOM-HEADER",
				Value: "value",
			},
		},
		TrackOpens: true,
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

func TestSendEmail(t *testing.T) {
	responseJSON := `{
		"To": "receiver@example.com",
		"SubmittedAt": "2014-02-17T07:25:01.4178645-05:00",
		"MessageID": "0a129aee-e1cd-480d-b08d-4f48548ff48d",
		"ErrorCode": 0,
		"Message": "OK"
	}`

	tMux.HandleFunc(pat.Post("/email"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	// Success
	res, err := client.SendEmail(context.Background(), getTestEmail())
	if err != nil {
		t.Fatalf("SendEmail: %s", err.Error())
	}

	if res.MessageID != "0a129aee-e1cd-480d-b08d-4f48548ff48d" {
		t.Fatalf("SendEmail: wrong id!")
	}

	// Failure
	responseJSON = `{
		"To": "receiver@example.com",
		"SubmittedAt": "2014-02-17T07:25:01.4178645-05:00",
		"MessageID": "0a129aee-e1cd-480d-b08d-4f48548ff48d",
		"ErrorCode": 401,
		"Message": "Sender signature not confirmed"
	}`

	_, err = client.SendEmail(context.Background(), getTestEmail())

	if err == nil {
		t.Fatalf("SendEmail should have failed")
	}
}

func TestSendEmailBatch(t *testing.T) {
	responseJSON := `[
	  {
		"ErrorCode": 0,
		"Message": "OK",
		"MessageID": "b7bc2f4a-e38e-4336-af7d-e6c392c2f817",
		"SubmittedAt": "2010-11-26T12:01:05.1794748-05:00",
		"To": "receiver1@example.com"
	  },
	  {
		"ErrorCode": 0,
		"Message": "OK",
		"MessageID": "e2ecbbfc-fe12-463d-b933-9fe22915106d",
		"SubmittedAt": "2010-11-26T12:01:05.1794748-05:00",
		"To": "receiver2@example.com"
	  }
	]`

	tMux.HandleFunc(pat.Post("/email/batch"), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(responseJSON))
	})

	testEmail := getTestEmail()
	res, err := client.SendEmailBatch(context.Background(), []Email{testEmail, testEmail})
	if err != nil {
		t.Fatalf("SendEmailBatch: %s", err.Error())
	}

	if len(res) != 2 {
		t.Fatalf("SendEmailBatch: wrong response array size!")
	}
}
