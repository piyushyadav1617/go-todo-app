package database

import (
	"fmt"
)

func Migrate() error {
	db := GetDB()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            completed BOOLEAN NOT NULL DEFAULT false
        );
	`)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)

	} else {
		return nil
	}
}
