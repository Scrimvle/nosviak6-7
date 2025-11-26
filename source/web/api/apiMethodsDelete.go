package api

import (
	"Nosviak4/source"
	"Nosviak4/source/functions"
	"Nosviak4/source/database"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/nicklaw5/go-respond"
)

// MethodsDeleteEndpoint
func MethodsDeleteEndpoint(writer http.ResponseWriter, request *http.Request) {
	responder := respond.NewResponse(writer)
	token, err := functions.ExtractToken(request)
	if err != nil || token == nil {
		responder.Unauthorized(&MiddlewareReject{Status: false, Message: "unauthorized"})
		return
	}

	user, err := token.GetUser()
	if err != nil || user == nil || !database.DB.IsSuperuser(user) {
		responder.Unauthorized(&MiddlewareReject{Status: false, Message: "unauthorized"})
		return
	}

	clone := make(map[string]*source.Method)
	for name, method := range source.Methods {
		clone[name] = method
	}

	if _, ok := clone[mux.Vars(request)["method"]]; !ok {
		responder.BadRequest(&MiddlewareReject{Status: false, Message: "method not found"})
		return
	}

	delete(clone, mux.Vars(request)["method"])
	content, err := json.MarshalIndent(clone, "", "\t")
	if err != nil {
		responder.InternalServerError(&MiddlewareReject{Status: false, Message: "error building output"})
		return
	}

	if err := os.WriteFile(filepath.Join(source.ASSETS, "attacks", "apis.json"), content, 0777); err != nil {
		responder.InternalServerError(&MiddlewareReject{Status: false, Message: "error writing output"})
		return
	}

	source.Methods = clone
	responder.Ok(&MiddlewareReject{
		Status: true,
	})
}
