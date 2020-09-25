package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"scotch-go-lang-rest-api/api/models"
	"scotch-go-lang-rest-api/api/responses"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) CreateVenue(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Venue Successfully Created"}

	user := r.Context().Value("userID").(float64)
	venue := &models.Venue{}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &venue)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	venue.Prepare()

	if err = venue.Validate(); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if vne, _ := venue.GetVenue(a.DB); vne != nil {
		resp["status"] = "failed"
		resp["message"] = "Venue already registered, please choose another one"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	venue.UserID = uint64(user)

	venueCreated, err := venue.Save(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["venue"] = venueCreated
	responses.JSON(w, http.StatusCreated, resp)
	return

}

func (a *App) GetVenues(w http.ResponseWriter, r *http.Request) {
	venues, err := models.GetVenues(a.DB)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, venues)
	return

}

func (a *App) UpdateVenue(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"satatus": "success", "message": "Venue Updated successfully"}

	vars := mux.Vars(r)

	user := r.Context().Value("userID").(float64)
	userID := uint64(user)

	id, _ := strconv.Atoi(vars["id"])

	venue, err := models.GetVenueById(id, a.DB)

	if venue.UserID != userID {
		resp["status"] = "failed"
		resp["message"] = "Unauthorized venue uodate"
		responses.JSON(w, http.StatusUnauthorized, resp)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	venueUpdate := models.Venue{}

	if err = json.Unmarshal(body, &venueUpdate); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	venueUpdate.Prepare()

	_, err = venueUpdate.UpdateVenue(id, a.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return

}

func (a *App) DeleteVenue(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Venue deleted successfully"}

	vars := mux.Vars(r)

	user := r.Context().Value("userID").(float64)
	userID := uint64(user)

	id, _ := strconv.Atoi(vars["id"])

	venue, err := models.GetVenueById(id, a.DB)

	if venue.UserID != userID {
		resp["status"] = "failed"
		resp["message"] = "Unauthorized venue delete"
		responses.JSON(w, http.StatusUnauthorized, resp)
		return
	}

	err = models.DeleteVenue(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return

}
