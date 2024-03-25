package controllers

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// @Summary Authenticate user
// @Description Authenticate the user by checking the provided credentials
// @Tags authentication
// @Accept json
// @Produce plain
// @Param credentials body models.User true "User credentials"
// @Success 200 {string} string "Authentication token"
// @Failure 400 {object} object "Bad Request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal Server Error"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User

	if err = json.Unmarshal(requestBody, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	userSavedDatabase, err := repository.SearchByEmail(user.Email)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.VerifyPassword(
		userSavedDatabase.Password, user.Password); err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	token, err := authentication.CreateToken(userSavedDatabase.ID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	formattedToken := fmt.Sprintf("Token: %s", token)
	w.Write([]byte(formattedToken))

}
