package middleware

import (
	"context"
	"net/http"

	"github.com/violarium/poplan/api"
	"github.com/violarium/poplan/user"
)

const authUserKey = "authUser"

type UserMiddleware struct {
	userRegistry *user.Registry
}

func NewUserMiddleware(userRegistry *user.Registry) *UserMiddleware {
	return &UserMiddleware{userRegistry: userRegistry}
}

func (m *UserMiddleware) AuthUserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		authUser, ok := m.userRegistry.Find(token)
		if !ok {
			api.SendMessage(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authUserKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
