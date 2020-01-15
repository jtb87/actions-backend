package app

import (
	"backend/entities"
	"backend/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (s *Server) getListOfActions(w http.ResponseWriter, r *http.Request) {
	profile := r.Context().Value("profile").(entities.Profile)
	qryParams := r.URL.Query()
	var categoryID string
	if val, ok := qryParams["category_id"]; ok {
		categoryID = val[0]
	} else {
		utils.RespondWithError(w, "query param 'category_id' required.")
		return
	}
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "query param 'category_id' not an integer")
		return
	}
	actions, err := s.Store.GetListOfActions(profile.ID, id)
	if err != nil {
		log.Warn("something has gone terribly wrong")
		utils.RespondWithError(w, "something has gone wrong")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, actions)
}

func (s *Server) getAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "something went wrong")
		return
	}
	act, err := s.Store.GetActionByID(id)
	if err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, act)
}

// createAction parses the json body and creates a new action in the database
func (s *Server) createAction(w http.ResponseWriter, r *http.Request) {
	var JSONbody struct {
		Subject     string `json:"subject"`
		Description string `json:"description"`
		CategoryID  int    `json:"category_id"`
		ActionDate  string `json:"action_date"`
		PlannedDate string `json:"planned_date,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&JSONbody); err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "Invalid request payload")
		return
	}
	profile := r.Context().Value("profile").(entities.Profile)
	var typeError []error
	action := entities.Action{
		Subject:     JSONbody.Subject,
		Description: JSONbody.Description,
		CategoryID:  JSONbody.CategoryID,
		ActionDate:  TransformToTime(JSONbody.ActionDate, &typeError),
		PlannedDate: TransformToTime(JSONbody.PlannedDate, &typeError),
		ProfileID:   profile.ID,
	}
	if typeError != nil {
		log.Warning(typeError)
		utils.RespondWithError(w, "something went wrong with typechecking")
		return
	}
	err := s.Store.CreateAction(&action)
	if err != nil {
		log.Warn("something has gone terribly wrong")
		utils.RespondWithError(w, "something has gone wrong")
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, action)
}

func (s *Server) deleteAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "something went wrong")
		return
	}
	err = s.Store.DeleteAction(id)
	if err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "something went wrong")
		return
	}
	utils.RespondWithJSON(w, http.StatusAccepted, nil)
}

func (s *Server) updateAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "something went wrong")
		return

	}
	var JSONbody struct {
		Subject     string `json:"subject"`
		Description string `json:"description"`
		CategoryID  int    `json:"category_id"`
		ActionDate  string `json:"action_date"`
		PlannedDate string `json:"planned_date,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&JSONbody); err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "Invalid request payload")
		return
	}
	profile := r.Context().Value("profile").(entities.Profile)
	var typeError []error
	action := entities.Action{
		ID:          id,
		Subject:     JSONbody.Subject,
		Description: JSONbody.Description,
		CategoryID:  JSONbody.CategoryID,
		ActionDate:  TransformToTime(JSONbody.ActionDate, &typeError),
		PlannedDate: TransformToTime(JSONbody.PlannedDate, &typeError),
		ProfileID:   profile.ID,
	}
	if typeError != nil {
		log.Warning(typeError)
		utils.RespondWithError(w, "something went wrong with typechecking")
		return
	}
	err = s.Store.UpdateAction(&action)
	if err != nil {
		log.Warn(err)
		utils.RespondWithError(w, "something has gone wrong")
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, action)
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
