package handler

import (
	"ai-photo-app/types"
	"context"
	"net/http"
	"strings"
)

const userKey = "user"

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return 
		}
		user := types.AuthenticatedUser{
			Email: "mgg@gmail.com",
			LoggedIn: true,
		}
		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
