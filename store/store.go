package store

import (
	"backend/entities"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	log "github.com/sirupsen/logrus"
)

// DbStore holds the
type DbStore struct {
	DB     *sqlx.DB
	Statsd string // timing
}

// InitializeStore inits datbase
func InitializeStore(psqlInfo string) (entities.StoreInterface, error) {
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("database connected")
	return &DbStore{db, ""}, nil
}
