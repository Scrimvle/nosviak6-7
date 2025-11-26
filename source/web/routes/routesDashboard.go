package routes

import (
	"Nosviak4/source"
	"Nosviak4/source/functions"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"net/http"
	"path/filepath"
)

// DashboardPage is a static file which is open for the conns to use and customize but not add too.
var DashboardPage = template.Must(template.ParseFiles(filepath.Join("assets", "schemas", "web", "dashboard.html")))

// DashboardHandler will handle all the requests being made to the dashboard of Nosviak4.
func DashboardHandler(writer http.ResponseWriter, request *http.Request) {
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

	// data is the queue we store information in so it gets displayed and processed.
	data := make(map[string]any)

	/* handles any form submission we recv. */
	switch request.Method {

	/* sending broadcast or kicking sessions/session */
	case http.MethodPost:
		if err := request.ParseForm(); err != nil {
			break
		}

		switch request.FormValue("FormName") {

		/* broadcast event */
		case "Broadcast":
			/*
				We write our own implementation of the sessions.Broadcast function to incorporate
				the swash language on the message branding body.
			*/

			for _, session := range sessions.Sessions {
				payload, err := session.ExecuteBrandingToString(map[string]any{"sender": user.User(), "message": request.FormValue("BroadcastMessage")}, "web_broadcast_recv.tfx")
				if err != nil || len(payload) == 0 {
					continue
				}

				/* sends the message to the reader */
				session.Reader.PostAlert(&terminal.Alert{
					AlertCode:    terminal.MESSAGE,
					AlertMessage: payload,
				})
			}
		}
	}

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

	/* metaData for the dashboard. */
	data["user"] = user
	data["users"] = len(users)
	data["attacks"] = len(attacks)
	data["methods"] = len(source.Methods)
	data["sessions"] = sessions.Sessions
	data["ongoingAttacks"] = ongoingAttacks
	data["ongoingAttacksMostUsed"] = ongoingAttacksMethods

	/* authentication information for redirects. */
	data["session"] = hex.EncodeToString(sha256.New().Sum(token.Signature))

	// we just execute the dashboard to the write with the information it requires to render properly, this means
	// that the data is server rendered and therefore is another safety precaution.
	DashboardPage.Execute(writer, data)
}
