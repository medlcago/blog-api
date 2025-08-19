package database

import (
	"blog-api/internal/models"
	"log"

	"gorm.io/gorm"
)

func RunMigrations() error {
	if DB == nil {
		return gorm.ErrInvalidDB
	}

	err := DB.AutoMigrate(
		&models.Post{},
	)
	if err != nil {
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil
}
