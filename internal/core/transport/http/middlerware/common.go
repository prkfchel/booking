package core_http_middleware

import (
	core_logger "booking/internal/core/logger"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
)

func RequestID() MiddleWare {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)
			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) MiddleWare {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestId),
				zap.String("url", r.URL.String()),
			)
			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() MiddleWare {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())

			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()

			log.Debug(
				">>> incoming http request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done http request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latenct", time.Now().Sub(before)),
			)
		})
	}
}

func Panic() MiddleWare {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

			defer func() {
				if r := recover(); r != nil {
					responseHandler.PanicResponse(r, "during handle HTTP request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
