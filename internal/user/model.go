package model

import (
	"fmt"
)

// count of users
var UserID = 1

type User struct {
	Name    string  `json: "name"`
	Age     string  `json: "age"`
	Friends []*User `json: "friends"`
}

// UserCreated оповещает о добавлении пользователя.
func (u *User) UserCreated() string {
	return fmt.Sprintf("User was created %s\nUserID : %d.\n", u.Name, UserID)
}

// NowFriends оповещает о дружбе двух пользователей.
func NowFriends(u1, u2 *User) string {
	return fmt.Sprintf("%s и %s теперь друзья.", u1.Name, u2.Name)
}

// UserDeleted оповещает об удалении пользователя.
func (u *User) UserDeleted() string {
	return fmt.Sprintf("User %s has been deleted", u.Name)
}

// Выводит информацию о друге у пользователя.
func (u *User) FriendsToString() string {
	return fmt.Sprintf("Name is %s and age is %s\n", u.Name, u.Age)
}

// NewAge меняет возраст пользователя.
func (u *User) NewAge(age string) string {
	u.Age = age
	return fmt.Sprintf("Возраст пользователя %s успешно обновлен на %s", u.Name, u.Age)
}
