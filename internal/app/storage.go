package app

type Service struct {
	Store map[int]*User
}

func New() *Service {
	return &Service{
		Store: make(map[int]*User),
	}
}
