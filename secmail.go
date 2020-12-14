package secmail

import (
	"context"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/mailgun/mailgun-go/v4"
)

// Recipient of the email. This type includes their email address and a public
// key for encryption.
type Recipient struct {
	Email string
	Key   string
}

// NewRecipient creates a new Recipient.
func NewRecipient(email, key string) *Recipient {
	return &Recipient{
		Email: email,
		Key:   key,
	}
}

// Mailer sends the email.
type Mailer struct {
	SenderEmail   string
	MailgunDomain string
	MailgunAPIKey string
}

// NewMailer returns a new Mailer.
func NewMailer(sender, domain, key string) *Mailer {
	return &Mailer{
		SenderEmail:   sender,
		MailgunDomain: domain,
		MailgunAPIKey: key,
	}
}

// Send the recipient an encrypted email.
func (m Mailer) Send(rcpt Recipient, subject string, message string) (string, error) {
	armor, err := helper.EncryptMessageArmored(rcpt.Key, message)
	if err != nil {
		return "", err
	}

	mg := mailgun.NewMailgun(m.MailgunDomain, m.MailgunAPIKey)

	msg := mg.NewMessage(m.SenderEmail, subject, armor, rcpt.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, id, err := mg.Send(ctx, msg)

	if err != nil {
		return "", err
	}

	return id, nil
}
