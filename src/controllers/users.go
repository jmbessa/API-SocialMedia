package controllers

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// @Summary Create a new user
// @Description Create a new user with the provided data
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "New user data"
// @Success 201 {object} models.User
// @Failure 400 {object} object "Bad Request"
// @Failure 500 {object} object "Internal Server Error"
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
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

	if err = user.Prepare("register"); err != nil {
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
	user.ID, err = repository.Create(user)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}

// @Summary Get all users
// @Description Retrieve all users, optionally filtered by name or nickname
// @Tags users
// @Accept json
// @Produce json
// @Param user query string false "Name or nickname of the user to filter"
// @Success 200 {array} models.User
// @Failure 500 {object} object "Internal Server Error"
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("user"))

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	users, err := repository.Search(nameOrNick)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

// @Summary Get user by ID
// @Description Retrieve a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} object "Bad Request"
// @Failure 500 {object} object "Internal Server Error"
// @Router /users/{userID} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
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
	user, err := repository.SearchByID(userId)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

// @Summary Delete user by ID
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 204 {object} object
// @Failure 400 {object} object "Bad Request"
// @Failure 403 {object} object "Forbidden"
// @Failure 500 {object} object "Internal Server Error"
// @Router /users/{userID} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIDToken, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDToken {
		responses.Error(w, http.StatusForbidden, errors.New("It is not possible to delete a user other than yours"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	if err := repository.Delete(userID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// @Summary Update user by ID
// @Description Update a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 204 {object} object
// @Failure 400 {object} object "Bad Request"
// @Failure 403 {object} object "Forbidden"
// @Failure 500 {object} object "Internal Server Error"
// @Router /users/{userID} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	userIDToken, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDToken {
		responses.Error(w, http.StatusForbidden, errors.New("It is not possible to update a user other than yours"))
		return
	}

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

	if err = user.Prepare("edit"); err != nil {
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
	if err = repository.Update(userID, user); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// @Summary Follow user by ID
// @Description Follow a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 204 {object} object
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal Server Error"
// @Router /users/{userID}/follow [post]
func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if userID == followerID {
		responses.Error(w, http.StatusForbidden, errors.New("Is not possible to follow yourself"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	if err = repository.Follow(userID, followerID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// @Summary Unfollow user by ID
// @Description Unfollow a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 204 {object} object
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal Server Error"
// @Router /users/{userID}/unfollow [post]
func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followerID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)

	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if userID == followerID {
		responses.Error(w, http.StatusForbidden, errors.New("Is not possible to unfollow yourself"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	if err = repository.Unfollow(userID, followerID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// @Summary Search followers of user
// @Description Search followers of a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {array} models.User
// @Failure 400 {object} object "Bad Request"
// @Failure 500 {object} object "Internal Server Error"
// @Router /users/{userID}/followers [get]
func SearchFollowers(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userID, err := strconv.ParseUint(parameters["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	followers, err := repository.SearchFollowers(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

// @Summary Search following users of user
// @Description Search users followed by a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {array} models.User
// @Failure 400 {object} object "Bad Request"
// @Failure 500 {object} object "Internal Server Error"
// @Router /users/{userID}/following [get]
func SearchFollowing(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userID, err := strconv.ParseUint(parameters["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	users, err := repository.SearchFollowing(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

// @Summary Update user password
// @Description Update the password of a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 204 {object} object
// @Failure 400 {object} object "Bad Request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal Server Error"
// @Router /users/{userID}/update-password [post]
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userIDToken, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userID"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if userIDToken != userID {
		responses.Error(w, http.StatusForbidden, errors.New("it is not possible to update the password of a user other than yours"))
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)

	var password models.Password

	if err = json.Unmarshal(requestBody, &password); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
	}
	defer db.Close()

	repository := repositories.NewUsersRepository(db)
	savedPassword, err := repository.SearchPassword(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.VerifyPassword(savedPassword, password.Current); err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("the current password does not match the one saved in the bank"))
		return
	}

	hashedPassword, err := security.Hash(password.New)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
	}

	if err = repository.UpdatePassword(userID, string(hashedPassword)); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
