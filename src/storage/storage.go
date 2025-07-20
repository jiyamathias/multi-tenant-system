// Package storage house all storage/database implementations that performs CRUD operations on our schema
package storage

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"codematic/pkg/environment"
	"codematic/pkg/gorm_sqlmock"
	"codematic/pkg/helper"
)

const packageName = "storage"

// Storage object
type Storage struct {
	Logger zerolog.Logger
	Env    *environment.Env
	DB     *gorm.DB
}

// New Storage, however should panic if it can't be pinged. System should be able to connect to the database
func New(z zerolog.Logger, env *environment.Env) *Storage {
	l := z.With().Str(helper.LogStrKeyModule, packageName).Logger()

	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s TimeZone=Africa/Lagos",
		env.Get("PG_ADDRESS"),
		env.Get("PG_PORT"),
		env.Get("PG_USER"),
		env.Get("PG_DATABASE"),
		env.Get("PG_PASSWORD"),
	)

	db, err := gorm.Open(
		postgres.Open(connString),
		&gorm.Config{},
	)
	if err != nil {
		l.Fatal().Err(err)
		panic(err)
	}

	fmt.Println("Database connected successfully")
	return &Storage{
		Logger: l,
		Env:    env,
		DB:     db,
	}
}

// GetStorage helper for tests/mock
// I expect our storage tests to use this helper going forward.
func GetStorage(t *testing.T) (sqlmock.Sqlmock, *Storage) {
	var (
		mock sqlmock.Sqlmock
		db   *gorm.DB
		err  error
	)

	db, mock, err = gorm_sqlmock.New(gorm_sqlmock.Config{
		Config:     &gorm.Config{},
		DriverName: "postgres",
		DSN:        "mock",
	})

	require.NoError(t, err)

	return mock, NewFromDB(db)
}

// NewFromDB created a new storage with just the database reference passed in
func NewFromDB(db *gorm.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

// Close securely closes the connection to the storage/database
func (s *Storage) Close() {
	sqlDD, _ := s.DB.DB()
	_ = sqlDD.Close()
}
