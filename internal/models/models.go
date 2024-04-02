package models

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type SenderCredentials struct {
	Domain     string
	Host       string
	Port       string
	Username   string
	Password   string
	Encryption string
}

type Message struct {
	FromAddress string `json:"from_address"`
	FromName    string `json:"from_name"`
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Text        string `json:"text"`
	Files       []File `json:"files"`
}

type File struct {
	Name string `json:"name"`
	Body string `json:"body"`
}
