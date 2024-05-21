package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	instance := logger.With("component", "http/request")

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			msg := instance.With(
				slog.String("method", r.Method),
				slog.String("url", r.RequestURI),
				slog.String("addr", r.RemoteAddr),
				slog.String("agent", r.UserAgent()),
				slog.String("id", middleware.GetReqID(r.Context())),
			)

			startTime := time.Now()
			defer func() {
				msg.Info(
					"request completed",
					slog.String("elapsed", time.Since(startTime).String()),
				)
			}()

			wrapWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(wrapWriter, r)
		}

		return http.HandlerFunc(fn)
	}
}
