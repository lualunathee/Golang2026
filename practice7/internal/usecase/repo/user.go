package repo

import (
	"fmt"
	"practice7/internal/entity"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (u *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	err := u.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) LoginUser(user *entity.LoginUserDTO) (*entity.User, error) {
	var userFromDB entity.User
	if err := u.DB.Where("username = ?", user.Username).First(&userFromDB).Error; err != nil {
		return nil, fmt.Errorf("Username Not Found: %v", err)
	}
	return &userFromDB, nil
}

func (u *UserRepo) GetUserByID(id string) (*entity.User, error) {
    var user entity.User
    if err := u.DB.Where("id = ?", id).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (u *UserRepo) UpdateUserRole(id string, role string) error {
    return u.DB.Model(&entity.User{}).Where("id = ?", id).Update("role", role).Error
}

