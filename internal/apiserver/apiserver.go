package apiserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"net/http"
	"strconv"

	"github.com/AlexCorn999/server-social-network/internal/logger"
	model "github.com/AlexCorn999/server-social-network/internal/user"
	"github.com/AlexCorn999/server-social-network/internal/usercase"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const hostName = ":9000"

// APIServer ...
type APIServer struct {
	usercase.UserCase
	router *chi.Mux
}

// New APIServer
func New() *APIServer {
	return &APIServer{
		router:   chi.NewRouter(),
		UserCase: usercase.NewUserCases(),
	}
}

// Start APIServer
func (s *APIServer) Start() error {
	s.router.Use(middleware.Logger)
	s.configureRouter()
	if err := s.UserCase.OpenConnection(); err != nil {
		log.Fatal(err)
	}
	defer s.UserCase.CloseConnection()
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

	// запись пользователя в базу
	result, err := s.UserCase.AddUser(&u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(result))
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

	// добавление пользователей в друзья в базу
	result, err := s.UserCase.MakeFriends(id1, id2)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}
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
	// удаление пользователя из базы
	result, err := s.UserCase.DeleteUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.ForError(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
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

	// получение всех друзей
	result, err := s.UserCase.GetFriends(id)
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

	// смена возраста пользователя
	result, err := s.UserCase.ChangeAge(id, requestNewAge.New_age)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		logger.ForError(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
