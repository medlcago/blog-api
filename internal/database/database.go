package database

import (
	"blog-api/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type DB struct {
	db *gorm.DB
}

func New(cfg config.DatabaseConfig) (*DB, error) {
	db, err := gorm.Open(postgres.Open(BuildDSN(cfg)), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return &DB{db: db}, nil
}

func (d *DB) Get() *gorm.DB {
	return d.db
}

func (d *DB) Close() error {
	sqlDB, err := d.Get().DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func BuildDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbName,
		cfg.DbPort,
	)
}
