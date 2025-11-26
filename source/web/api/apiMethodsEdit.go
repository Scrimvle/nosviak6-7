package api

import (
	"Nosviak4/source"
	"Nosviak4/source/functions"
	"Nosviak4/source/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nicklaw5/go-respond"
)

// MethodsEditEndpoint will edit the method based on the params
func MethodsEditEndpoint(writer http.ResponseWriter, request *http.Request) {
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

	content, err := io.ReadAll(request.Body)
	if err != nil {
		responder.BadRequest(&MiddlewareReject{Status: false, Message: "missing body"})
		return
	}

	method := new(source.Method)
	defer request.Body.Close()
	if err := json.Unmarshal(content, &method); err != nil {
		responder.BadRequest(&MiddlewareReject{Status: false, Message: "bad body"})
		return
	}

	clone := make(map[string]*source.Method)
	for name, method := range source.Methods {
		clone[name] = method
	}

	fmt.Println(clone)
}
