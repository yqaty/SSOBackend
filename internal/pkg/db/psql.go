package db

import (
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/hints"
)

var (
	ErrorNilOption = errors.New("option is nil")
)

type Option struct {
	DSN          string
	MaxOpenConn  int
	MaxIdleConn  int
	MaxLifetime  time.Duration
	PreStatement bool
}

func NewPostgreConn(opt *Option) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(opt.DSN), &gorm.Config{PrepareStmt: opt.PreStatement})
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	// limit connection limit
	sqlDb.SetMaxOpenConns(opt.MaxOpenConn)
	sqlDb.SetMaxIdleConns(opt.MaxIdleConn)
	sqlDb.SetConnMaxIdleTime(opt.MaxLifetime)

	// !PostgresQL not support hints system
	db = db.Clauses(hints.New(" MAX_EXECUTION_TIME(1000) "))

	return db, nil
}
