package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func serve(app *application) error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting server on port %d", app.port)

	return server.ListenAndServe()
}

/*
The serve function sets up an HTTP server with specific configurations
like address, handler, and timeouts.
It uses the routes function to get the handler (Gin instance) for the server.
The server is started with ListenAndServe,
if there is an error it will log the error and exit the program.
*/
