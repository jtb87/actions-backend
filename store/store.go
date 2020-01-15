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
func InitializeStore(psqlInfo string) (entities.StoreInterface, error) {
	// db, err := sqlx.Connect("postgres", psqlInfo)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("database connected")
	return &DbStore{db, ""}, nil
}
