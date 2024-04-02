# README

`mailSender` is a Go library for sending emails with attachments via SMTP. It also supports sending emails via form data and JSON.

### Installation

To install `mailSender`, you need to have Go installed on your machine. Once you have Go set up, you can fetch the library using `go get`:

```bash
go get https://github.com/sillamilla/mailSender.git
```

### Usage

Here is a basic example of how to use the `mailSender` library:

```go
package main

import (
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
	
	// You can use different authorization methods from standard library (CRAMMD5Auth)
    auth := smtp.PlainAuth("", sender.Username, sender.Password, sender.Host)

    service := service.New(&sender, auth)
    handler := handler.New(service)
    builder := builder.New()

    // Build the message
    msg, err := builder.
        FromAddress("fromTest@gmail.com").
        FromName("Ilia").
        To("test1@gmail.com", "test2@gmail.com").
        AddSubject("Test message").
        AddText("Hello world!").
        AddFiles(nil).
        Build()
    if err != nil {
        log.Printf("Error building message: %v", err)
        return
    }

    // Send the message
    err = service.SendSMTPMessage(msg)
    if err != nil {
        log.Printf("Error sending SMTP message: %v", err)
        return
    }

    // Start the server
    srv := &http.Server{
        Addr:    fmt.Sprintf(":%s", cfg.WebPort.Port),
        Handler: handler.Routes(),
    }

    err = srv.ListenAndServe()
    if err != nil {
        log.Panicf("Error starting server: %v", err)
    }
}
```

This example shows how to use the `mailSender` library to send an email with attachments. The `SenderCredentials` struct is used to provide the SMTP server details and the sender's email credentials. The `builder` is used to construct the email message, and the `service` is used to send the email.

In addition to the above, `mailSender` also supports sending emails via form data and JSON. This can be done by using the `ReadForm` and `ReadJSON` helper functions respectively. These functions read the email details from the form data or JSON and use them to send the email.

Please refer to the library's API documentation for more detailed information on how to use `mailSender`.