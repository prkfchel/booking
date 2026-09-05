package core_http_server

import (
	core_http_middleware "booking/internal/core/transport/http/middlerware"
	"net/http"
)

type Route struct {
	Method    string
	Path      string
	Handler   http.HandlerFunc
	MiddeWare []core_http_middleware.MiddleWare
}

func (r *Route) WithMiddlware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r.Handler,
		r.MiddeWare...,
	)
}
