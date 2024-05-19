package service

import (
	"bytes"
	"encoding/base64"
	"mail_microservice/internal/models"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"net/textproto"
	"strings"
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
			h = make(textproto.MIMEHeader)
			//h.Set("Content-Disposition", `form-data; name="file"; filename="`+fileData.Name+`"`)
			h.Set("Content-Disposition", `attachment;filename="`+fileData.Name+`"`)
			contentType := http.DetectContentType(fileData.Body)
			h.Set("Content-Type", contentType)
			h.Set("Content-Transfer-Encoding", "base64")

			part, err = writer.CreatePart(h)
			if err != nil {
				return err
			}

			encoded := base64.StdEncoding.EncodeToString(fileData.Body)
			_, err = part.Write([]byte(encoded))
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
	header["To"] = strings.Join(msg.To, ",")
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
	//TCPconn, err := net.DialTCP("TCP")
	//smtp.NewClient()
	err := smtp.SendMail(s.sender.Host+":"+s.sender.Port, s.auth, s.msg.FromAddress, s.msg.To, []byte(message))
	return err
}
