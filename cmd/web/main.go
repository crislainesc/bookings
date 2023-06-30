package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/crislainesc/bookings/internal/config"
	"github.com/crislainesc/bookings/internal/handlers"
	"github.com/crislainesc/bookings/internal/render"
)

const (
	portNumber = ":8080"
)

var (
	app     config.AppConfig
	session *scs.SessionManager
)

// main is the main function
func main() {

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tcache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tcache
	app.UseCache = false

	repository := handlers.NewRepository(&app)
	handlers.NewHandlers(repository)

	render.NewTemplates(&app)
	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Printf("Starting application on port %s\n", portNumber)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
