package dbrepo

import (
	"context"
	"time"

	"github.com/crislainesc/bookings/internal/models"
)

func (repository *postgresDBRepo) AllUsers() bool {
	return true
}

func (repository *postgresDBRepo) InsertReservation(res models.Reservation) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `INSERT INTO 	reservations (first_name, last_name, email, phone, start_date, end_date,
		room_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := repository.DB.ExecContext(context, query,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
