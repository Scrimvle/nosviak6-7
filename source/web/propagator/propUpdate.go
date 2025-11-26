package propagator

import (
	"Nosviak4/source/database"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nicklaw5/go-respond"
	"golang.org/x/exp/slices"
)

// UpdatePropagate will update the item to the propagator buffer
func UpdatePropagate(w http.ResponseWriter, r *http.Request) {
	responder, fields := respond.NewResponse(w), mux.Vars(r)

	// decides what value we're changing
	if !r.URL.Query().Has("value") && r.Method == http.MethodPost {
		responder.BadRequest(map[string]any{"status": false, "message": "missing field"})
		return
	}

	token, ok := fields["key"]
	if !ok || len(token) == 0 {
		responder.BadRequest(map[string]any{"status": false, "message": "missing field"})
		return
	}

	propToken, ok := fields["prop_key"]
	if !ok || len(propToken) == 0 {
		responder.BadRequest(map[string]any{"status": false, "message": "missing field"})
		return
	}

	user, err := database.DB.GetUserAPIKey(token)
	if err != nil || !user.API || !slices.Contains(user.Roles, "admin") {
		responder.Unauthorized(map[string]any{"status": false, "message": "unauthorized"})
		return
	}

	appName, err := propagation.GetFromKey(propToken)
	if err != nil || len(appName) == 0 {
		responder.Unauthorized(map[string]any{"status": false, "message": "unauthorized"})
		return
	}

	switch r.Method {

	// Get value
	case http.MethodGet:
		mutex.RLock()
		defer mutex.RUnlock()

		val, ok := PropagatedFields[appName]
		if !ok {
			responder.BadRequest(map[string]any{"status": false, "message": "item not found"})
			return
		}

		responder.Ok(map[string]any{"status": true, "name": appName, "value": val})
		return

	// Edit/Update value
	case http.MethodPost:
		mutex.Lock()
		defer mutex.Unlock()

		PropagatedFields[appName] = r.URL.Query().Get("value")

	// Delete
	case http.MethodDelete:
		mutex.Lock()
		defer mutex.Unlock()

		if _, ok := PropagatedFields[appName]; !ok {
			responder.BadRequest(map[string]any{"status": false, "message": "item not found"})
			return
		}

		delete(PropagatedFields, appName)
		responder.Ok(map[string]any{"status": true, "message": "item deleted"})
		return

	default:
		responder.BadRequest(map[string]any{"status": false, "message": "unsupported method"})
		return
	}

	responder.Ok(map[string]any{"status": true, "name": appName, "value": PropagatedFields[appName]})
}