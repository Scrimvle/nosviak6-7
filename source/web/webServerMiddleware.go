package web

import (
	"Nosviak4/modules/gologr"

	"net/http"
	"github.com/gorilla/mux"
)

// MiddlewareFunc implements the logging interface for the middleware
func (s *Serve) MiddlewareFunc() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.Logger.WithTerminal().WriteLog(gologr.DEBUG, "[HTTP-SERVER:%s] Request to %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			h.ServeHTTP(w, r)
		})
	}
}
