package serviceMiddleware

import (
	"blog-service/internal/http-server/api/response"
	"blog-service/internal/lib"
	"blog-service/pkg/logger"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	tokenManager lib.TokenManager
}

func NewAuthMiddleware(tokenManager lib.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tokenManager}
}
func (a *AuthMiddleware) CheckAuthorizationToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, response.Error("token is empty"))
		}

		token = strings.TrimPrefix(token, "Bearer ")
		decodedUserId, err := a.tokenManager.Parse(token)
		if err != nil {
			fmt.Printf("failed to parse token: %s", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to parse token"))
		} else {
			r.Header.Set("user_id", decodedUserId)
			next.ServeHTTP(w, r)
		}
	})
}
