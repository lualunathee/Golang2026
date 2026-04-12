package usecase

import (
	"practice7/internal/entity"
)

type (
	UserInterface interface {
		RegisterUser(user *entity.User) (*entity.User, string, error)
		LoginUser(user *entity.LoginUserDTO) (string, error)
		GetUserByID(id string) (*entity.User, error)
		UpdateUserRole(id string, role string) error
	}
)