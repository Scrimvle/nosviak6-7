package api

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	components "Nosviak4/source/functions"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nicklaw5/go-respond"
)

// ViewMethods will return a list of all the methods the api key holder can access.
func ViewMethods(writer http.ResponseWriter, request *http.Request) {
	responder := respond.NewResponse(writer)
	user, err := database.DB.GetUserAPIKey(mux.Vars(request)["key"])
	if err != nil || !user.API {
		responder.Unauthorized(&MiddlewareReject{Status: false, Message: "Unauthorized"})
		return
	}

	// Methods stores all the methods they can access with some information too.
	Methods := make(map[string]map[string]any)

	for method, object := range source.Methods {
		if !components.CanAccessThemPermissions(user, object.Access...) || !object.Options.API {
			continue
		}

		// Inserts the information about the method
		Methods[method] = map[string]any{
			"description": object.Description,
			"ipAllowed":   object.IPAllowed,
			"urlAllowed":  object.URLAllowed,
		}
	}

	content, err := json.MarshalIndent(Methods, "", "\t")
	if err != nil {
		responder.ServiceUnavailable(&MiddlewareReject{Status: false, Message: "Error"})
		return
	}

	if _, err := writer.Write([]byte(content)); err != nil {
		responder.ServiceUnavailable(&MiddlewareReject{Status: false, Message: "Error"})
		return
	}
}
