package service

import "app/internal/app/mailer/model"

type MailerService interface {
	Send(mail model.Mail) error
}
