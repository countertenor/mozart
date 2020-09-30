package rest

import (
	"log"
	"net/http"
	"time"

	"github.com/prashantgupta24/mozart/pkg/rest/route"
)

//StartServer starts the REST server
func StartServer() {

	log.Println("Starting REST server ...")

	router := route.NewRouter()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
