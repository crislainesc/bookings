package dbrepo

import (
	"errors"
	"time"

	"github.com/crislainesc/bookings/internal/models"
)

func (repository *testDBRepo) AllUsers() bool {
	return true
}

func (repository *testDBRepo) InsertReservation(reservation models.Reservation) (int, error) {
	if reservation.RoomID == 2 {
		return 0, errors.New("can't insert reservation")
	}

	return 1, nil
}

func (repository *testDBRepo) InsertRoomRestriction(restriction models.RoomRestriction) error {
	if restriction.RoomID == 1000 {
		return errors.New("can't insert room restriction")
	}

	return nil
}

func (repository *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

func (repository *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	return rooms, nil
}

func (repository *testDBRepo) GetRoomByID(roomID int) (models.Room, error) {
	var room models.Room

	if roomID > 2 {
		return room, errors.New("not found")
	}

	return room, nil
}
