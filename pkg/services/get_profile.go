package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	goframework "github.com/a-novel/go-framework"
	"time"
)

type GetProfileService interface {
	Get(ctx context.Context, tokenRaw string, now time.Time) (*models.Profile, error)
}

func NewGetProfileService(
	profileDAO dao.ProfileRepository,
	introspectTokenService IntrospectTokenService,
) GetProfileService {
	return &getProfileServiceImpl{
		profileDAO:             profileDAO,
		IntrospectTokenService: introspectTokenService,
	}
}

type getProfileServiceImpl struct {
	profileDAO dao.ProfileRepository
	IntrospectTokenService
}

func (s *getProfileServiceImpl) Get(ctx context.Context, tokenRaw string, now time.Time) (*models.Profile, error) {
	token, err := s.IntrospectToken(ctx, tokenRaw, now, false)
	if err != nil {
		return nil, goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return nil, goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidToken)
	}

	profile, err := s.profileDAO.GetProfile(ctx, token.Token.Payload.ID)
	if err != nil {
		return nil, goerrors.Join(ErrGetProfile, err)
	}

	return &models.Profile{
		Username: profile.Username,
		Slug:     profile.Slug,
	}, nil
}
