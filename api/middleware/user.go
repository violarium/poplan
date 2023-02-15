package middleware

import (
	"context"
	"net/http"

	"github.com/violarium/poplan/user"
)

type User struct {
	userRegistry *user.Registry
}

func NewUser(userRegistry *user.Registry) *User {
	return &User{userRegistry: userRegistry}
}

func (m *User) AuthUserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if authUser, ok := m.userRegistry.Find(token); ok {
			ctx := context.WithValue(r.Context(), "authUser", authUser)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (m *User) RequireAuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Value("authUser").(*user.User); !ok {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}
