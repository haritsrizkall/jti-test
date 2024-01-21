package middlewares

import (
	"context"
	"net/http"

	"github.com/haritsrizkall/jti-test/pkg"
)

type ContextKey int

const (
	UserID ContextKey = iota
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// get token from cookie
		token, err := req.Cookie("token")
		if err != nil {
			resp.WriteHeader(http.StatusUnauthorized)
			resp.Header().Set("Content-Type", "application/json")
			resp.Write([]byte("Unauthorized"))
			return
		}

		// validate token
		payload, err := pkg.ExtractPayload(token.Value)
		if err != nil {
			resp.WriteHeader(http.StatusUnauthorized)
			resp.Header().Set("Content-Type", "application/json")
			resp.Write([]byte("Unauthorized"))
			return
		}

		// set user id to context
		ctx := req.Context()
		newCtx := context.WithValue(ctx, UserID, payload.UserID)

		next.ServeHTTP(resp, req.WithContext(newCtx))
	})
}
