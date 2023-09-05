package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
)

type SlugExistsService interface {
	SlugExists(ctx context.Context, slug string) (bool, error)
}

func NewSlugExistsService(
	profileDAO dao.ProfileRepository,
) SlugExistsService {
	return &slugExistsServiceImpl{
		profileDAO: profileDAO,
	}
}

type slugExistsServiceImpl struct {
	profileDAO dao.ProfileRepository
}

func (s *slugExistsServiceImpl) SlugExists(ctx context.Context, slug string) (bool, error) {
	ok, err := s.profileDAO.SlugExists(ctx, slug)
	if err != nil {
		return false, goerrors.Join(ErrSlugExists, err)
	}

	return ok, nil
}
