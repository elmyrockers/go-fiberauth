package session


type MailAddress struct {
    Name  string
    Email string
}

type MailTemplate struct {
    Subject string
    Body    string
    AltBody string
}

type MailAdapter interface {
    SetFrom(address MailAddress)
    SetTo(address MailAddress)
    SetReplyTo(address MailAddress)

    SetTemplate(template MailTemplate)

    Send() error
}