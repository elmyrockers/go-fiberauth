package ses

import (
	sesmailer "github.com/elmyrockers/go-sesmailer"
	fiberauth "github.com/elmyrockers/go-fiberauth/session"
)

type Adapter struct {
	mailer *sesmailer.Mailer
}

// New() creates a fresh adapter wrapping a new sesmailer.Mailer instance.
func New() *Adapter {
	return &Adapter{mailer: sesmailer.New()}
}

func (a *Adapter) SetFrom(address fiberauth.MailAddress) {
	a.mailer.SetFrom(address.Email, address.Name)
}

func (a *Adapter) SetTo(address fiberauth.MailAddress) {
	a.mailer.AddTo(address.Email, address.Name)
}

func (a *Adapter) SetReplyTo(address fiberauth.MailAddress) {
	a.mailer.AddReplyTo(address.Email, address.Name)
}

func (a *Adapter) SetTemplate(template fiberauth.MailTemplate) {
	a.mailer.SetSubject(template.Subject)
	a.mailer.SetBody(template.Body)
	a.mailer.AsHTML()
	if template.AltBody != "" {
		a.mailer.SetAltBody(template.AltBody)
	}
}

func (a *Adapter) Send() error {
	_, err := a.mailer.Send()
	return err
}