package apiserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"net/http"
	"strconv"

	"github.com/AlexCorn999/server-social-network/internal/logger"
	"github.com/AlexCorn999/server-social-network/internal/store"
	model "github.com/AlexCorn999/server-social-network/internal/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const hostName = ":9000"

// APIServer ...
type APIServer struct {
	storage *store.PostgresStore
	store.Store
	router *chi.Mux
}

// New APIServer
func New() *APIServer {
	return &APIServer{
		router:  chi.NewRouter(),
		storage: store.NewStore(),
	}
}

// Start APIServer
func (s *APIServer) Start() error {
	// добавил интерфейс
	s.Store = s.storage

	s.router.Use(middleware.Logger)
	s.configureRouter()
	if err := s.storage.Open(); err != nil {
		log.Fatal(err)
	}
	defer s.storage.Close()

	fmt.Println("Starting api server")

	return http.ListenAndServe(hostName, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.Post("/users", s.Create)
	s.router.Post("/friends", s.MakeFriends)
	s.router.Delete("/users/{user_id}", s.DeleteUser)
	s.router.Get("/friends/{user_id}", s.GetUser)
	s.router.Put("/users/{user_id}", s.ChangeAge)
}

// Create отвечает за создание пользователя и добавления в хранилище.
func (s *APIServer) Create(w http.ResponseWriter, r *http.Request) {
	content, err := io.ReadAll(r.Body)
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

	// записываем пользователя в хранилище
	u.Id = model.UserID
	s.Store.Create(&u)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(u.UserCreated()))
	model.UserID++
}

// MakeFriends добавляет пользователей в друзья.
func (s *APIServer) MakeFriends(w http.ResponseWriter, r *http.Request) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	defer r.Body.Close()

	type MakeFriends struct {
		Source_id string `json:"source_id"`
		Target_id string `json:"target_id"`
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

	// проверка на наличие пользователей в базе
	u1, err := s.Store.GetUser(id1)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.ForError(err)
		return
	}

	u2, err := s.Store.GetUser(id2)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.ForError(err)
		return
	}

	// дружба
	if err := s.Store.AddFriends(id1, id2); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	result := fmt.Sprintf("%s и %s теперь друзья.", u1.Name, u2.Name)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

// DeleteUser удаляет пользователя.
func (s *APIServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	// проверка на наличие пользователя в базе
	u, err := s.Store.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.ForError(err)
		return
	}

	if err := s.Store.Delete(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.ForError(err)
	}

	text := fmt.Sprintf("User %s has been deleted", u.Name)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

// GetUser выводит друзей пользователя.
func (s *APIServer) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	// проверка на наличие пользователя в базе
	_, err = s.Store.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.ForError(err)
		return
	}

	result, err := s.Store.AllFriends(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.ForError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

// ChangeAge меняет возраст пользователя.
func (s *APIServer) ChangeAge(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}
	defer r.Body.Close()

	type UserNewAge struct {
		New_age string `json:"new_age"`
	}

	var requestNewAge UserNewAge
	if err := json.Unmarshal(content, &requestNewAge); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	// проверка на наличие пользователя в базе
	u, err := s.Store.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.ForError(err)
		return
	}

	if err := s.Store.NewAge(id, requestNewAge.New_age); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}

	newU, err := s.Store.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.ForError(err)
		return
	}

	result := fmt.Sprintf("User %s, new_age %s", u.Name, newU.Age)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
