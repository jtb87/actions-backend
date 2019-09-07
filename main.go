package main

import (
	"backend/entities"
	"backend/store"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// DateFormat holds the date format to be used
const DateFormat = "2006-01-02"

// DateTimeFormat is the date-time format to use
const DateTimeFormat = "2006-01-02 15:04:05"

var storeInterface entities.StoreInterface

func main() {
	db, err := store.InitializeStore()
	if err != nil {
		log.Fatal(err)
	}
	app := App{
		Config: Config{
			port:    ":9090",
			timeout: time.Second * 15,
		},
		Store: db,
	}

	// initialize router
	app.NewRouter()
	// initialize log
	initLog()
	log.Infof("Server running on http://localhost%s", app.Config.port)
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
	port         string
	timeout      time.Duration
	databaseHost string
	databasPort  string
}

// startserver
func (a *App) startServer() {
	allowedHeaders := handlers.AllowedHeaders([]string{"content-type"})
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:8080"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	s := &http.Server{
		Addr:    a.Config.port,
		Handler: http.TimeoutHandler(handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods, handlers.AllowCredentials())(a.Router), a.Config.timeout, "timeout"),
	}
	log.Fatal(s.ListenAndServe())
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
	a.Router.Use(logRequest)
}

// logRequest logging middleware
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		}).Info("http-request")
		next.ServeHTTP(w, r)
	})
}

// func authMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var header = r.Header.Get("x-access-token") //Grab the token from the header
// 		header = strings.TrimSpace(header)
// 		if header == "" {
// 			//Token is missing, returns with error code 403 Unauthorized
// 			w.WriteHeader(http.StatusForbidden)
// 			json.NewEncoder(w).Encode(Exception{Message: "Missing auth token"})
// 			return
// 		}
// 		// tk := &models.Token{}

// 		_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
// 			return []byte("secret"), nil
// 		})

// 		if err != nil {
// 			w.WriteHeader(http.StatusForbidden)
// 			json.NewEncoder(w).Encode(Exception{Message: err.Error()})
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), "user", tk)
// 		next.ServeHTTP(w, r.WithContext(ctx))

// 	})
// }

// Set logging to std out or logfile
func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	// file, err := os.OpenFile("api_logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Info("Failed to log to file! using default stderr")
	// }
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}
