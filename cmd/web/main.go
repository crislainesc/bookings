package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/crislainesc/bookings/pkg/config"
	"github.com/crislainesc/bookings/pkg/handlers"
	"github.com/crislainesc/bookings/pkg/render"
)

const portNumber = ":8080"

// main is the main function
func main() {
	tcache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app := config.AppConfig{TemplateCache: tcache, UseCache: false}

	repository := handlers.NewRepository(&app)
	handlers.NewHandlers(repository)

	render.NewTemplates(&app)
	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Printf("Starting application on port %s", portNumber)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
