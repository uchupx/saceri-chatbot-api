package service

import (
	"context"

	"github.com/labstack/echo/v4"
	_type "github.com/uchupx/saceri-chatbot-api/internal/api/handlers/type"
	"github.com/uchupx/saceri-chatbot-api/internal/models"
	"github.com/uchupx/saceri-chatbot-api/internal/repository"
	"github.com/uchupx/saceri-chatbot-api/pkg/apierror"
)

type UserService struct {
	UserRepo repository.UserRepoInterface
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.UserModel, *apierror.APIerror) {
	users, err := s.UserRepo.GetAllUsers(100, 0)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	return users, nil
}

func (s *UserService) Register(ctx context.Context, body _type.RegisterRequest) (*models.UserModel, error) {
	user := models.UserModel{
		Username: body.Username,
		Password: body.Password,
		Name:     body.Name,
	}

	resp, err := s.UserRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
