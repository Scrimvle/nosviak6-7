package propagator

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nicklaw5/go-respond"
	"golang.org/x/exp/slices"
)

type propagator struct {
	Enabled bool `toml:"enabled"`
	Key     []struct {
		Name string `toml:"name"`
		Key  string `toml:"key"`
	} `toml:"key"`
}

// GetFromKey will attempt to find the propagator key from the config
func (p *propagator) GetFromKey(key string) (string, error) {
	for _, name := range p.Key {
		if name.Key != key {
			continue
		}

		return name.Name, nil
	}

	return "", errors.New("key not found")
}

// Propagation stores all the fields required
var propagation *propagator = new(propagator)

// PropagateMiddleware authenticates the middleware
func PropagateMiddleware(logger *gologr.FileLogger) mux.MiddlewareFunc {
	if err := source.OPTIONS.MarshalFromPath(&propagation, "web", "propagator"); err != nil {
		propagation = &propagator{
			Enabled: false,
		}
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responder, fields := respond.NewResponse(w), mux.Vars(r)

			token, ok := fields["key"]
			if !ok || len(token) == 0 {
				responder.BadRequest(map[string]any{"status": false, "message": "missing field: key"})
				return
			}

			propToken, ok := fields["prop_key"]
			if !ok || len(propToken) == 0 {
				responder.BadRequest(map[string]any{"status": false, "message": "missing field: prop_key"})
				return
			}

			user, err := database.DB.GetUserAPIKey(token)
			if err != nil || !user.API || !slices.Contains(user.Roles, "admin") {
				responder.Unauthorized(map[string]any{"status": false, "message": "unauthorized: key"})
				return
			}

			appName, err := propagation.GetFromKey(propToken)
			if err != nil || len(appName) == 0 {
				responder.Unauthorized(map[string]any{"status": false, "message": "unauthorized: prop_key"})
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
