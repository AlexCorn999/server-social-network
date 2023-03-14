package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const hostName = ":8080"

var user_id = 1

type User struct {
	Name    string  `json: "name"`
	Age     string  `json: "age"`
	Friends []*User `json: "friends"`
}

func (u *User) toString() string {
	return fmt.Sprintf("Name is %s and age is %s\n", u.Name, u.Age)
}

func (u *User) nowFriends(u1, u2 *User) string {
	return fmt.Sprintf("%s и %s теперь друзья", u1.Name, u2.Name)
}

type service struct {
	store map[int]*User
}

func main() {

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	fmt.Println("Starting server ...")

	srv := service{make(map[int]*User)}
	router.Post("/create", srv.Create)
	router.Post("/make_friends", srv.MakeFriends)
	router.Delete("/user", srv.DeleteUser)
	router.Get("/friends/{user_id}", srv.GetFriends)
	router.Put("/{user_id}", srv.ChangeAge)

	log.Fatal(http.ListenAndServe(hostName, router))
}

func (s *service) Create(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	var u User
	if err := json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// записываем нового user в мапу под id
	s.store[user_id] = &u
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User was created " + u.Name + "\n" + "User_id is " + strconv.Itoa(user_id)))
	user_id++
}

func (s *service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	type MakeFriends struct {
		Source_id string `json: "source_id"`
		Target_id string `json: "target_id"`
	}

	var mk MakeFriends

	if err := json.Unmarshal(content, &mk); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	id1, err := strconv.Atoi(mk.Source_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't convert user_id ..."))
		return
	}

	id2, err := strconv.Atoi(mk.Target_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't convert user_id ..."))
		return
	}

	// достаем юзера из мапы
	id_1, ok1 := s.store[id1]
	id_2, ok2 := s.store[id2]

	if !ok1 || !ok2 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("user not found"))
		return
	}

	// дружим пользователей, путем добавления в слайс
	id_1.Friends = append(id_1.Friends, id_2)
	id_2.Friends = append(id_2.Friends, id_1)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id_1.nowFriends(id_1, id_2)))
}

func (s *service) GetFriends(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// если пользователь найден
	_, ok := s.store[id]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("user not found ...."))
		return
	}

	userFriends, ok := s.store[id]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("User not found..."))
		return
	}

	response := ""
	for _, friend := range userFriends.Friends {
		response += friend.toString()
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))

}

func (s *service) DeleteUser(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	type UserIdForDelete struct {
		Target_id string `json: "target_id"`
	}

	var ud UserIdForDelete

	if err := json.Unmarshal(content, &ud); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// получаем user_id
	idUser, err := strconv.Atoi(ud.Target_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("user not found ...."))
		return
	}

	// если пользователь найден
	_, ok := s.store[idUser]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("user not found ...."))
		return
	}

	response := fmt.Sprintf("User %s has been deleted", s.store[idUser].Name)

	user_id--
	s.store[idUser] = nil

	delete(s.store, idUser)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (s *service) ChangeAge(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// если пользователь найден
	_, ok := s.store[id]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("user not found ...."))
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	type UserNewAge struct {
		New_age string `json: "new_age"`
	}

	var una UserNewAge

	if err := json.Unmarshal(content, &una); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	s.store[id].Age = una.New_age

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Возраст пользователя %s успешно обновлен на %s", s.store[id].Name, s.store[id].Age)
	return
}
