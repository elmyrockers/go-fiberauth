package session




type MailAdapter interface {
    SetFrom(address MailAddress)
    SetTo(address MailAddress)
    SetReplyTo(address MailAddress)

    SetTemplate(template MailTemplate)

    Send() error
}