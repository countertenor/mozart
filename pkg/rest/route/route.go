package route

import (
	"errors"
	"io/fs"
	"net/http"

	"github.com/countertenor/mozart/pkg/rest/handler"
	"github.com/countertenor/mozart/static"
	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

//UIRouter creates a router for the UI
func UIRouter() (*mux.Router, error) {
	router := mux.NewRouter()
	if static.WebappBuildType == "" {
		return nil, errors.New("ui not included")
	}
	fsys := fs.FS(static.GetEmbedFS(static.WebappBuildType))
	contentStatic, _ := fs.Sub(fsys, string(static.WebappBuildType))
	router.PathPrefix("/").Handler(http.FileServer(http.FS(contentStatic)))
	return router, nil
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
