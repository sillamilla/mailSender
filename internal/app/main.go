package main

import (
	"fmt"
	"log"
	"mail_microservice/internal/config"
	"mail_microservice/internal/handler"
	"mail_microservice/internal/helper"
	"mail_microservice/internal/models"
	"mail_microservice/internal/service"
	"mail_microservice/internal/service/builder"
	"net/http"
	"net/smtp"
	"time"
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

	files := helper.GetFilesFromDirectory("./files")

	msg, err := builder.
		FromAddress("test@cock.li").
		FromName("test").
		To("test@cock.li").
		AddSubject("TEST").
		AddText("TEST MESSAGE!").
		AddFiles(files).
		Build()
	if err != nil {
		log.Printf("Error building message: %v", err)
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	counter := 0
	go func() {
		for range ticker.C {
			counter++
			fmt.Printf("\rSending...: %d", counter)
		}
	}()

	err = service.SendSMTPMessage(msg)
	if err != nil {
		log.Printf("Error sending SMTP message: %v", err)
		return
	}

	ticker.Stop()

	fmt.Println("\n message sent")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.WebPort.Port),
		Handler: handler.Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panicf("Error starting server: %v", err)
	}
}
