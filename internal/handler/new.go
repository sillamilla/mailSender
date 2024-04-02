package handler

import "mail_microservice/internal/service"

type Handler struct {
	srv service.Service
}

func New(srv *service.Service) Handler {
	return Handler{
		srv: *srv,
	}
}
