package core_http_response

import "net/http"

var (
	StatusCodeUninitialized = -1 //чтобы не было конфликта с реальными статус кодами, которые всегда >= 100, используем -1 как значение по умолчанию
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int //нам этот статус код нужно получать, но публичным это поле делать не будем, а напишем отдельный метод для получения
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

// переопределяем метод WriteHeader, чтобы сохранять статус код ответа в поле statusCode
// метод WriteHeader отввечает не за хэдеры, а конкретно за статус код ответа !!!
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("status code is uninitialized")
	}
	return rw.statusCode
}
