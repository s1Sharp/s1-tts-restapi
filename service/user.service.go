package service

import "github.com/s1Sharp/s1-tts-restapi/internal/models"

type UserService interface {
	FindUserById(id string) (*models.DBUserResponse, error)
	FindUserByEmail(email string) (*models.DBUserResponse, error)
	UpdateUserById(id string, data *models.UpdateInput) (*models.DBUserResponse, error)
	RemoveUserById(id string) error
}
