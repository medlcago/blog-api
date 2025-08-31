package database

import (
	"blog-api/internal/models"
)

func (d *DB) RunMigrations() error {
	db := d.Get()
	return db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.PostEntity{},
		&models.Reaction{},
	)
}
