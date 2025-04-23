package database

import (
	"fmt"
)

func Migrate() error {
	db := GetDB()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(255) NOT NULL,
            password VARCHAR(255) NOT NULL
        );
	`)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            completed BOOLEAN NOT NULL DEFAULT false,
			user_id Int NOT NULL REFERENCES users(id)
        );
	`)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)

	} else {
		return nil
	}
}
