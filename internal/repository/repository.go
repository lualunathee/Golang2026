package repository

import (
	"golang/internal/repository/_postgres"
	pgUsers "golang/internal/repository/_postgres/users"
	"golang/pkg/modules"
)

type UserRepository interface {
	CreateUser(u modules.User) (int, error)
	UpdateUser(id int, u modules.User) error
	GetUserByID(id int) (*modules.User, error)
	DeleteUserByID(id int) (int64, error)
	GetUsers() ([]modules.User, error)
}

type Repositories struct {
	UserRepository UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: pgUsers.NewUserRepository(db),
	}
}
