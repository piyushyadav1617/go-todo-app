package database

import (
	"database/sql"
	"fmt"
	"sync"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

func InitDB(connectionString string) error {
	var initErr error
	dbOnce.Do(func() {
		var err error
		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			initErr = fmt.Errorf("failed to open database: %w", err)
			return
		}
		if err = db.Ping(); err != nil {
			initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}
	})
	return initErr
}

func GetDB() *sql.DB {
	if db == nil {
		panic("database not initialized - call InitDB first")
	}
	return db
}
