package routes

import (
	"Nosviak4/source/database"
	"Nosviak4/source/functions"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/dgrijalva/jwt-go"
)

// LoginPage is a static file which is open for the conns to use and customize but not add too.
var LoginPage = template.Must(template.ParseFiles(filepath.Join("assets", "schemas", "web", "login.html")))

// LoginHandler will handle the login which has been received from the connection
func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	data := make(map[string]any)

	switch request.Method {

	case http.MethodPost:
		if err := request.ParseForm(); err != nil {
			data["error"] = "Client error occurred. Please try again."
			break
		}

		user, err := database.DB.GetUser(request.FormValue("Username"))
		if err != nil || user == nil {
			data["error"] = "Incorrect username or password. Please try again."
			break
		}

		// Tries to cross-reference the users password to the password provided.
		if !user.IsPassword([]byte(request.FormValue("Password"))) {
			data["error"] = "Incorrect username or password. Please try again."
			break
		}

		if !functions.CanAccessThemPermissions(user, "admin") {
			data["error"] = "Missing admin permission. Please try again."
			break
		}

		token, err := functions.NewToken(jwt.MapClaims{"user": user.Username})
		if err != nil || token == nil {
			data["error"] = "Server error occurred. Please try again."
			break
		}

		http.SetCookie(writer, token.Cookie)
		http.Redirect(writer, request, fmt.Sprintf("/%s/", hex.EncodeToString(sha256.New().Sum(token.Signature))), http.StatusFound)
		return
	}

	LoginPage.Execute(writer, data)
}
