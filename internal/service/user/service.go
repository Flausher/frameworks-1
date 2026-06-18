package user

import (
	"errors"

	"user.com/m/internal/domain"
	"user.com/m/internal/repository/user"
)

type UserService struct {
	ur *user.UserRepository
}

func NewUserService(ur *user.UserRepository) *UserService {
	return &UserService{
		ur: ur,
	}
}

func (us *UserService) CreateUser(user *domain.User) error {
	if len(user.Name) == 0 {
		return errors.New("name cannot be empty")
	}

	if len(user.Email) == 0 {
		return errors.New("email cannot be empty")
	}

	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if user.Age <= 0 {
		return errors.New("age must be greater than 0")
	}

	return us.ur.CreateUser(user)
}

func (us *UserService) GetUserById(id int) (*domain.User, error) {
	return us.ur.GetUserById(id)
}

func (us *UserService) GetUsers() []*domain.User {
	return us.ur.GetUsers()
}
