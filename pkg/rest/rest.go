package rest

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/countertenor/mozart/pkg/rest/route"
	"github.com/rs/cors"
)

//Ports on which to start
const restPort = "8080"
const uiPort = "8081"

//StartServer starts the REST and UI server
func StartServer() {
	restRouter := route.RestRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:" + uiPort},
		AllowCredentials: true,
	})
	restServer := &http.Server{
		Addr:           ":" + restPort,
		Handler:        c.Handler(restRouter),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	uiRouter := route.UIRouter()
	uiServer := &http.Server{
		Addr:           ":" + uiPort,
		Handler:        uiRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		log.Fatal(restServer.ListenAndServe())
		wg.Done()
	}()

	go func() {
		log.Fatal(uiServer.ListenAndServe())
		wg.Done()
	}()

	log.Printf("Started REST server at port %v ... \n", restPort)
	log.Printf("Started UI server at port %v ... \n", uiPort)
	wg.Wait()
}
