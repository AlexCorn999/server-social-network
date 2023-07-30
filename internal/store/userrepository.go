package store

import (
	"fmt"

	model "github.com/AlexCorn999/server-social-network/internal/user"
)

type UserRepository struct {
	store *Store
}

// AddUser добавляет пользователя в хранилище.
func (u *UserRepository) Create(user *model.User) error {
	_, err := u.store.db.Exec("insert into users (user_id,username,age) values ($1,$2,$3)",
		user.Id, user.Name, user.Age)
	return err
}

// getUser выводит пользователя из базы данных по id
func (u *UserRepository) GetUser(id int) (*model.User, error) {
	var user model.User
	err := u.store.db.QueryRow("select * from users where user_id=$1", id).
		Scan(&user.Id, &user.Name, &user.Age)
	return &user, err
}

// NewUserAge устанавливает новый возраст пользователя.
func (u *UserRepository) NewAge(userID int, newAge string) error {
	_, err := u.store.db.Exec("update users set age=$1 where user_id=$2",
		newAge, userID)
	return err
}

// AddFriends добавляет друзей пользователю.
func (u *UserRepository) AddFriends(user1, user2 int) error {
	_, err := u.store.db.Exec("insert into friends (friend_one,friend_two) values ($1,$2)",
		user1, user2)
	return err
}

// AllUserFriends выводит всех друзей пользователя.
func (u *UserRepository) AllFriends(id int) (string, error) {
	rows, err := u.store.db.Query("select u.username, u.age from friends f join users u on u.user_id=f.friend_one or u.user_id=f.friend_two where (f.friend_one=$1 or f.friend_two=$2) and u.user_id!=$3;", id, id, id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var u model.User
		if err = rows.Scan(&u.Name, &u.Age); err != nil {
			return "", err
		}
		users = append(users, u)
	}

	if rows.Err(); err != nil {
		return "", err
	}

	result := ""
	for _, user := range users {
		result += fmt.Sprintf("User %s, age %s\n", user.Name, user.Age)
	}

	return result, nil
}

// UserDelete удаляет пользователя из хранилища.
func (u *UserRepository) Delete(userID int) error {
	if _, err := u.store.db.Exec("delete from friends where friend_one=$1 or friend_two=$2",
		userID, userID); err != nil {
		return err
	}

	if _, err := u.store.db.Exec("delete from users where user_id=$1",
		userID); err != nil {
		return err
	}

	return nil
}
