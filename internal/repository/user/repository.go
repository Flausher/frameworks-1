package user

import (
	"errors"

	"user.com/m/internal/domain"
)

type UserRepository struct {
	bd     map[int]*domain.User
	nextID int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		bd:     make(map[int]*domain.User),
		nextID: 1,
	}
}

func (ur *UserRepository) CreateUser(user *domain.User) error {
	user.ID = ur.nextID
	ur.nextID++

	ur.bd[user.ID] = user
	return nil
}

func (ur *UserRepository) GetUserById(id int) (*domain.User, error) {
	user, ok := ur.bd[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (ur *UserRepository) GetUsers() []*domain.User {
	users := make([]*domain.User, 0)
	for _, user := range ur.bd {
		users = append(users, user)
	}
	return users
}
