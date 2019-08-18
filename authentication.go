package main

import (
	"backend/entities"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func (a *App) authenticate(w http.ResponseWriter, r *http.Request) {
	p := entities.Profile{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	fmt.Printf("%+v", p)
	err := a.Store.ProfileAuthentication(&p)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{Name: "codiq_session", Value: p.Token, Expires: expiration, Path: "/", HttpOnly: false}
	http.SetCookie(w, &cookie)

	act := map[string]string{"status": "authenticated"}
	respondWithJSON(w, http.StatusOK, act)
}

func (a *App) authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from header
		c, err := r.Cookie("codiq_session")
		if err != nil {
			respondWithError(w, err.Error())
			log.Errorf("error: %v", err)
			return
		}
		profile, err := a.Store.AuthorizeToken(c.Value)
		if err != nil {
			respondWithJSON(w, http.StatusUnauthorized, "{'error': 'User not Authenticated'}")
			return
		}
		log.Infof("Authorization middleware succesfull for %s", profile.Username)
		// set profile as headers in context
		ctx := context.WithValue(r.Context(), "profile", profile)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
