package app

import (
	"backend/entities"
	"backend/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// DateFormat holds the date format to be used
const DateFormat = "2006-01-02"

// DateTimeFormat is the date-time format to use
const DateTimeFormat = "2006-01-02 15:04:05"

// Server main struct
type Server struct {
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
	ServerEngine     bool          `json:"Serverengine"`
}

// NewRouter creates a new muxrouter
func (s *Server) NewRouter() {
	s.Router = mux.NewRouter().StrictSlash(true)
	// initialize routes
	s.initializeAuth()
	s.initializeAPI()
	// initialize global middleware
	s.Router.Use(utils.LogRequest)
}

// InitializeRoutes Initialize routes
func (s *Server) initializeAuth() {
	s.Router.HandleFunc("/auth/login", s.authenticate).Methods("POST")
}

func (s *Server) initializeAPI() {
	api := s.Router.PathPrefix("/v1").Subrouter()
	api.Use(s.authorizationMiddleware)

	api.HandleFunc("/action_list", s.getListOfActions).Methods("GET")
	api.HandleFunc("/action", s.createAction).Methods("POST")
	api.HandleFunc("/action/{id}", s.getAction).Methods("GET")
	api.HandleFunc("/action/{id}", s.deleteAction).Methods("DELETE")
	api.HandleFunc("/action/{id}/update", s.updateAction).Methods("POST")

	api.HandleFunc("/category", s.getCategory).Methods("GET")
	api.HandleFunc("/category/list", s.getListOfCategories).Methods("GET")
	// api.HandleFunc("/category/update", s.updateCategory).Methods("POST")
	api.HandleFunc("/category/create", s.createCategory).Methods("POST")
	api.HandleFunc("/category/delete", s.deleteCategory).Methods("POST")
}

// StartServer starts the server
func (s *Server) StartServer() {
	allowedHeaders := handlers.AllowedHeaders([]string{"content-type"})
	allowedOrigins := handlers.AllowedOrigins([]string{"https://www.codiq.eu", "https://codiq.eu", "http://localhost:8080"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	timeout := time.Second * s.Config.Timeout
	port := os.Getenv("PORT")
	if port == "" {
		port = s.Config.Port
		log.Printf("Defaulting to port %s", port)
	}

	serv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: http.TimeoutHandler(handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods, handlers.AllowCredentials())(s.Router), timeout, "timeout"),
	}
	log.Fatal(serv.ListenAndServe())
}
