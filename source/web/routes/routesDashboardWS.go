package routes

import (
	"Nosviak4/source"
	"Nosviak4/source/functions"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/sessions"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketDashboardHandler will allow us to concurrently update the information
func WebSocketDashboardHandler(writer http.ResponseWriter, request *http.Request) {
	token, err := functions.ExtractToken(request)
	if err != nil || token == nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	user, err := token.GetUser()
	if err != nil || user == nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	var socket = &websocket.Upgrader{}
	conn, err := socket.Upgrade(writer, request, nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	/* cache, acts as a cache map */
	cache := make(map[string]any)
	defer conn.Close()

	/* ranges through the ticker every 1 second */
	for range time.NewTicker(500 * time.Millisecond).C {
		users, err := database.DB.GetUsers()
		if err != nil || users == nil {
			users = make([]*database.User, 0)
		}

		attacks, err := database.DB.GetAttacks()
		if err != nil || attacks == nil {
			attacks = make([]*database.Attack, 0)
		}

		ongoingAttacks, err := database.DB.GetOngoing()
		if err != nil || ongoingAttacks == nil {
			ongoingAttacks = make([]*database.Attack, 0)
		}

		ongoingAttacksMethods, err := database.DB.GetOngoingMostUsedMethod()
		if err != nil || len(ongoingAttacksMethods) == 0 {
			ongoingAttacksMethods = "NOG"
		}

		data := make(map[string]any)

		/* metaData for the dashboard. */
		data["user"] = user
		data["users"] = len(users)
		data["attacks"] = len(attacks)
		data["methods"] = len(source.Methods)
		data["sessions"] = len(sessions.Sessions)
		data["ongoingAttacks"] = ongoingAttacks
		data["ongoingAttacksMostUsed"] = ongoingAttacksMethods

		/* authentication information for redirects. */
		data["session"] = hex.EncodeToString(sha256.New().Sum(token.Signature))
		if reflect.DeepEqual(data, cache) {
			continue
		}

		cache = data

		/* writes the data to the webSocket */
		if err := conn.WriteJSON(data); err != nil {
			fmt.Println(err)
			return
		}
	}
}
