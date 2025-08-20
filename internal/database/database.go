package database

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	initOnce  sync.Once
	initError error
)
var (
	ErrNotInitialized = fmt.Errorf("database not initialized")
	ErrEmptyDSN       = fmt.Errorf("DSN cannot be empty")

	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type PoolConfig struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func InitDB(dsn string, poolConfig *PoolConfig) error {
	initOnce.Do(func() {
		if dsn == "" {
			initError = ErrEmptyDSN
			return
		}

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		})

		if err != nil {
			initError = err
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			initError = err
			return
		}

		if err := sqlDB.Ping(); err != nil {
			initError = err
			return
		}

		if poolConfig != nil {
			sqlDB.SetMaxIdleConns(poolConfig.MaxIdleConns)
			sqlDB.SetMaxOpenConns(poolConfig.MaxOpenConns)
			sqlDB.SetConnMaxLifetime(poolConfig.ConnMaxLifetime)
		}

		DB = db
		log.Println("✅ Successfully connected to the database")
	})

	return initError
}

func GetDb() *gorm.DB {
	if DB == nil {
		panic("DB not initialized. Call InitDB() first")
	}
	return DB
}

func BuildDSN(DbHost, DbUser, DbPassword, DbName, DbPort string) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		DbHost,
		DbUser,
		DbPassword,
		DbName,
		DbPort,
	)
}
