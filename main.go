package main

import (
	"backend/entities"
	"backend/store"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jtb87/goconfig"
	log "github.com/sirupsen/logrus"
)

// DateFormat holds the date format to be used
const DateFormat = "2006-01-02"

// DateTimeFormat is the date-time format to use
const DateTimeFormat = "2006-01-02 15:04:05"

var storeInterface entities.StoreInterface

func main() {
	var c Config
	err := goconfig.ParseConfig("config.json", &c)
	if err != nil {
		log.Fatal(err)
	}
	var psqlInfo string
	if c.AppEngine {
		psqlInfo = fmt.Sprintf("postgres://%s:%s@/postgres?host=/cloudsql/%s", c.DatabaseUsername, c.DatabasePassword, c.DatabaseHost)
	} else {
		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			c.DatabaseHost, c.DataBasePort, c.DatabaseUsername, c.DatabasePassword, c.DatabaseName)
	}

	db, err := store.InitializeStore(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	app := App{
		Config: c,
		Store:  db,
	}

	// initialize router
	app.NewRouter()
	// initialize log
	InitLog()
	log.Infof("Server running on http://localhost:%s", app.Config.Port)
	app.startServer()
}

// App main struct
type App struct {
	Router *mux.Router
	Config Config
	Store  entities.StoreInterface
}

// Config holds all the initialization information
type Config struct {
	Port             string        `json:"port"`
	Timeout          time.Duration `json:"timeout"`
	DatabaseHost     string        `json:"db_host"`
	DataBasePort     int           `json:"db_port"`
	DatabaseName     string        `json:"db_name"`
	DatabaseUsername string        `json:"db_username"`
	DatabasePassword string        `json:"db_password"`
	AppEngine        bool          `json:"appengine"`
}

// startserver
func (a *App) startServer() {
	allowedHeaders := handlers.AllowedHeaders([]string{"content-type"})
	allowedOrigins := handlers.AllowedOrigins([]string{"https://www.codiq.eu", "https://codiq.eu", "http://localhost:8080"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	timeout := time.Second * a.Config.Timeout
	port := os.Getenv("PORT")
	if port == "" {
		port = a.Config.Port
		log.Printf("Defaulting to port %s", port)
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: http.TimeoutHandler(handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods, handlers.AllowCredentials())(a.Router), timeout, "timeout"),
	}
	log.Fatal(s.ListenAndServe())
}

func respondWithError(w http.ResponseWriter, message string) {
	respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

// NewRouter creates a new muxrouter
func (a *App) NewRouter() {
	a.Router = mux.NewRouter().StrictSlash(true)
	// initialize routes
	a.initializeAuth()
	a.initializeAPI()
	// initialize global middleware
	a.Router.Use(LogRequest)
}

// InitializeRoutes Initialize routes
func (a *App) initializeAuth() {
	a.Router.HandleFunc("/auth/login", a.authenticate).Methods("POST")
}

func (a *App) initializeAPI() {
	api := a.Router.PathPrefix("/api").Subrouter()
	api.Use(a.authorizationMiddleware)

	api.HandleFunc("/action", a.getListOfActions).Methods("GET")
	api.HandleFunc("/action", a.createAction).Methods("POST")
	api.HandleFunc("/action/{id}", a.getAction).Methods("GET")
	api.HandleFunc("/action/{id}", a.deleteAction).Methods("DELETE")
	api.HandleFunc("/action/{id}/update", a.updateAction).Methods("POST")
}
