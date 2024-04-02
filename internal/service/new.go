package service

import (
	"mail_microservice/internal/models"
	"net/smtp"
)

type Service struct {
	sender *models.SenderCredentials
	msg    *models.Message
	auth   smtp.Auth
}

func New(sender *models.SenderCredentials, auth smtp.Auth) *Service {
	return &Service{
		sender: sender,
		auth:   auth,
	}
}
