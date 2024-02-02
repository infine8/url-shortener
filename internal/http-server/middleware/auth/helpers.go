package auth

import (
	"context"
	"net/http"
	"strings"
)

const (
	UidKey     = "uid"
	ErrorKey   = "error"
	IsAdminKey = "is_admin"
)

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func UIDFromContext(ctx context.Context) (int64, bool) {
	uid, ok := ctx.Value(UidKey).(int64)
	return uid, ok
}

func ErrorFromContext(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(ErrorKey).(error)
	return err, ok
}
