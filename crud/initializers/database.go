package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	// postgres://eamqxkre:0cv5hH7oY3TMPRjWB_gIWxppNpVjq2Ki@bubble.db.elephantsql.com/eamqxkre
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB")
	}
}
