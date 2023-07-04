package models

import "time"

type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
