package repository

import (
	"time"

	"github.com/crislainesc/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(reservation models.Reservation) (int, error)
	InsertRoomRestriction(restriction models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(roomID int) (models.Room, error)
}
