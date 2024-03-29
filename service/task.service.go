package service

import "github.com/s1Sharp/s1-tts-restapi/internal/models"

type TaskService interface {
	CreateTask(*models.CreateTaskScheme) (*models.DBTaskScheme, error)
	GetUserTasks(string, int, int) ([]*models.DBTaskScheme, error)
	TaskById(string) (*models.DBTaskScheme, error)
}
