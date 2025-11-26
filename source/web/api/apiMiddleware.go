package api

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source/database"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nicklaw5/go-respond"
)

// MiddlewareReject is the rejection response reporter to the request, all
// of our responses should inherit the 2 fields inside the rejector to ensure
// we have some sort of consistency within the responses, making sure all fields
// are omitted is critical too.
type MiddlewareReject struct {
	Status  bool   `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// HandleAuthMiddleware will ensure the incoming request is authenticated,
// every single API route must have a key variable associated within the URL.
func HandleAuthMiddleware(logger *gologr.FileLogger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := database.DB.GetUserAPIKey(mux.Vars(r)["key"])
			if err != nil || !user.API {
				respond.NewResponse(w).Unauthorized(&MiddlewareReject{Status: false, Message: "Unauthorized"})
				return
			}

			// implement's antiPanic into the middleware for the API, stops API panics from crashing the entire CNC.
			defer func() {
				err := recover()
				if err == nil {
					return
				}
		
				logger.WithTerminal().WriteLog(gologr.ERROR, "Panic caught for %s@%s: %v", user.Username, r.RemoteAddr, err)
			}()

			logger.WithTerminal().WriteLog(gologr.DEBUG, "authorized %s: %s", user.Username, r.URL.String())

			// continues through the processing stage, authentication has been completed.
			h.ServeHTTP(w, r)
		})
	}
}