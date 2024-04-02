package service

import (
	"bytes"
	"mail_microservice/internal/models"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
)

func (s *Service) SendSMTPMessage(msg *models.Message) error {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", "text/plain")
	h.Set("Content-Transfer-Encoding", "base64")

	part, err := writer.CreatePart(h)
	if err != nil {
		return err
	}

	_, err = part.Write([]byte(msg.Text))
	if err != nil {
		return err
	}

	if msg.Files != nil {
		for _, fileData := range msg.Files {
			part, err = writer.CreateFormFile("file", fileData.Name)
			if err != nil {
				return err
			}

			_, err = part.Write([]byte(fileData.Body))
			if err != nil {
				return err
			}
		}
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	header := make(map[string]string)
	header["FromAddress"] = msg.FromName + " <" + msg.FromAddress + ">"
	header["To"] = msg.To
	header["Subject"] = msg.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "multipart/mixed; boundary=" + writer.Boundary()

	var message string
	for k, v := range header {
		message += k + ": " + v + "\n"
	}
	message += "\n" + body.String()

	s.msg = msg
	err = s.sendMessage(message)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) sendMessage(message string) error {
	err := smtp.SendMail(s.sender.Host+":"+s.sender.Port, s.auth, s.msg.FromAddress, []string{s.msg.To}, []byte(message))
	return err
}
