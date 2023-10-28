package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/pkg/errors"
)

func NewClient(url string) (*gorm.DB, func(), error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, func() {}, errors.Wrap(err, "gorm open postgresql connection")
	}
	postgresDB, err := db.DB()
	if err != nil {
		return nil, func() {}, errors.Wrap(err, "getting sql.DB entity")
	}

	if err = postgresDB.Ping(); err != nil {
		return nil, func() {}, errors.Wrap(err, "postgres db ping")
	}

	return db, func() {
		postgresDB.Close()
	}, nil
}
