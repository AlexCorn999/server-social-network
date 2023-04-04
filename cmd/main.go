package main

import (
	"fmt"
	"log"
	"net/http"

	handler "github.com/AlexCorn999/server-social-network/internal/actions"
	repository "github.com/AlexCorn999/server-social-network/internal/storage"
	model "github.com/AlexCorn999/server-social-network/internal/user"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const hostName = ":8080"

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	fmt.Println("Starting server ...")

	// storage creation
	rs := &repository.Service{Store: make(map[int]*model.User)}
	srv := handler.Storage{rs}

	router.Post("/users", srv.Create)
	router.Post("/friends", srv.MakeFriends)
	router.Delete("/users/{user_id}", srv.DeleteUser)
	router.Get("/friends/{user_id}", srv.GetUser)
	router.Put("/users/{user_id}", srv.ChangeAge)

	log.Fatal(http.ListenAndServe(hostName, router))
}
