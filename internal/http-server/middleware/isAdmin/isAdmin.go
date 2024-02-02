package isAdmin

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/http-server/middleware/auth"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	const op = "middleware.isAdmin.New"
    log = log.With(slog.String("op", op))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			isAdmin, ok := r.Context().Value(auth.IsAdminKey).(bool)

			if !ok || !isAdmin {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}