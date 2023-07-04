package dbrepo

import (
	"context"
	"time"

	"github.com/crislainesc/bookings/internal/models"
)

func (repository *postgresDBRepo) AllUsers() bool {
	return true
}

func (repository *postgresDBRepo) InsertReservation(reservation models.Reservation) (int, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		INSERT INTO 
			reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id
	`

	var newID int

	err := repository.DB.QueryRowContext(context, query,
		reservation.FirstName,
		reservation.LastName,
		reservation.Email,
		reservation.Phone,
		reservation.StartDate,
		reservation.EndDate,
		reservation.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (repository *postgresDBRepo) InsertRoomRestriction(restriction models.RoomRestriction) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		INSERT INTO 
			room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := repository.DB.ExecContext(context, query,
		restriction.StartDate,
		restriction.EndDate,
		restriction.RoomID,
		restriction.ReservationID,
		time.Now(),
		time.Now(),
		restriction.RestrictionID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repository *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		SELECT 
			count(id)
		FROM
		  room_restrictions
		WHERE
			room_id = $1 AND
		  $2 < end_date AND $3 > start_date
	`

	var numRows int

	row := repository.DB.QueryRowContext(context, query, roomID, start, end)
	err := row.Scan(&numRows)

	if err != nil {
		return false, err
	}

	return numRows == 0, nil
}
