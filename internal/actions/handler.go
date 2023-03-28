package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	model "github.com/AlexCorn999/server-social-network/internal/user"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	Store map[int]*model.User
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	var u model.User
	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	// write new user in the map under the user_id
	s.Store[model.User_id] = &u
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(u.UserCreated()))
	model.User_id++
}

func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
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
		return
	}

	id1, err := strconv.Atoi(request.Source_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can't convert user_id ..."))
		return
	}

	id2, err := strconv.Atoi(request.Target_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can't convert user_id ..."))
		return
	}

	id_1, ok1 := s.Store[id1]
	id_2, ok2 := s.Store[id2]

	if !ok1 || !ok2 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	id_1.Friends = append(id_1.Friends, id_2)
	id_2.Friends = append(id_2.Friends, id_1)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(model.NowFriends(id_1, id_2)))
}

func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	// if the user is found
	_, ok := s.Store[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found..."))
		return
	}

	response := s.Store[id].UserDeleted()
	model.User_id--

	s.Store[id].DeleteFromFriends()
	delete(s.Store, id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (s *Service) GetUser(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	// if the user is found
	userFriends, ok := s.Store[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found..."))
		return
	}

	response := ""
	for _, friend := range userFriends.Friends {
		response += friend.FriendsToString()
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (s *Service) ChangeAge(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	// if the user is found
	_, ok := s.Store[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found ...."))
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
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
		return
	}

	s.Store[id].Age = requestNewAge.New_age

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s.Store[id].NewAge()))
}
