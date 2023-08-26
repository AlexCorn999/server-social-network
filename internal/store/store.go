package store

import (
	"database/sql"
	"fmt"
	"os"

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

// host=172.28.0.2 port=5432 user=skillbox dbname=skillbox sslmode=disable password=5427
// Open открывает подключение к базе данных.
func (s *Store) Open() error {
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
