package route

import (
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/countertenor/mozart/pkg/rest/handler"
	"github.com/countertenor/mozart/statik"
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
	attachProfiler(router)
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

func attachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// Manually add support for paths linked to by index page at /debug/pprof/
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
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
