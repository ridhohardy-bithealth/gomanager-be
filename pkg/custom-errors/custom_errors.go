package customErrors

import (
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

var (
	ErrNotFound     = pgx.ErrNoRows
	ErrConflict     = errors.New("conflict")
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
)
