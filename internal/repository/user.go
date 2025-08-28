package repository

import (
	"context"

	"github.com/uchupx/saceri-chatbot-api/internal/models"
)

type UserRepoInterface interface {
	GetUser(ctx context.Context, id string) (*models.UserModel, error)
	CreateUser(ctx context.Context, user models.UserModel) (*models.UserModel, error)
	UpdateUser(ctx context.Context, user models.UserModel) (*models.UserModel, error)
	DeleteUser(ctx context.Context, id string) error
	GetAllUsers(ctx context.Context, keyword *string, limit, offset int) ([]models.UserModel, error)
	GetUserByOauthID(ctx context.Context, oauthID string) (*models.UserModel, error)
}
