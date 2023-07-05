package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/crislainesc/bookings/internal/config"
	"github.com/crislainesc/bookings/internal/driver"
	"github.com/crislainesc/bookings/internal/forms"
	"github.com/crislainesc/bookings/internal/helpers"
	"github.com/crislainesc/bookings/internal/models"
	"github.com/crislainesc/bookings/internal/render"
	"github.com/crislainesc/bookings/internal/repository"
	"github.com/crislainesc/bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

var (
	Repo       *Repository
	dateLayout = "2006-01-02"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

type JsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func NewRepository(app *config.AppConfig, db *driver.Database) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresRepo(db.SQL, app),
	}
}

func NewHandlers(repo *Repository) {
	Repo = repo
}

// Home is the handler for the home page
func (repository *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl.html", &models.TemplateData{})
}

// About is the handler for the about page
func (repository *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl.html", &models.TemplateData{})
}

// Reservation is the handler for the reservation page
func (repository *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repository.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		repository.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := repository.DB.GetRoomByID(reservation.RoomID)

	if err != nil {
		repository.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.Room.RoomName = room.RoomName

	repository.App.Session.Put(r.Context(), "reservation", reservation)

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format(dateLayout)
	ed := reservation.EndDate.Format(dateLayout)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "make-reservation.page.tmpl.html", &models.TemplateData{
		Data:      data,
		Form:      forms.New(nil),
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (repository *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repository.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		repository.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {

		repository.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	startDate, err := time.Parse(dateLayout, sd)
	if err != nil {
		repository.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(dateLayout, ed)
	if err != nil {
		repository.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))

	if err != nil {
		repository.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")
	reservation.StartDate = startDate
	reservation.EndDate = endDate
	reservation.RoomID = roomID

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "Invalid data", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.tmpl.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := repository.DB.InsertReservation(reservation)
	if err != nil {
		repository.App.Session.Put(r.Context(), "error", "can't create new reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	repository.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = repository.DB.InsertRoomRestriction(restriction)
	if err != nil {
		repository.App.Session.Put(r.Context(), "error", "can't finish reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	repository.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals is the handler for the generals page
func (repository *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl.html", &models.TemplateData{})
}

// Majors is the handler for the majors page
func (repository *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl.html", &models.TemplateData{})
}

// Availability is the handler for the search availability page
func (repository *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl.html", &models.TemplateData{})
}

// PostAvailability is the handler for the search availability
func (repository *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	startDate, err := time.Parse(dateLayout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(dateLayout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}
	rooms, err := repository.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	for _, room := range rooms {
		repository.App.InfoLog.Println("ROOM:", room.ID, room.RoomName)
	}

	if len(rooms) == 0 {
		repository.App.Session.Put(r.Context(), "error", "No available rooms")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	repository.App.Session.Put(r.Context(), "reservation", reservation)

	render.Template(w, r, "choose-room.page.tmpl.html", &models.TemplateData{
		Data: data,
	})
}

// AvailabilityJSON is the handler for the search availability
func (repository *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	startDate, _ := time.Parse(dateLayout, sd)
	endDate, _ := time.Parse(dateLayout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := repository.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)

	resp := JsonResponse{
		OK:        available,
		Message:   "",
		RoomID:    strconv.Itoa(roomID),
		StartDate: sd,
		EndDate:   ed,
	}

	out, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		helpers.ServerError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact is the handler for the contact page
func (repository *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl.html", &models.TemplateData{})
}

func (repository *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repository.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		repository.App.ErrorLog.Println("cannot get item from session")
		repository.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	repository.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format(dateLayout)
	ed := reservation.EndDate.Format(dateLayout)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (repository *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation, ok := repository.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		helpers.ServerError(w, errors.New("type assertion to string failed"))
		return
	}

	reservation.RoomID = roomID
	repository.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (repository *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	startDate, _ := time.Parse(dateLayout, sd)
	endDate, _ := time.Parse(dateLayout, ed)

	room, err := repository.DB.GetRoomByID(roomID)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var reservation models.Reservation

	reservation.Room.RoomName = room.RoomName
	reservation.RoomID = roomID
	reservation.StartDate = startDate
	reservation.EndDate = endDate

	repository.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
