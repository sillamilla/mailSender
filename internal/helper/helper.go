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
	"os"
	"path/filepath"
	"strings"
)

const maxBytes = 10 << 20    // 10 MB
const maxFileSize = 26214400 // 25MB

func GetFilesFromDirectory(dirPath string) []models.File {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	var fileList []models.File
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dirPath, file.Name())

			fileInfo, err := os.Stat(filePath)
			if err != nil {
				log.Fatalf("Failed to get file info: %v", err)
			}

			if fileInfo.Size() > maxFileSize {
				log.Printf("Skipping file %s as it exceeds the maximum allowed size of %d bytes", file.Name(), maxFileSize)
				continue
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Failed to read file: %v", err)
			}

			fileList = append(fileList, models.File{
				Name: file.Name(),
				Body: content,
			})
		}
	}

	return fileList
}

func ReadForm(r *http.Request, data *models.Message) error {
	if err := r.ParseMultipartForm(maxBytes); err != nil {
		return err
	}

	data.FromAddress = r.FormValue("FromAddress")
	data.FromName = r.FormValue("FromName")
	data.To = strings.Split(r.FormValue("To"), ",")
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
			Body: content,
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

	address := fmt.Sprintf("%s <%s>", data.FromName, data.FromAddress)

	if len(address) > 320 || len(address) < 3 {
		errs = append(errs, errors.New("FROM: quantity of characters"))

	} else if _, err := mail.ParseAddress(address); err != nil { //todo replace with another validation
		errs = append(errs, errors.New("FROM: invalid email format"))
	}

	for _, reciver := range data.To {
		address = fmt.Sprintf(" <%s>", reciver)
		if len(address) > 320 || len(address) < 3 {
			errs = append(errs, errors.New("TO: quantity of characters"))
		} else if _, err := mail.ParseAddress(address); err != nil {
			errs = append(errs, errors.New("TO: invalid email format"))
		}
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

	if errs != nil {
		return errs
	}

	return nil
}
