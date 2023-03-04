package util

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// GetDB returns a database connection
func GetDB() *pgx.Conn {
	db, err := pgx.Connect(context.Background(), "postgresql://postgres:postgres@localhost")

	if err != nil {
		panic(err)
	}
	return db
}
