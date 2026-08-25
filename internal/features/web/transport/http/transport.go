package web_transport_http

import (
	core_http_server "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/server"
)

type WebHTTPHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
}

func NewWebHTTPHandler(
	webService WebService,
) *WebHTTPHandler {
	return &WebHTTPHandler{
		webService: webService,
	}
}

func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Path:    "/", //на корневом паттерне отдаём главную страницу нашего веб-приложения
			Handler: h.GetMainPage,
		},
	}
}
