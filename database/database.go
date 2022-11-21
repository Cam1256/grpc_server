package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	Database *gorm.DB
}

// Init SQLite database using Gorm
func InitDatabase() *DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	log.Println("DB connected")
	return &DB{db}
}
