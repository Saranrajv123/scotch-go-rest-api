package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"scotch-go-lang-rest-api/api/models"
	"scotch-go-lang-rest-api/api/responses"
	"scotch-go-lang-rest-api/utils"
)

func (a *App) UserSignUp(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Registered Successfully"}

	user := &models.User{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	usr, _ := user.GetUser(a.DB)

	if usr != nil {
		resp["status"] = "failed"
		resp["message"] = "User already registered, please login"

		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	user.Prepare()

	err = user.Validate("")

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userCreated, err := user.SaveUser(a.DB)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["user"] = userCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "login successfully"}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user.Prepare()

	err = user.Validate("login")

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	usr, err := user.GetUser(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if usr == nil {
		resp["status"] = "failed"
		resp["message"] = "Login failed"

		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	err = models.CheckPasswordHash(user.Password, usr.Password)

	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "Login failed"

		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	token, err := utils.EncodeAuthToken(usr.ID)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["token"] = token
	responses.JSON(w, http.StatusOK, resp)
	return
}
