package core_http_middleware

import "net/http"

// в этом файле создаем алиас для типа функции, которые будем использовать как middleware в разных местах приложения
type Middleware func(http.Handler) http.Handler

func ChainMiddleware(
	h http.Handler,
	m ...Middleware,
) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- { //чтобы сверху навешивать внешние middleware
		h = m[i](h)
	}
	return h
}
