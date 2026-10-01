package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/stevanusy21/golang_sandbox/pkg/response"
	"github.com/stevanusy21/golang_sandbox/pkg/token"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

type contextKey string

const userIdKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			response.Error(w, http.StatusUnauthorized, utils.ErrTokenMissing.Error())
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, http.StatusUnauthorized, utils.ErrTokenInvalid.Error())
			return
		}

		claims, err := token.ValidateToken(parts[1])
		if err != nil {
			response.Error(w, http.StatusUnauthorized, fmt.Sprintf("%s: %v", utils.ErrTokenInvalid, err))
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIdFromContext(r *http.Request) (int, bool) {
	userId, ok := r.Context().Value(userIdKey).(int)
	return userId, ok
}
