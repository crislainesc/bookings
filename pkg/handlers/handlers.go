package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crislainesc/bookings/pkg/config"
	"github.com/crislainesc/bookings/pkg/models"
	"github.com/crislainesc/bookings/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

type JsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
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

	render.RenderTemplate(w, r, "home.page.tmpl.html", &models.TemplateData{})
}

// About is the handler for the about page
func (repository *Repository) About(w http.ResponseWriter, r *http.Request) {
	remoteIP := repository.App.Session.GetString(r.Context(), "remote_ip")
	stringMap := map[string]string{
		"test":      "Hello world",
		"remote_ip": remoteIP,
	}

	render.RenderTemplate(w, r, "about.page.tmpl.html", &models.TemplateData{StringMap: stringMap})
}

// Reservation is the handler for the reservation page
func (repository *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl.html", &models.TemplateData{})
}

// Generals is the handler for the generals page
func (repository *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl.html", &models.TemplateData{})
}

// Majors is the handler for the majors page
func (repository *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl.html", &models.TemplateData{})
}

// Availability is the handler for the search availability page
func (repository *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl.html", &models.TemplateData{})
}

// PostAvailability is the handler for the search availability
func (repository *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	result := fmt.Sprintf("start date is %s and end date is %s", start, end)

	w.Write([]byte(result))
}

// AvailabilityJSON is the handler for the search availability
func (repository *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := JsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact is the handler for the contact page
func (repository *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl.html", &models.TemplateData{})
}
