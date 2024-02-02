package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/grpc"
	"url-shortener/lib/logger/sl"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrFailedIsAdminCheck = errors.New("failed to check if user is admin")
)

func New(log *slog.Logger, jwtAppSecret string, grpcClient *grpc.Client) func(next http.Handler) http.Handler {
	const op = "middleware.auth.New"

	log = log.With(slog.String("op", op))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtAppSecret), nil
			})

			if err != nil {
				log.Warn("failed to parse token", sl.Err(err))

				ctx := context.WithValue(r.Context(), ErrorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			mapClaims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				ctx := context.WithValue(r.Context(), ErrorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			log.Info("user authorized", slog.Any("claims", mapClaims))

			uid, ok := mapClaims["uid"].(float64)
			if !ok {
				ctx := context.WithValue(r.Context(), ErrorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			isAdmin, err := grpcClient.IsAdmin(r.Context(), int64(uid))
			if err != nil {
				ctx := context.WithValue(r.Context(), ErrorKey, ErrFailedIsAdminCheck)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			ctx := context.WithValue(r.Context(), UidKey, uid)
			ctx = context.WithValue(r.Context(), IsAdminKey, isAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
