package rest

import (
	"log"
	"net/http"
	"time"

	"github.com/prashantgupta24/mozart/pkg/rest/route"
)

//Port on which to start
const Port = "8080"

//StartServer starts the REST server
func StartServer() {

	log.Printf("Starting REST server at port %v ... \n", Port)

	router := route.NewRouter()

	s := &http.Server{
		Addr:           ":" + Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
