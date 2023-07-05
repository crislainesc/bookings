package dbrepo

import (
	"database/sql"

	"github.com/crislainesc/bookings/internal/config"
	"github.com/crislainesc/bookings/internal/repository"
)

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewTestRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
