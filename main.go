package main

import (
	"backend/app"
	"backend/entities"
	"backend/store"
	"backend/utils"
	"fmt"

	"github.com/jtb87/goconfig"
	log "github.com/sirupsen/logrus"
)

var storeInterface entities.StoreInterface

func main() {
	var c app.Config
	err := goconfig.ParseConfig("config_prod_local.json", &c)
	if err != nil {
		log.Fatal(err)
	}
	var psqlInfo string
	if c.ServerEngine {
		psqlInfo = fmt.Sprintf("%s:%s@cloudsql(%s)/", c.DatabaseUsername, c.DatabasePassword, c.DatabaseHost)
		// psqlInfo = fmt.Sprintf("postgres://%s:%s@/postgres?host=/cloudsql/%s", c.DatabaseUsername, c.DatabasePassword, c.DatabaseHost)
	} else {
		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			c.DatabaseHost, c.DataBasePort, c.DatabaseUsername, c.DatabasePassword, c.DatabaseName)
	}

	db, err := store.InitializeStore(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	server := app.Server{
		Config: c,
		Store:  db,
	}

	// initialize router
	server.NewRouter()
	// initialize log
	utils.InitLog()
	log.Infof("Server running on http://localhost:%s", server.Config.Port)
	server.StartServer()
}
