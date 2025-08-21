package repository

import "github.com/uchupx/saceri-chatbot-api/internal/models"

type UserRepoInterface interface {
	GetUser(id string) (*models.UserModel, error)
	CreateUser(user models.UserModel) (*models.UserModel, error)
	UpdateUser(user models.UserModel) (*models.UserModel, error)
	DeleteUser(id string) error
	GetAllUsers(limit, offset int) ([]models.UserModel, error)
}
