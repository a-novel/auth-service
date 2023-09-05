package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/samber/lo"
)

type PreviewService interface {
	// Preview returns a subset of data for the requested user.
	Preview(ctx context.Context, slug string) (*models.UserPreview, error)
}

func NewPreviewService(profileDAO dao.ProfileRepository, identityDAO dao.IdentityRepository) PreviewService {
	return &previewServiceImpl{
		profileDAO:  profileDAO,
		identityDAO: identityDAO,
	}
}

type previewServiceImpl struct {
	profileDAO  dao.ProfileRepository
	identityDAO dao.IdentityRepository
}

func (s *previewServiceImpl) Preview(ctx context.Context, slug string) (*models.UserPreview, error) {
	profile, err := s.profileDAO.GetProfileBySlug(ctx, slug)
	if err != nil {
		return nil, goerrors.Join(ErrGetProfileBySlug, err)
	}

	identity, err := s.identityDAO.GetIdentity(ctx, profile.ID)
	if err != nil {
		return nil, goerrors.Join(ErrGetIdentity, err)
	}

	return &models.UserPreview{
		// Don't return user real name if a username is set.
		FirstName: lo.Ternary(profile.Username == "", identity.FirstName, ""),
		LastName:  lo.Ternary(profile.Username == "", identity.LastName, ""),
		Username:  profile.Username,
		Slug:      profile.Slug,
		CreatedAt: profile.CreatedAt,
	}, nil
}
