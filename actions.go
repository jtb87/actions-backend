package main

import (
	"backend/entities"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a *App) getListOfActions(w http.ResponseWriter, r *http.Request) {
	profile := r.Context().Value("profile").(entities.Profile)
	actions, err := a.Store.GetListOfActions(profile.ID)
	if err != nil {
		log.Warn("something has gone terribly wrong")
		respondWithError(w, "something has gone wrong")
		return
	}
	respondWithJSON(w, http.StatusOK, actions)
}

func (a *App) getAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Warning(err)
		respondWithError(w, "something went wrong")
		return
	}
	act, err := a.Store.GetActionByID(id)
	if err != nil {
		log.Warning(err)
		respondWithError(w, "something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, act)
}

// createAction parses the json body and creates a new action in the database
func (a *App) createAction(w http.ResponseWriter, r *http.Request) {
	var JSONbody struct {
		Subject     string `json:"subject"`
		Description string `json:"description"`
		Category    string `json:"category"`
		ActionDate  string `json:"action_date"`
		PlannedDate string `json:"planned_date,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&JSONbody); err != nil {
		log.Warning(err)
		respondWithError(w, "Invalid request payload")
		return
	}
	profile := r.Context().Value("profile").(entities.Profile)
	var typeError []error
	action := entities.Action{
		Subject:     JSONbody.Subject,
		Description: JSONbody.Description,
		Category:    JSONbody.Category,
		ActionDate:  TransformToTime(JSONbody.ActionDate, &typeError),
		PlannedDate: TransformToTime(JSONbody.PlannedDate, &typeError),
		ProfileID:   profile.ID,
	}
	if typeError != nil {
		log.Warning(typeError)
		respondWithError(w, "something went wrong with typechecking")
		return
	}
	err := a.Store.CreateAction(&action)
	if err != nil {
		log.Warn("something has gone terribly wrong")
		respondWithError(w, "something has gone wrong")
		return
	}
	respondWithJSON(w, http.StatusCreated, action)
}

func (a *App) deleteAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Warning(err)
		respondWithError(w, "something went wrong")
		return
	}
	err = a.Store.DeleteAction(id)
	if err != nil {
		log.Warning(err)
		respondWithError(w, "something went wrong")
		return
	}
	respondWithJSON(w, http.StatusAccepted, nil)
}

func (a *App) updateAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Warning(err)
		respondWithError(w, "something went wrong")
		return

	}
	var JSONbody struct {
		Subject     string `json:"subject"`
		Description string `json:"description"`
		Category    string `json:"category"`
		ActionDate  string `json:"action_date"`
		PlannedDate string `json:"planned_date,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&JSONbody); err != nil {
		log.Warning(err)
		respondWithError(w, "Invalid request payload")
		return
	}
	profile := r.Context().Value("profile").(entities.Profile)
	var typeError []error
	action := entities.Action{
		ID:          id,
		Subject:     JSONbody.Subject,
		Description: JSONbody.Description,
		Category:    JSONbody.Category,
		ActionDate:  TransformToTime(JSONbody.ActionDate, &typeError),
		PlannedDate: TransformToTime(JSONbody.PlannedDate, &typeError),
		ProfileID:   profile.ID,
	}
	if typeError != nil {
		log.Warning(typeError)
		respondWithError(w, "something went wrong with typechecking")
		return
	}
	err = a.Store.UpdateAction(&action)
	if err != nil {
		log.Warn(err)
		respondWithError(w, "something has gone wrong")
		return
	}
	respondWithJSON(w, http.StatusCreated, action)
}

// TransformToTime takes a string and a format and transforms it to a *time.Time
func TransformToTime(datetimeRaw string, typeError *[]error) *time.Time {
	if datetimeRaw == "" {
		return nil
	}
	prsdTime, err := time.Parse(DateFormat, datetimeRaw)
	if err != nil {
		log.Println(err)
		*typeError = append(*typeError, err)
	}
	return &prsdTime
}

// TransformStringToInt takes a string and a format and transforms it to a *time.Time
func TransformStringToInt(intString string, typeError *[]error) int {
	intResult, err := strconv.Atoi(intString)
	if err != nil {
		log.Warning(err)
		*typeError = append(*typeError, err)
	}
	return intResult
}
