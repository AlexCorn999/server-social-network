package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/AlexCorn999/server-social-network/internal/logger"
	repository "github.com/AlexCorn999/server-social-network/internal/storage"
	model "github.com/AlexCorn999/server-social-network/internal/user"
	"github.com/go-chi/chi/v5"
)

type Storage struct {
	*repository.Service
}

func (s *Storage) Create(w http.ResponseWriter, r *http.Request) {

	content, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	defer r.Body.Close()

	var u model.User
	if err = json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	// write new user in the map under the user_id
	s.AddUser(&u)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(u.UserCreated()))
	model.User_id++
}

func (s *Storage) MakeFriends(w http.ResponseWriter, r *http.Request) {

	content, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	defer r.Body.Close()

	type MakeFriends struct {
		Source_id string `json: "source_id"`
		Target_id string `json: "target_id"`
	}

	var request MakeFriends

	if err := json.Unmarshal(content, &request); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	id1, err := strconv.Atoi(request.Source_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can't convert user_id ..."))
		logger.ForError(fmt.Errorf("can't convert user_id : %v", err))
		return
	}

	id2, err := strconv.Atoi(request.Target_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can't convert user_id ..."))
		logger.ForError(fmt.Errorf("can't convert user_id : %v", err))
		return
	}

	text, err := s.AddFriends(id1, id2)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

func (s *Storage) DeleteUser(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))

	logger.ForError(err)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	// if the user is found
	text, err := s.UserDelete(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found..."))
		logger.ForError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

func (s *Storage) GetUser(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	// if the user is found
	text, err := s.AllUserFriends(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found..."))
		logger.ForError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

func (s *Storage) ChangeAge(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	content, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	defer r.Body.Close()

	type UserNewAge struct {
		New_age string `json: "new_age"`
	}

	var requestNewAge UserNewAge

	if err := json.Unmarshal(content, &requestNewAge); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	// if the user is found
	text, err := s.NewUserAge(id, requestNewAge.New_age)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found ...."))
		logger.ForError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}
