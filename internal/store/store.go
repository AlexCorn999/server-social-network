package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

func NewStore() *Store {
	return &Store{
		db:             nil,
		userRepository: nil,
	}
}

// Open открывает подключение к базе данных.
func (s *Store) Open() error {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=ninja dbname=ninja sslmode=disable password=5427")
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
func (s *Store) Close() error {
	return s.db.Close()
}

// User позволяет инкапсулировать работу с юзером в хранилище.
func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
