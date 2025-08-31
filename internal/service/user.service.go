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

func (s *UserService) GetUsers(ctx context.Context, params _type.GetQuery) ([]models.UserModel, *apierror.APIerror) {
	users, err := s.UserRepo.GetAllUsers(ctx, params.Keyword, params.Limit(), params.Offset())
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	return users, nil
}

func (s *UserService) GetUserByOauthID(ctx context.Context, id string) (*models.UserModel, *apierror.APIerror) {
	user, err := s.UserRepo.GetUserByOauthID(ctx, id)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	if user == nil {
		return nil, apierror.NewAPIError(echo.ErrNotFound.Code, apierror.ERR_NOT_FOUND)
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.UserModel, *apierror.APIerror) {
	user, err := s.UserRepo.GetUser(ctx, id)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	if user == nil {
		return nil, apierror.NewAPIError(echo.ErrNotFound.Code, apierror.ERR_NOT_FOUND)
	}

	return user, nil
}

func (s *UserService) Register(ctx context.Context, body _type.RegisterRequest) (*models.UserModel, error) {
	user := models.UserModel{
		Username: body.Username,
		Name:     body.Name,
		OauthID:  body.OauthID,
	}

	resp, err := s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, body _type.UserUpdateRequest) (*models.UserModel, *apierror.APIerror) {
	user, err := s.GetUserByOauthID(ctx, body.ID())
	if err != nil {
		return nil, err
	}

	if body.Name != nil {
		user.Name = *body.Name
	}

	user, er := s.UserRepo.UpdateUser(ctx, *user)
	if er != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, er)
	}

	return user, nil
}
