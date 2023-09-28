package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ListService interface {
	List(ctx context.Context, ids []uuid.UUID) ([]*models.UserPreview, error)
}

func NewListService(userDAO dao.UserRepository) ListService {
	return &listServiceImpl{
		userDAO: userDAO,
	}
}

type listServiceImpl struct {
	userDAO dao.UserRepository
}

func (s *listServiceImpl) List(ctx context.Context, ids []uuid.UUID) ([]*models.UserPreview, error) {
	users, err := s.userDAO.List(ctx, ids)
	if err != nil {
		return nil, goerrors.Join(ErrListUsers, err)
	}

	return lo.Map(users, func(item *dao.UserModel, _ int) *models.UserPreview {
		return &models.UserPreview{
			ID:        item.ID,
			FirstName: lo.Ternary(item.Profile.Username == "", item.Identity.FirstName, ""),
			LastName:  lo.Ternary(item.Profile.Username == "", item.Identity.LastName, ""),
			Username:  item.Profile.Username,
			Slug:      item.Profile.Slug,
			CreatedAt: item.CreatedAt,
		}
	}), nil
}
