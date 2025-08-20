package database

import (
	"blog-api/internal/models"
	"log"
)

func RunMigrations() error {
	if DB == nil {
		return ErrNotInitialized
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.Post{},
	)
	if err != nil {
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil
}
