package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	goframework "github.com/a-novel/go-framework"
	"github.com/samber/lo"
)

const (
	MaxUserSearchLimit = 100
)

type SearchService interface {
	Search(ctx context.Context, query string, limit int, offset int) ([]*models.UserPreview, int, error)
}

func NewSearchService(userDAO dao.UserRepository) SearchService {
	return &searchServiceImpl{
		userDAO: userDAO,
	}
}

type searchServiceImpl struct {
	userDAO dao.UserRepository
}

func (s *searchServiceImpl) Search(ctx context.Context, query string, limit int, offset int) ([]*models.UserPreview, int, error) {
	if err := goframework.CheckMinMax(limit, 1, MaxUserSearchLimit); err != nil {
		return nil, 0, goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidSearchLimit, err)
	}

	users, total, err := s.userDAO.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, goerrors.Join(ErrSearchUsers, err)
	}

	return lo.Map(users, func(item *dao.UserModel, _ int) *models.UserPreview {
		return &models.UserPreview{
			FirstName: lo.Ternary(item.Profile.Username == "", item.Identity.FirstName, ""),
			LastName:  lo.Ternary(item.Profile.Username == "", item.Identity.LastName, ""),
			Username:  item.Profile.Username,
			Slug:      item.Profile.Slug,
			CreatedAt: item.CreatedAt,
		}
	}), total, nil
}
