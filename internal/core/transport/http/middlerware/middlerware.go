package core_http_middleware

import "net/http"

type MiddleWare func(http.Handler) http.Handler

func ChainMiddleware(
	h http.Handler,
	m ...MiddleWare,
) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}
