package service

import "github.com/s1Sharp/s1-tts-restapi/internal/models"

type UserService interface {
	FindUserById(id string) (*models.DBResponse, error)
	FindUserByEmail(email string) (*models.DBResponse, error)
	UpdateUserById(id string, data *models.UpdateInput) (*models.DBResponse, error)
	RemoveUserById(id string) error
}
