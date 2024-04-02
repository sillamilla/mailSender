package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mail_microservice/internal/models"
	"net/http"
	"net/mail"
	"strings"
)

const maxBytes = int64(1048576) // 1GB
const maxFileSize = 26214400    // 25MB

func ReadForm(r *http.Request, data *models.Message) error {
	if err := r.ParseMultipartForm(maxBytes); err != nil {
		return err
	}

	data.FromAddress = r.FormValue("FromAddress")
	data.FromName = r.FormValue("FromName")
	data.To = r.FormValue("To")
	data.Subject = r.FormValue("Subject")
	data.Text = r.FormValue("Text")

	fileHeaders := r.MultipartForm.File["Files"]
	for _, header := range fileHeaders {
		file, err := header.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		newFile := models.File{
			Name: header.Filename,
			Body: string(content),
		}

		data.Files = append(data.Files, newFile)
	}

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data *models.Message) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	out, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		log.Println(err)
	}
}

func ErrorJSON(w http.ResponseWriter, err error) {
	statusCode := http.StatusBadRequest

	var payload models.JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	WriteJSON(w, statusCode, payload)
}

type customErr []error

func (e customErr) Error() string {
	builder := strings.Builder{}
	for _, err := range e {
		builder.WriteString(err.Error() + "\n")
	}

	return builder.String()
}

func Validate(data *models.Message) error {
	var errs customErr

	if len(data.FromAddress) > 320 || len(data.FromAddress) < 3 {
		errs = append(errs, errors.New("FROM: quantity of characters"))
	} else if _, err := mail.ParseAddress(data.FromAddress); err != nil { //todo replace with another validation
		errs = append(errs, errors.New("FROM: invalid email format"))
	}

	if len(data.To) > 320 || len(data.To) < 3 {
		errs = append(errs, errors.New("TO: quantity of characters"))
	} else if _, err := mail.ParseAddress(data.To); err != nil {
		errs = append(errs, errors.New("TO: invalid email format"))
	}

	if len(data.Subject) > 78 {
		errs = append(errs, errors.New("SUBJECT: quantity of characters"))
	}

	if len(data.Text) > 1000000 {
		errs = append(errs, errors.New("TEXT: quantity of characters"))
	}

	for i, file := range data.Files {
		if len(file.Body) > maxFileSize {
			errs = append(errs, fmt.Errorf("FILES[%d]: file size", i))
		}
	}

	return errs
}
