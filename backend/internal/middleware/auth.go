package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dhruv15803/budgeting-app/internal/auth"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func AuthMiddleware(jwtSigner *auth.JWTSigner) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeUnauthorized(w, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeUnauthorized(w, "invalid authorization header format")
				return
			}

			claims, err := jwtSigner.ParseAccessToken(parts[1])
			if err != nil {
				writeUnauthorized(w, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) *auth.Claims {
	claims, _ := ctx.Value(ClaimsKey).(*auth.Claims)
	return claims
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	//nolint:errcheck
	w.Write([]byte(`{"success":false,"message":"` + message + `","status_code":401}`))
}
