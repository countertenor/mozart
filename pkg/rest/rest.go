package rest

import (
	"fmt"
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

	doIncludeUI := true
	var uiServer *http.Server

	restRouter := route.RestRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:" + uiPort, "http://localhost:3000"},
		AllowCredentials: true,
	})
	restServer := &http.Server{
		Addr:           ":" + restPort,
		Handler:        c.Handler(restRouter),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	uiRouter, err := route.UIRouter()
	if err != nil {
		doIncludeUI = false
	} else {
		uiServer = &http.Server{
			Addr:           ":" + uiPort,
			Handler:        uiRouter,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	}

	var wg sync.WaitGroup

	//start UI server
	if doIncludeUI {
		wg.Add(1)
		go func() {
			log.Fatal(uiServer.ListenAndServe())
			wg.Done()
		}()
		fmt.Printf("Started UI server at port %v ... \n", uiPort)
	}

	//start REST server
	wg.Add(1)
	go func() {
		log.Fatal(restServer.ListenAndServe())
		wg.Done()
	}()
	fmt.Printf("Started REST server at port %v ... \n", restPort)
	if !doIncludeUI {
		fmt.Println(("(UI is not included in this build. If you want to include the UI, build using '-tags=ui')"))
	}

	wg.Wait()
}
