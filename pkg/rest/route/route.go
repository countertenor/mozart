package route

import (
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
func UIRouter() *mux.Router {
	router := mux.NewRouter()
	// statikFS, err := statik.GetStaticFS(statik.Webapp)
	// if err != nil {
	// 	log.Fatalf("could not get static files for UI, err : %v", err)
	// }
	// sub, _ := fs.Sub(static.Webapp, "static")
	// router.PathPrefix("/").Handler(http.FileServer(http.FS(sub)))

	fsys := fs.FS(static.Webapp)
	contentStatic, _ := fs.Sub(fsys, string(static.WebappBuild))

	router.PathPrefix("/").Handler(http.FileServer(http.FS(contentStatic)))
	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.Webapp))))
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Webapp))))
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
