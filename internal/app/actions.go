package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// count of users
var user_id = 1

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	var u User
	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	// write new user in the map under the user_id
	s.Store[user_id] = &u
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(u.UserCreated()))
	user_id++
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
	w.Write([]byte(NowFriends(id_1, id_2)))
}

func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request) {

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	type UserIdForDelete struct {
		Target_id string `json: "target_id"`
	}

	var forDelete UserIdForDelete

	if err := json.Unmarshal(content, &forDelete); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	// get the userID
	idUser, err := strconv.Atoi(forDelete.Target_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found ...."))
		return
	}

	// looking for a user
	_, ok := s.Store[idUser]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found ...."))
		return
	}

	response := s.Store[idUser].UserDeleted()
	user_id--

	s.Store[idUser].DeleteFromFriends()
	delete(s.Store, idUser)

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
