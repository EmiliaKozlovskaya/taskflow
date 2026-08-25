package web_transport_http

import (
	"net/http"

	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
)

// http обработчик, который будет отдавать главную html страницу нашего веб-приложения
func (h *WebHTTPHandler) GetMainPage(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)

	//хотим отдать в ответе html файл, достанем его из уровня сервиса
	//html := service.GetMainPage()
	html, err := h.webService.GetMainPage()
	if err != nil {
		responseHandler.ErrorResponse(rw, err, "failed to get index.html for main page")
	}
	responseHandler.HTMLResponse(rw, html)
}
