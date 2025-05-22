package services

import (
	"go-fiber/data/repositories"
)

type UserService interface {
	GetAllUsers()
}

type UserServiceImpl struct {
	userRepo *repositories.UserRepositoryImpl
}

// GetAllUsers implements UserService.
func (u *UserServiceImpl) GetAllUsers() {
	panic("unimplemented")
}

func NewUserService(userRepo *repositories.UserRepositoryImpl) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}