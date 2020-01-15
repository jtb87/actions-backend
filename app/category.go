package app

import (
	"backend/entities"
	"backend/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Data is standard output for lists
type Data struct {
	Data interface{} `json:"data"`
}

func (s *Server) getListOfCategories(w http.ResponseWriter, r *http.Request) {
	profile := r.Context().Value("profile").(entities.Profile)
	categories, err := s.Store.GetListOfCategories(profile.ID)
	if err != nil {
		log.Warn("something has gone terribly wrong")
		utils.RespondWithError(w, "something has gone wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, Data{Data: categories})
}

func (s *Server) getCategory(w http.ResponseWriter, r *http.Request) {
	qryParams := r.URL.Query()
	var categoryID string
	if val, ok := qryParams["id"]; ok {
		categoryID = val[0]
	} else {
		utils.RespondWithError(w, "query param 'id' required.")
		return
	}

	id, err := strconv.Atoi(categoryID)
	if err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "query param 'id' not an integer")
		return
	}
	cat, err := s.Store.GetCategory(id)
	if err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "something went wrong")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, cat)
}

// createCategory creates a new category
func (s *Server) createCategory(w http.ResponseWriter, r *http.Request) {
	var JSONbody struct {
		Name         string  `json:"name"`
		Interval     *int    `json:"interval"`
		IntervalType *string `json:"interval_type"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&JSONbody); err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "Invalid request payload")
		return
	}
	cat := entities.Category{
		Name:         JSONbody.Name,
		Interval:     JSONbody.Interval,
		IntervalType: JSONbody.IntervalType,
	}
	err := cat.Validate()
	if err != nil {
		log.Info("Validation failed")
		utils.RespondWithError(w, err.Error())
		return
	}
	err = s.Store.CreateCategory(&cat)
	if err != nil {
		log.Warn(err.Error())
		utils.RespondWithError(w, "Could not create category.")
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, cat)
}

func (s *Server) deleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "Wrong parameter type")
		return
	}
	err = s.Store.DeleteCategory(id)
	if err != nil {
		log.Warning(err)
		utils.RespondWithError(w, "Could not delete category.")
		return
	}
	utils.RespondWithJSON(w, http.StatusAccepted, nil)
}
