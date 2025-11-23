package usecase

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/schema"
	"github.com/rafaeldepontes/go-full-crud/internal/repository"
	"github.com/rafaeldepontes/go-full-crud/internal/util"
	log "github.com/sirupsen/logrus"
)

type userUpdateForm struct {
	Email     *string
	Birthdate *string
}

type userRequest struct {
	Username string
	Password string
}

type usernameRequest struct {
	Username string
}

type UserHandler struct {
	UserRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepo: userRepo,
	}
}

func (uh *UserHandler) FindUserById(w http.ResponseWriter, r *http.Request) {
	if r.PathValue("id") == "" {
		log.Error(util.ErrorBlankId)
		util.BadRequestErrorHandler(w, util.ErrorBlankId)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Error(err)
		util.InternalErrorHandler(w)
		return
	}

	user, err := uh.UserRepo.FindUserById(&id)
	if err != nil {
		log.Error(err)
		util.RequestErrorHandler(w, util.ErrorUserNotFound, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) FindByUsername(w http.ResponseWriter, r *http.Request) {
	var username usernameRequest
	var decoder *schema.Decoder = schema.NewDecoder()
	err := decoder.Decode(&username, r.URL.Query())

	if err != nil {
		log.Error(err)
		util.InternalErrorHandler(w)
		return
	}

	if username.Username == "" {
		log.Error(util.ErrorBlankUsername)
		util.BadRequestErrorHandler(w, util.ErrorBlankUsername)
		return
	}

	user, err := uh.UserRepo.FindUserByUsername(&username.Username)
	if err != nil {
		log.Error(err)
		util.RequestErrorHandler(w, util.ErrorUserNotFound, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var params userRequest
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Error(err)
		util.InternalErrorHandler(w)
		return
	}

	if params.Username == "" || params.Password == "" {
		log.Error(util.ErrorUserCredentials)
		util.BadRequestErrorHandler(w, util.ErrorUserCredentials)
		return
	}

	var created_at time.Time = time.Now()
	var user repository.User = repository.User{
		Username:  &params.Username,
		Password:  &params.Password,
		CreatedAt: &created_at,
	}

	err = uh.UserRepo.CreateUser(&user)
	if err != nil {
		log.Error(err)
		util.InternalErrorHandler(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (uh *UserHandler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.PathValue("id") == "" {
		log.Error(util.ErrorBlankId)
		util.BadRequestErrorHandler(w, util.ErrorBlankId)
		return
	}

	var params userUpdateForm
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Error(err)
		util.InternalErrorHandler(w)
		return
	}

	timestamp := time.Now()
	user := repository.User{
		Email:     params.Email,
		Birthdate: params.Birthdate,
		UpdatedAt: &timestamp,
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Error(err)
		util.InternalErrorHandler(w)
		return
	}

	err = uh.UserRepo.UpdateUserDetails(&user, id)
	if err != nil {
		log.Error(err)
		util.InternalErrorHandler(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	if r.PathValue("id") == "" {
		log.Error(util.ErrorBlankId)
		util.BadRequestErrorHandler(w, util.ErrorBlankId)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Error(err)
		util.InternalErrorHandler(w)
		return
	}

	err = uh.UserRepo.DeleteUserById(id)
	if err != nil {
		log.Error(err)
		util.InternalErrorHandler(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
