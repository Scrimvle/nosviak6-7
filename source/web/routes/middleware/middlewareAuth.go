package middleware

import (
	"Nosviak4/modules/gologr"
	components "Nosviak4/source/functions"
	"net/http"

	"github.com/gorilla/mux"
)

// ServeMiddlewareForAuth will handle the middleware authentication for the wall
func ServeMiddlewareForAuth(logger *gologr.FileLogger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := components.ExtractToken(r)
			if err != nil || token == nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			user, err := token.GetUser()
			if err != nil || user == nil || !components.CanAccessThemPermissions(user, "admin") {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
