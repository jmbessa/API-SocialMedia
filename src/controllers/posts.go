package controllers

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// @Summary Create a new post
// @Description Create a new post with the data sent in the request body
// @Tags posts
// @Accept  json
// @Produce  json
// @Param        post  body      models.Post  true  "Create Post"
// @Success      201  {object}  models.Post
// @Failure      400  {object}  object       "Bad Request"
// @Failure      401  {object}  object       "Unauthorized"
// @Failure      422  {object}  object       "Unprocessable Entity"
// @Failure      500  {object}  object       "Internal Server Error"
// @Router       /posts [post]
// @Security ApiKeyAuth
func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post
	if err = json.Unmarshal(bodyRequest, &post); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	post.AuthorID = userID

	if err = post.Prepare(); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewPostsRepository(db)
	post.ID, err = repository.Create(post)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)
}

// @Summary Get all posts
// @Description Retrieve all posts from the database
// @Tags posts
// @Produce json
// @Success 200 {array} models.Post
// @Failure 500 {object} object "Internal Server Error"
// @Router /posts [get]
func GetPosts(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewPostsRepository(db)
	posts, err := repository.Search(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

// @Summary Get a post by ID
// @Description Retrieve a post by its ID from the database
// @Tags posts
// @Produce json
// @Param postId path int true "Post ID"
// @Success 200 {object} models.Post
// @Failure 400 {object} object "Bad Request"
// @Failure 500 {object} object "Internal Server Error"
// @Router /posts/{postId} [get]
func GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID, err := strconv.ParseUint(params["postId"], 10, 64)
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

	repository := repositories.NewPostsRepository(db)
	post, err := repository.SearchByID(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, post)
}

// @Summary Update a post
// @Description Update a post with the data sent in the request body
// @Tags posts
// @Accept json
// @Produce json
// @Param postId path int true "Post ID"
// @Param post body models.Post true "Post data"
// @Success 204 {object} object
// @Failure 400 {object} object "Bad Request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 422 {object} object "Unprocessable Entity"
// @Failure 500 {object} object "Internal Server Error"
// @Router /posts/{postId} [put]
// @Security ApiKeyAuth
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	postID, err := strconv.ParseUint(params["postId"], 10, 64)
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

	repository := repositories.NewPostsRepository(db)
	postOnBank, err := repository.SearchByID(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if postOnBank.AuthorID != userID {
		responses.Error(w, http.StatusInternalServerError, errors.New("It is not possible to update a post that is not yours"))
		return
	}

	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post
	if err = json.Unmarshal(bodyRequest, &post); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = post.Prepare(); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = repository.Update(postID, post); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// @Summary Delete a post
// @Description Delete a post by its ID
// @Tags posts
// @Param postId path int true "Post ID"
// @Success 204 {object} object
// @Failure 400 {object} object "Bad Request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal Server Error"
// @Router /posts/{postId} [delete]
// @Security ApiKeyAuth
func DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	postID, err := strconv.ParseUint(params["postId"], 10, 64)
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

	repository := repositories.NewPostsRepository(db)
	postOnBank, err := repository.SearchByID(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if postOnBank.AuthorID != userID {
		responses.Error(w, http.StatusInternalServerError, errors.New("It is not possible to delete a post that is not yours"))
		return
	}

	if err = repository.Delete(postID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// @Summary Get all posts by user
// @Description Retrieve all posts created by a specific user
// @Tags posts
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {array} models.Post
// @Failure 400 {object} object "Bad Request"
// @Failure 500 {object} object "Internal Server Error"
// @Router /users/{userId}/posts [get]
func GetPostsPerUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["postId"], 10, 64)
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

	repository := repositories.NewPostsRepository(db)
	posts, err := repository.SearchByUser(userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

// @Summary Like a post
// @Description Increment the likes count for a post
// @Tags posts
// @Param postId path int true "Post ID"
// @Success 204 {object} object
// @Failure 400 {object} object "Bad Request"
// @Failure 500 {object} object "Internal Server Error"
// @Router /posts/{postId}/like [post]
// @Security ApiKeyAuth
func LikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["postId"], 10, 64)
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

	repository := repositories.NewPostsRepository(db)
	if err = repository.Like(userID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}

// @Summary Dislike a post
// @Description Decrement the likes count for a post
// @Tags posts
// @Param postId path int true "Post ID"
// @Success 204 {object} object
// @Failure 400 {object} object "Bad Request"
// @Failure 500 {object} object "Internal Server Error"
// @Router /posts/{postId}/dislike [post]
// @Security ApiKeyAuth
func DislikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["postId"], 10, 64)
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

	repository := repositories.NewPostsRepository(db)
	if err = repository.Dislike(userID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}
