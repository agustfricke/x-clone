package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() {
    var err error

    dbPath := "db.sqlite3"

    DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

    if err != nil {
        panic("Failed to connect to the database")
    }

    fmt.Println("Connection opened to the database")
    DB.AutoMigrate(&models.User{}, &models.Post{})
    fmt.Println("Database migrated")
}
