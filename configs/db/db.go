package db

import "database/sql"

func Connect(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("aviapocket", connectionString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
