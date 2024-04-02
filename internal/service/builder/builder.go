package builder

import (
	"encoding/base64"
	"mail_microservice/internal/helper"
	"mail_microservice/internal/models"
	"strings"
)

type MessageBuilder struct {
	msg *models.Message
}

func New() *MessageBuilder {
	return &MessageBuilder{
		msg: &models.Message{},
	}

}

func (b *MessageBuilder) FromAddress(email string) *MessageBuilder {
	b.msg.FromAddress = email
	return b
}

func (b *MessageBuilder) FromName(name string) *MessageBuilder {
	b.msg.FromName = name
	return b
}

func (b *MessageBuilder) To(emails ...string) *MessageBuilder {
	b.msg.To = strings.Join(emails, ", ") //todo coma check at the end
	return b
}

func (b *MessageBuilder) AddSubject(subject string) *MessageBuilder {
	b.msg.Subject = subject
	return b
}

func (b *MessageBuilder) AddText(text string) *MessageBuilder {
	b.msg.Text = base64.StdEncoding.EncodeToString([]byte(text))
	return b
}

func (b *MessageBuilder) AddFiles(files []models.File) *MessageBuilder {
	for _, file := range files {
		encodedFile := models.File{
			Name: file.Name,
			Body: base64.StdEncoding.EncodeToString([]byte(file.Body)),
		}

		b.msg.Files = append(b.msg.Files, encodedFile)
	}

	return b
}

func (b *MessageBuilder) Build() (*models.Message, error) {
	err := helper.Validate(b.msg)
	if err != nil {
		return nil, err
	}

	return b.msg, nil
}
