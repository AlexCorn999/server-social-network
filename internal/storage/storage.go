package storage

import (
	"fmt"

	model "github.com/AlexCorn999/server-social-network/internal/user"
)

type Service struct {
	store map[int]*model.User
}

func NewService() *Service {
	return &Service{
		store: make(map[int]*model.User),
	}
}

// AddUser добавляет пользователя в хранилище.
func (s *Service) AddUser(user *model.User) {
	s.store[model.UserID] = user
}

// AddFriends добавляет друзей пользователю.
func (s *Service) AddFriends(userID1, userID2 int) (string, error) {
	user1, ok1 := s.store[userID1]
	if !ok1 {
		return "", fmt.Errorf("user %d not found", userID1)
	}

	user2, ok2 := s.store[userID2]
	if !ok2 {
		return "", fmt.Errorf("user %d not found", userID2)
	}

	user1.Friends = append(user1.Friends, user2)
	user2.Friends = append(user2.Friends, user1)
	return model.NowFriends(user1, user2), nil
}

// UserDelete удаляет пользователя из хранилища и удаляет его у всех его друзей.
func (s *Service) UserDelete(userID int) (string, error) {
	userForDelete, ok := s.store[userID]
	if !ok {
		return "", fmt.Errorf("user %d not found", userID)
	}

	for _, friend := range userForDelete.Friends {
		for i := 0; i < len(friend.Friends); i++ {
			if friend.Friends[i] == userForDelete {
				friend.Friends = append(friend.Friends[:i], friend.Friends[i+1:]...)
			}
		}
	}

	response := userForDelete.UserDeleted()
	delete(s.store, userID)
	model.UserID--
	return response, nil
}

// AllUserFriends выводит всех друзей пользователя.
func (s *Service) AllUserFriends(userID int) (string, error) {
	userFriends, ok := s.store[userID]
	if !ok {
		return "", fmt.Errorf("user %d not found", userID)
	}

	response := ""
	for _, friend := range userFriends.Friends {
		response += friend.FriendsToString()
	}
	return response, nil
}

// NewUserAge устанавливает новый возраст пользователя.
func (s *Service) NewUserAge(userID int, newAge string) (string, error) {
	_, ok := s.store[userID]
	if !ok {
		return "", fmt.Errorf("user %d not found", userID)
	}
	return s.store[userID].NewAge(newAge), nil
}
