package handlers

import (
	"net/http"

	"github.com/crislainesc/bookings/pkg/config"
	"github.com/crislainesc/bookings/pkg/models"
	"github.com/crislainesc/bookings/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepository(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
	}
}

func NewHandlers(repo *Repository) {
	Repo = repo
}

// Home is the handler for the home page
func (repository *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	repository.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl.html", &models.TemplateData{})
}

// About is the handler for the about page
func (repository *Repository) About(w http.ResponseWriter, r *http.Request) {
	remoteIP := repository.App.Session.GetString(r.Context(), "remote_ip")
	stringMap := map[string]string{
		"test":      "Hello world",
		"remote_ip": remoteIP,
	}

	render.RenderTemplate(w, "about.page.tmpl.html", &models.TemplateData{StringMap: stringMap})
}
