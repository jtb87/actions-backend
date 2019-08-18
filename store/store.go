package store

import (
	"backend/entities"
	_ "database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// DbStore holds the
type DbStore struct {
	DB     *sqlx.DB
	Statsd string // timing
}

// InitializeStore inits datbase
func InitializeStore() (entities.StoreInterface, error) {
	// this Pings the database trying to connect, panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("postgres", "user=admin password=solarrules sslmode=disable")
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	log.Infoln("database connected")
	return &DbStore{db, ""}, nil
}
