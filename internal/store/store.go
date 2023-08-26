package store

import (
	"database/sql"
	"fmt"
	"os"

	model "github.com/AlexCorn999/server-social-network/internal/user"
	_ "github.com/lib/pq"
)

type Store interface {
	Open() error
	Close() error
	Create(user *model.User) error
	GetUser(id int) (*model.User, error)
	NewAge(userID int, newAge string) error
	AddFriends(user1, user2 int) error
	AllFriends(id int) (string, error)
	Delete(userID int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewStore() *PostgresStore {
	return &PostgresStore{
		db: nil,
	}
}

// host=172.28.0.2 port=5432 user=skillbox dbname=skillbox sslmode=disable password=5427
// Open открывает подключение к базе данных.
func (s *PostgresStore) Open() error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	fmt.Println("CONECTED DB")

	s.db = db
	return nil
}

// Close закрывает соединение с базой данных.
func (s *PostgresStore) Close() error {
	return s.db.Close()
}

// AddUser добавляет пользователя в хранилище.
func (s *PostgresStore) Create(user *model.User) error {
	_, err := s.db.Exec("insert into users (user_id,username,age) values ($1,$2,$3)",
		user.Id, user.Name, user.Age)
	return err
}

// getUser выводит пользователя из базы данных по id
func (s *PostgresStore) GetUser(id int) (*model.User, error) {
	var user model.User
	err := s.db.QueryRow("select * from users where user_id=$1", id).
		Scan(&user.Id, &user.Name, &user.Age)
	return &user, err
}

// NewUserAge устанавливает новый возраст пользователя.
func (s *PostgresStore) NewAge(userID int, newAge string) error {
	_, err := s.db.Exec("update users set age=$1 where user_id=$2",
		newAge, userID)
	return err
}

// AddFriends добавляет друзей пользователю.
func (s *PostgresStore) AddFriends(user1, user2 int) error {
	_, err := s.db.Exec("insert into friends (friend_one,friend_two) values ($1,$2)",
		user1, user2)
	return err
}

// AllUserFriends выводит всех друзей пользователя.
func (s *PostgresStore) AllFriends(id int) (string, error) {
	rows, err := s.db.Query("select u.username, u.age from friends f join users u on u.user_id=f.friend_one or u.user_id=f.friend_two where (f.friend_one=$1 or f.friend_two=$2) and u.user_id!=$3;", id, id, id)
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
func (s *PostgresStore) Delete(userID int) error {
	if _, err := s.db.Exec("delete from friends where friend_one=$1 or friend_two=$2",
		userID, userID); err != nil {
		return err
	}

	if _, err := s.db.Exec("delete from users where user_id=$1",
		userID); err != nil {
		return err
	}

	return nil
}
