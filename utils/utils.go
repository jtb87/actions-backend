package utils

import (
	"encoding/json"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// LogRequest middleware for logging requests
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		}).Info("http-request")
		next.ServeHTTP(w, r)
	})
}

// InitLog set logging to std out or logfile
func InitLog() {
	log.SetFormatter(&log.JSONFormatter{})
	// file, err := os.OpenFile("api_logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Info("Failed to log to file! using default stderr")
	// }
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

// RespondWithError respond with error
func RespondWithError(w http.ResponseWriter, message string) {
	RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": message})
}

// RespondWithJSON respond with json
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
