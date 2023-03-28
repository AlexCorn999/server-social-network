package repository

import (
	handler "github.com/AlexCorn999/server-social-network/internal/actions"
	model "github.com/AlexCorn999/server-social-network/internal/user"
)

func New() *handler.Service {
	return &handler.Service{
		Store: make(map[int]*model.User),
	}
}
