package usercase

import (
	"github.com/AlexCorn999/server-social-network/internal/store"
	model "github.com/AlexCorn999/server-social-network/internal/user"
)

type UserCase interface {
	AddUser(usr *model.User) (string, error)
	MakeFriends(usr1, usr2 int) (string, error)
	DeleteUser(id int) (string, error)
	GetFriends(id int) (string, error)
	ChangeAge(id int, newAge string) (string, error)
	OpenConnection() error
	CloseConnection() error
}

type UserCases struct {
	store store.Store
}

func NewUserCases() *UserCases {
	cases := store.NewStore()
	return &UserCases{
		store: cases,
	}
}

// OpenConnection открывает соединение с базой
func (u *UserCases) OpenConnection() error {
	return u.store.Open()
}

// CloseConnection закрывает подключение с базой
func (u *UserCases) CloseConnection() error {
	return u.store.Close()
}

// AddUser добавляет пользователя в хранилище
func (u *UserCases) AddUser(usr *model.User) (string, error) {
	usr.Id = model.UserID
	model.UserID++
	if err := u.store.Create(usr); err != nil {
		return "", err
	}
	return usr.UserCreated(), nil
}

// MakeFriends добавляет пользователей в друзья
func (u *UserCases) MakeFriends(usr1, usr2 int) (string, error) {
	u1, err := u.store.GetUser(usr1)
	if err != nil {
		return "", err
	}
	u2, err := u.store.GetUser(usr2)
	if err != nil {
		return "", err
	}

	if err := u.store.AddFriends(u1.Id, u2.Id); err != nil {
		return "", err
	}

	return model.NowFriends(u1, u2), nil
}

// DeleteUser удаляет пользователя из базы
func (u *UserCases) DeleteUser(id int) (string, error) {
	user, err := u.store.GetUser(id)
	if err != nil {
		return "", err
	}

	if err := u.store.Delete(user.Id); err != nil {
		return "", err
	}
	return user.UserDeleted(), nil
}

// GetFrineds выводит всех друзей пользователя
func (u *UserCases) GetFriends(id int) (string, error) {
	_, err := u.store.GetUser(id)
	if err != nil {
		return "", err
	}
	return u.store.AllFriends(id)
}

// ChangeAge меняет возраст пользователя
func (u *UserCases) ChangeAge(id int, newAge string) (string, error) {
	user, err := u.store.GetUser(id)
	if err != nil {
		return "", err
	}

	if err := u.store.NewAge(id, newAge); err != nil {
		return "", err
	}

	return user.NewAge(newAge), nil
}
