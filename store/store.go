package store

import (
	"backend/entities"
	_ "database/sql"
	"fmt"

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
func InitializeStore(user string, password string, port int, dbname string, host string) (entities.StoreInterface, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", psqlInfo)
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
