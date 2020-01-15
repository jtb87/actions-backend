package app

import (
	"backend/entities"
	"backend/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func (s *Server) authenticate(w http.ResponseWriter, r *http.Request) {
	p := entities.Profile{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		utils.RespondWithError(w, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	err := s.Store.ProfileAuthentication(&p)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "{'error': 'User not authenticated'}")
		log.Errorf("failed authentication request: %v", err)
		return
	}
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{Name: "codiq_session", Value: p.Token, Expires: expiration, Path: "/", HttpOnly: false}
	http.SetCookie(w, &cookie)

	act := map[string]string{"status": "authenticated"}
	utils.RespondWithJSON(w, http.StatusOK, act)
}

func (s *Server) authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("codiq_session")
		if err != nil {
			utils.RespondWithJSON(w, http.StatusUnauthorized, "{'error': 'User not authenticated'}")
			log.Errorf("error: %v", err)
			return
		}
		profile, err := s.Store.AuthorizeToken(c.Value)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusForbidden, "{'error': 'User not authorized'}")
			log.Errorf("error: %v", err)
			return
		}
		ctx := context.WithValue(r.Context(), "profile", profile)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
