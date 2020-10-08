package rest

import (
	"log"
	"net/http"
	"time"

	"github.com/prashantgupta24/mozart/pkg/rest/route"
	"github.com/rs/cors"
)

//Port on which to start
const Port = "8080"

//StartServer starts the REST server
func StartServer() {

	log.Printf("Starting REST server at port %v ... \n", Port)

	router := route.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowCredentials: true,
	})

	s := &http.Server{
		Addr:           ":" + Port,
		Handler:        c.Handler(router),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	uiRouter := route.NewUIRouter()

	s1 := &http.Server{
		Addr:           ":8081",
		Handler:        uiRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Fatal(s1.ListenAndServe())
	}()
	log.Fatal(s.ListenAndServe())

	// statikFS, err := fs.New()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // Serve the contents over HTTP.
	// http.Handle("/", http.FileServer(statikFS))
	// http.ListenAndServe(":8080", nil)
}
