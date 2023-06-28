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
	app := config.AppConfig{}

	tcache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tcache
	app.UseCache = false

	repository := handlers.NewRepository(&app)
	handlers.NewHandlers(repository)

	render.NewTemplates(&app)

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Staring application on port %s", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}
