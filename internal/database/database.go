package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) error {
	if DB != nil {
		return nil
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return err
	}

	log.Println("✅ Successfully connected to the database")
	return nil
}

func GetDb() *gorm.DB {
	return DB
}
