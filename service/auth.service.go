package service

import "github.com/s1Sharp/s1-tts-restapi/internal/models"

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.DBUserResponse, error)
	SignInUser(*models.SignInInput) (*models.DBUserResponse, error)
}
