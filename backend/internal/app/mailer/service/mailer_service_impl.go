package service

import "app/internal/app/mailer/model"

type MailProvider interface {
	Send(to, subject, body string) error
}

type mailerService struct {
	provider MailProvider
}

func NewMailerService(p MailProvider) MailerService {
	return &mailerService{provider: p}
}

func (m *mailerService) Send(mail model.Mail) error {
	return m.provider.Send(mail.To, mail.Subject, mail.Body)
}
