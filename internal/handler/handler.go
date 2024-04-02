package handler

import (
	"fmt"
	"log"
	"mail_microservice/internal/helper"
	"mail_microservice/internal/models"
	"mail_microservice/internal/service/builder"
	"net/http"
)

func (h *Handler) SendMail(w http.ResponseWriter, r *http.Request) {
	var requestPayload models.Message
	var err error

	switch r.Header.Get("Content-Type") {
	case "application/json":
		err = helper.ReadJSON(w, r, &requestPayload)
		if err != nil {
			log.Println(err)
			helper.ErrorJSON(w, err)
			return
		}
	case "multipart/form-data":
		err = helper.ReadForm(r, &requestPayload)
		if err != nil {
			log.Println(err)
			helper.ErrorJSON(w, err)
			return
		}
	default:
		http.Error(w, "Invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	builder := builder.MessageBuilder{}
	msg, err := builder.BuildMessage().
		FromAddress(requestPayload.FromAddress).
		FromName(requestPayload.FromName).
		To(requestPayload.To).
		AddSubject(requestPayload.Subject).
		AddText(requestPayload.Text).
		AddFiles(requestPayload.Files).
		Build()
	if err != nil {
		log.Println(err)
		helper.ErrorJSON(w, err)
		return
	}

	err = h.srv.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		helper.ErrorJSON(w, err)
		return
	}

	payload := models.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("sent to %s", requestPayload.To),
	}

	helper.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		log.Println(err)
		helper.ErrorJSON(w, err)
		return
	}
}
