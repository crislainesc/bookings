package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/crislainesc/bookings/internal/config"
	"github.com/crislainesc/bookings/internal/forms"
	"github.com/crislainesc/bookings/internal/models"
	"github.com/crislainesc/bookings/internal/render"
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
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.tmpl.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.tmpl.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
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
