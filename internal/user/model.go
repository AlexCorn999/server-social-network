package model

import (
	"fmt"
)

// count of users
var User_id = 1

type User struct {
	Name    string  `json: "name"`
	Age     string  `json: "age"`
	Friends []*User `json: "friends"`
}

func (u *User) UserCreated() string {
	return fmt.Sprintf("User was created %s\nUser_id : %d.\n", u.Name, User_id)
}

func NowFriends(u1, u2 *User) string {
	return fmt.Sprintf("%s и %s теперь друзья.", u1.Name, u2.Name)
}

func (u *User) UserDeleted() string {
	return fmt.Sprintf("User %s has been deleted", u.Name)
}

func (u *User) FriendsToString() string {
	return fmt.Sprintf("Name is %s and age is %s\n", u.Name, u.Age)
}

func (u *User) NewAge() string {
	return fmt.Sprintf("Возраст пользователя %s успешно обновлен на %s", u.Name, u.Age)
}

func (u *User) DeleteFromFriends() {
	var allFriends []*User

	// adding friends of the user we want to delete
	for _, friend := range u.Friends {
		allFriends = append(allFriends, friend)
	}

	// deleting a user from each friend
	for i := 0; i < len(allFriends); i++ {

		for j := 0; j < len(allFriends[i].Friends); j++ {

			if allFriends[i].Friends[j] == u {
				newFriends := append(allFriends[i].Friends[:j], allFriends[i].Friends[j+1:]...)
				allFriends[i].Friends = newFriends
			}

		}
	}
}
