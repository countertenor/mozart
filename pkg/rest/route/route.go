package route

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prashantgupta24/mozart/pkg/rest/handler"
	"github.com/prashantgupta24/mozart/statik"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

//UIRouter creates a router for the UI
func UIRouter() *mux.Router {
	router := mux.NewRouter()
	statikFS, err := statik.GetStaticFS(statik.Webapp)
	if err != nil {
		log.Fatalf("could not get static files for UI, err : %v", err)
	}
	router.PathPrefix("/").Handler(http.FileServer(statikFS))
	return router
}

//RestRouter creates a new mux router for application
func RestRouter() *mux.Router {
	router := mux.NewRouter()
	restServer := router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)

	restServer.Use(loggingMiddleware)
	restServer.Use(panicHandlerMiddleware)
	for _, route := range routesForApp {
		restServer.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

var routesForApp = routes{
	route{
		"Index Page",
		"GET",
		"/",
		handler.Index,
	},
	route{
		"Web socket",
		"GET",
		"/ws",
		handler.WebSocket,
	},
	route{
		"Get modules",
		"GET",
		"/modules",
		handler.GetModules,
	},
	route{
		"Configuration Toggle",
		"POST",
		"/config",
		handler.Configuration,
	},
	route{
		"Execute dir",
		"POST",
		"/execute",
		handler.ExecuteDir,
	},
	route{
		"Get state",
		"GET",
		"/state",
		handler.GetState,
	},
	route{
		"Cancel execution",
		"PUT",
		"/cancel",
		handler.Cancel,
	},
}
