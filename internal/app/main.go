package main

import (
	"fmt"
	"log"
	"mail_microservice/internal/config"
	"mail_microservice/internal/handler"
	"mail_microservice/internal/models"
	"mail_microservice/internal/service"
	"mail_microservice/internal/service/builder"
	"net/http"
	"net/smtp"
	"os"
)

func main() {
	cfg := config.GetConfig()

	sender := models.SenderCredentials{
		Host:       cfg.Mail.Host,
		Port:       cfg.Mail.Port,
		Username:   cfg.Mail.Username,
		Password:   cfg.Mail.Password,
		Encryption: cfg.Mail.Encryption,
	}
	auth := smtp.PlainAuth("", sender.Username, sender.Password, sender.Host)

	service := service.New(&sender, auth)
	handler := handler.New(service)
	builder := builder.New()

	videoData, err := os.ReadFile("files/video.mp4")
	if err != nil {
		log.Fatal(err)
	}
	txtData, err := os.ReadFile("files/file.txt")
	if err != nil {
		log.Fatal(err)
	}
	files := []models.File{
		{Name: "video.mp4", Body: string(videoData)},
		{Name: "file.txt", Body: string(txtData)},
	}

	msg, err := builder.
		FromAddress("test@gmail.com").
		FromName("Ilia").
		To("test1@gmail.com", "test2@gmail.com").
		AddSubject("Test").
		AddText("TEST MESSAGE!").
		AddFiles(files).
		Build()
	if err != nil {
		log.Printf("Error building message: %v", err)
		return
	}
	err = service.SendSMTPMessage(msg)
	if err != nil {
		log.Printf("Error sending SMTP message: %v", err)
		return
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.WebPort.Port),
		Handler: handler.Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panicf("Error starting server: %v", err)
	}
}
