package dbrepo

import (
	"context"
	"time"

	"github.com/crislainesc/bookings/internal/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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

func (repository *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		SELECT 
			r.id, r.room_name
		FROM
		  rooms r
		WHERE r.id NOT IN
			(SELECT 
				room_id
			FROM
				room_restrictions rr
			WHERE
				$1 < rr.end_date AND $2 > rr.start_date)
	`

	var rooms []models.Room

	rows, err := repository.DB.QueryContext(context, query, start, end)

	if err != nil {
		return rooms, err
	}

	if rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

func (repository *postgresDBRepo) GetRoomByID(roomID int) (models.Room, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		SELECT id, room_name, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`

	var room models.Room

	row := repository.DB.QueryRowContext(context, query, roomID)

	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}

func (repository *postgresDBRepo) GetUserByID(userID int) (models.User, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		SELECT id, first_name, last_name, email, password, access_level, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User

	row := repository.DB.QueryRowContext(context, query, userID)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (repository *postgresDBRepo) UpdateUser(user models.User) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
	`

	_, err := repository.DB.ExecContext(context, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (repository *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var id int
	var hashedPassword string

	query := `
		SELECT id, password
		FROM users
		WHERE email = $1
	`

	row := repository.DB.QueryRowContext(context, query, email)
	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func getReservations(query string, repository *postgresDBRepo) ([]models.Reservation, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var reservations []models.Reservation

	rows, err := repository.DB.QueryContext(context, query)

	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(
			&r.ID,
			&r.FirstName,
			&r.LastName,
			&r.Email,
			&r.Phone,
			&r.StartDate,
			&r.EndDate,
			&r.RoomID,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.Processed,
			&r.Room.ID,
			&r.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, r)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil

}

func (repository *postgresDBRepo) GetAllReservations() ([]models.Reservation, error) {
	query := `
		SELECT 
			r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
			r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
			rm.id, rm.room_name
		FROM reservations r
		LEFT JOIN rooms rm ON (r.room_id = rm.id)
		ORDER BY r.start_date asc
	`

	reservations, err := getReservations(query, repository)
	if err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (repository *postgresDBRepo) GetAllNewReservations() ([]models.Reservation, error) {
	query := `
		SELECT 
			r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
			r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
			rm.id, rm.room_name
		FROM reservations r
		LEFT JOIN rooms rm ON (r.room_id = rm.id)
		WHERE processed = 0
		ORDER BY r.start_date asc
	`

	reservations, err := getReservations(query, repository)

	if err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (repository *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `
			SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id,
			r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
			FROM reservations r 
			LEFT JOIN rooms rm on (r.room_id = rm.id)
			WHERE r.id = $1`

	row := repository.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)

	if err != nil {
		return res, err
	}

	return res, err
}

func (repository *postgresDBRepo) UpdateReservation(reservation models.Reservation) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		UPDATE reservations
		SET first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := repository.DB.ExecContext(context, query,
		reservation.FirstName,
		reservation.LastName,
		reservation.Email,
		reservation.Phone,
		time.Now(),
		reservation.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repository *postgresDBRepo) DeleteReservation(id int) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		DELETE FROM reservations
		WHERE id = $1
	`

	_, err := repository.DB.ExecContext(context, query, id)

	if err != nil {
		return err
	}

	return nil
}

func (repository *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		UPDATE reservations
		SET processed = $1
		WHERE id = $2
	`

	_, err := repository.DB.ExecContext(context, query,
		processed,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
