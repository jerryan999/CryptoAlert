package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Initialize initialises the database
func Initialize(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil || db == nil {
		panic("Error connecting to database")
	}

	return db
}

// Migrate migrates the database
func Migrate(db *sql.DB) {
	sql := `
            CREATE TABLE IF NOT EXISTS alerts(
                    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					crypto TEXT NOT NULL,
                    price FLOAT NOT NULL,
                    direction BOOLEAN NOT NULL
            );
       `

	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
}
