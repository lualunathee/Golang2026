package service

import (
	"errors"
	"practice8/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.CreateUser(user)
	assert.NoError(t, err)
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)
	user := &repository.User{ID: 2, Name: "John Doe"}
	email := "test@test.com"

	t.Run("User already exists", func(t *testing.T) {
		mockRepo.EXPECT().GetByEmail(email).Return(user, nil)
		err := userService.RegisterUser(user, email)
		assert.EqualError(t, err, "user with this email already exists")
	})

	t.Run("New User Success", func(t *testing.T) {
		mockRepo.EXPECT().GetByEmail(email).Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(nil)
		err := userService.RegisterUser(user, email)
		assert.NoError(t, err)
	})

	t.Run("Repository error on CreateUser", func(t *testing.T) {
		mockRepo.EXPECT().GetByEmail(email).Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(errors.New("db error"))
		err := userService.RegisterUser(user, email)
		assert.EqualError(t, err, "db error")
	})
}

func TestUpdateUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)
	user := &repository.User{ID: 2, Name: "Old Name"}

	t.Run("Empty name", func(t *testing.T) {
		err := userService.UpdateUserName(2, "")
		assert.EqualError(t, err, "name cannot be empty")
	})

	t.Run("User not found / repo error", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(2).Return(nil, errors.New("not found"))
		err := userService.UpdateUserName(2, "New Name")
		assert.EqualError(t, err, "not found")
	})

	t.Run("Successful update", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(2).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(user).DoAndReturn(func(u *repository.User) error {
			assert.Equal(t, "New Name", u.Name) // Verify name was changed
			return nil
		})
		err := userService.UpdateUserName(2, "New Name")
		assert.NoError(t, err)
	})

	t.Run("UpdateUser Fails", func(t *testing.T) {
		user.Name = "Old Name" // Reset
		mockRepo.EXPECT().GetUserByID(2).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(user).Return(errors.New("update failed"))
		err := userService.UpdateUserName(2, "New Name")
		assert.EqualError(t, err, "update failed")
	})
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	t.Run("Attempt to delete admin", func(t *testing.T) {
		err := userService.DeleteUser(1)
		assert.EqualError(t, err, "it is not allowed to delete admin user")
	})

	t.Run("Successful delete", func(t *testing.T) {
		mockRepo.EXPECT().DeleteUser(2).Return(nil)
		err := userService.DeleteUser(2)
		assert.NoError(t, err)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.EXPECT().DeleteUser(2).Return(errors.New("db error"))
		err := userService.DeleteUser(2)
		assert.EqualError(t, err, "db error")
	})
}