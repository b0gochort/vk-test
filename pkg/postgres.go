package pkg

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgres() (*sql.DB, error) {
	connStr := "postgresql://postgres:king1337@localhost:5432/vkbot?sslmode=disable"
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
