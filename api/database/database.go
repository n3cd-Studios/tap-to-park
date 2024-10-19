package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Connect() {
	database, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	err = database.AutoMigrate(&Organization{}, &Invite{}, &Spot{}, &Reservation{}, &User{}, &Session{})
	if err != nil {
		panic("Failed to automigrate models.")
	}

	Db = database

}
