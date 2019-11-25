package email

import (
	"github.com/ansel1/merry"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var errorNotAccepted = merry.New("email request not accepted by SendGrid")
var sendGridAPIKey string

// Init sets the API key for SendGrid
func Init(SendGridAPIKey string) {
	sendGridAPIKey = SendGridAPIKey
}

// Send a single email to a single user without attachments.
func Send(fromName string, fromEmail string,
	toName string, toEmail string,
	subject string, plainMsg string, htmlMsg string) error {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, plainMsg, htmlMsg)
	client := sendgrid.NewSendClient(sendGridAPIKey)
	response, err := client.Send(message)
	if err != nil {
		return merry.Wrap(err)
	} else if response.StatusCode != 202 {
		return errorNotAccepted.Here()
	}
	return nil
}
