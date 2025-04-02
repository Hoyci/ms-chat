package middlewares

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/hoyci/ms-chat/core/types"
	"github.com/hoyci/ms-chat/core/utils"
)

func AuthMiddleware(next http.Handler, publicKey *rsa.PublicKey) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(
				w,
				http.StatusUnauthorized,
				fmt.Errorf("user did not send an authorization header"),
				"AuthMiddleware",
				types.UnauthorizedResponse{Error: "Missing authorization header"},
			)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteError(
				w,
				http.StatusUnauthorized,
				fmt.Errorf("user sent an authorization header out of format"),
				"AuthMiddleware",
				types.UnauthorizedResponse{Error: "Invalid authorization header format"},
			)
			return
		}

		token := parts[1]

		claims, err := utils.VerifyJWT(token, publicKey)
		if err != nil {
			utils.WriteError(
				w,
				http.StatusUnauthorized,
				fmt.Errorf("user sent an invalid or expired authorization header"),
				"AuthMiddleware",
				types.UnauthorizedResponse{Error: "Invalid or expired token"},
			)
			return
		}

		ctx := r.Context()
		ctx = utils.SetClaimsToContext(ctx, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
