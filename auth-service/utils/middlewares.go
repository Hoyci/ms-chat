package utils

import (
	"fmt"
	"net/http"
	"strings"

	coreTypes "github.com/hoyci/ms-chat/core/types"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			coreUtils.WriteError(
				w,
				http.StatusUnauthorized,
				fmt.Errorf("user did not send an authorization header"),
				"AuthMiddleware",
				coreTypes.UnauthorizedResponse{Error: "Missing authorization header"},
			)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			coreUtils.WriteError(
				w,
				http.StatusUnauthorized,
				fmt.Errorf("user sent an authorization header out of format"),
				"AuthMiddleware",
				coreTypes.UnauthorizedResponse{Error: "Invalid authorization header format"},
			)
			return
		}

		token := parts[1]

		claims, err := VerifyJWT(token, &PrivateKeyAccess.PublicKey)
		if err != nil {
			coreUtils.WriteError(
				w,
				http.StatusUnauthorized,
				fmt.Errorf("user sent an invalid or expired authorization header"),
				"AuthMiddleware",
				coreTypes.UnauthorizedResponse{Error: "Invalid or expired token"},
			)
			return
		}

		ctx := r.Context()
		ctx = SetClaimsToContext(ctx, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
