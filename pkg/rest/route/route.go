package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prashantgupta24/mozart/pkg/rest/handler"
	"github.com/prashantgupta24/mozart/pkg/spa"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

// //NewRouter creates a new mux router for application
// func NewRouter() *mux.Router {
// 	r := mux.NewRouter()
// 	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./webapp/build"))))
// 	//r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./webapp/build"))))

// 	r.PathPrefix("/").Handler(spa.Handler{StaticPath: "./webapp/build", IndexPath: "index.html"})
// 	return r
// }

func NewUIRouter() *mux.Router {
	router := mux.NewRouter()

	// ui := router.PathPrefix("/ui/")
	// ui.Handler(http.StripPrefix("", spa.Handler{StaticPath: "./webapp/build", IndexPath: "index.html"}))

	ui := router.PathPrefix("/")
	ui.Handler(spa.Handler{StaticPath: "./webapp/build", IndexPath: "index.html"})

	return router
}

//NewRouter creates a new mux router for application
func NewRouter() *mux.Router {

	router := mux.NewRouter()

	// ui := router.PathPrefix("/")
	// ui.Handler(spa.Handler{StaticPath: "./webapp/build", IndexPath: "index.html"})

	restServer := router.PathPrefix("/api/v1").Subrouter().StrictSlash(false)
	// router.PathPrefix("/api/v1").Handler(negroni.New(
	// 	negroni.NewRecovery(),
	// 	negroni.NewLogger(),
	// 	negroni.Wrap(subrouter),
	// ))

	// subrouter.Handle("/", handlers.LoggingHandler(logFile, finalHandler))
	// subrouter.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./webapp/build"))))
	// subrouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("webapp/build"))))
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
