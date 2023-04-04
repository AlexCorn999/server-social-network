package repository

import (
	"fmt"

	model "github.com/AlexCorn999/server-social-network/internal/user"
)

type Service struct {
	Store map[int]*model.User
}

func (s *Service) AddUser(user *model.User) {
	s.Store[model.User_id] = user
}

func (s *Service) AddFriends(userID1, userID2 int) (string, error) {

	user1, ok1 := s.Store[userID1]
	if !ok1 {
		return "", fmt.Errorf("user %d not found ...", userID1)
	}

	user2, ok2 := s.Store[userID2]
	if !ok2 {
		return "", fmt.Errorf("user %d not found ...", userID2)
	}

	user1.Friends = append(user1.Friends, user2)
	user2.Friends = append(user2.Friends, user1)
	result := model.NowFriends(user1, user2)
	return result, nil
}

func (s *Service) UserDelete(userID int) (string, error) {
	_, ok := s.Store[userID]
	if !ok {
		return "", fmt.Errorf("user %d not found ...", userID)
	}

	response := s.Store[userID].UserDeleted()
	s.Store[userID].DeleteFromFriends()
	delete(s.Store, userID)
	model.User_id--
	return response, nil
}

func (s *Service) AllUserFriends(userID int) (string, error) {
	userFriends, ok := s.Store[userID]
	if !ok {
		return "", fmt.Errorf("user %d not found ...", userID)
	}

	response := ""
	for _, friend := range userFriends.Friends {
		response += friend.FriendsToString()
	}
	return response, nil
}

func (s *Service) NewUserAge(userID int, newAge string) (string, error) {
	_, ok := s.Store[userID]
	if !ok {
		return "", fmt.Errorf("user %d not found ...", userID)
	}
	s.Store[userID].Age = newAge
	response := s.Store[userID].NewAge()
	return response, nil
}
