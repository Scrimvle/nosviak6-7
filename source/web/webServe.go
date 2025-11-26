package web

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/web/api"
	"Nosviak4/source/web/propagator"
	"Nosviak4/source/web/routes"
	"Nosviak4/source/web/routes/middleware"
	"path/filepath"

	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

// Serve maintains all the controller and services required for the web application
type Serve struct {
	Bind   *http.Server
	Config *ServeConfig
	Router *mux.Router
	Logger *gologr.FileLogger
}

// ServeConfig is polled from GoConfig for the appropriate configuration
type ServeConfig struct {
	Enabled bool   `json:"enabled"`
	Address string `json:"address"`
	Port    int    `json:"port"`
	SSL     struct {
		Enabled bool   `json:"enabled"`
		Cert    string `json:"cert"`
		Key     string `json:"key"`
	} `json:"ssl"`
}

// NewWebServe will build our Serve and initialize the server entirely
func NewWebServe() error {
	serve := new(Serve)
	err := source.OPTIONS.MarshalFromPath(&serve.Config, "web")
	if err != nil || !serve.Config.Enabled {
		return err
	}

	serve.Logger = source.LOGGER.NewFileLogger(filepath.Join(source.ASSETS, "logs", "web.log"), int64(source.OPTIONS.Ints("branding", "recycle_log")))
	if serve.Logger.Err != nil {
		return serve.Logger.Err
	}

	return serve.bind()
}

// bind will bind to the http port
func (s *Serve) bind() error {
	errRecv := make(chan error)
	go func() {
		s.Router = mux.NewRouter()
		s.Bind = &http.Server{
			Handler: s.Router,
			Addr:    fmt.Sprintf("%s:%d", s.Config.Address, s.Config.Port),
		}

		time.AfterFunc(1*time.Second, func() {
			errRecv <- nil
			source.LOGGER.AggregateTerminal().WriteLog(gologr.DEFAULT, "[Successfully initialized the web server] (%s)", s.Bind.Addr)
		})

		/* listens based on what type we're using */
		if s.Config.SSL.Enabled {
			errRecv <- s.Bind.ListenAndServeTLS(s.Config.SSL.Cert, s.Config.SSL.Key)
		} else {
			errRecv <- s.Bind.ListenAndServe()
		}
	}()

	defer func() {
		s.Router.Use(s.MiddlewareFunc())
		s.Router.HandleFunc("/login", routes.LoginHandler)
		s.Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/login", http.StatusFound)
		})

		/* requests for this handler are required to be authenticated */
		authServed := s.Router.PathPrefix("/{id}/").Subrouter()
		authServed.Use(middleware.ServeMiddlewareForAuth(s.Logger))

		/* methods api endpoints through the session */
		authServed.HandleFunc("/api/attacks/edit/{method}", api.MethodsEditEndpoint)
		authServed.HandleFunc("/api/attacks/delete/{method}", api.MethodsDeleteEndpoint)

		/* dashboard and websocket handler */
		authServed.HandleFunc("/", routes.DashboardHandler)
		authServed.HandleFunc("/ws", routes.WebSocketDashboardHandler)

		/* attacks */
		authServed.HandleFunc("/attacks/apis", routes.AttacksAPI)

		/* API handler */
		apiRouter := s.Router.PathPrefix("/api/{key}").Subrouter()
		apiRouter.Use(api.HandleAuthMiddleware(s.Logger))
		apiRouter.HandleFunc("/attack", api.LaunchAPIAttack)
		apiRouter.HandleFunc("/methods", api.ViewMethods)

		prop := apiRouter.PathPrefix("/propagator/{prop_key}/").Subrouter()
		prop.Use(propagator.PropagateMiddleware(s.Logger))
		prop.HandleFunc("/", propagator.UpdatePropagate)

		url := url.URL{
			Host:   s.Bind.Addr,
			Scheme: "http",
			Path:   "/login",
		}

		source.LOGGER.AggregateTerminal().WriteLog(gologr.ALERT, "- Web server is now available: %s", url.String())
	}()

	return <-errRecv
}
