package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	goframework "github.com/a-novel/go-framework"
	"time"
)

type UpdateProfileService interface {
	UpdateProfile(ctx context.Context, tokenRaw string, now time.Time, form models.UpdateProfileForm) error
}

func NewUpdateProfileService(ProfileDAO dao.ProfileRepository, introspectTokenService IntrospectTokenService) UpdateProfileService {
	return &updateProfileServiceImpl{
		profileDAO:             ProfileDAO,
		IntrospectTokenService: introspectTokenService,
	}
}

type updateProfileServiceImpl struct {
	profileDAO dao.ProfileRepository
	IntrospectTokenService
}

func (s *updateProfileServiceImpl) UpdateProfile(ctx context.Context, tokenRaw string, now time.Time, form models.UpdateProfileForm) error {
	token, err := s.IntrospectToken(ctx, tokenRaw, now, false)
	if err != nil {
		return goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidToken)
	}

	if err := goframework.CheckMinMax(form.Slug, 1, MaxSlugLength); err != nil {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidSlug, err)
	}
	if err := goframework.CheckMinMax(form.Username, -1, MaxUsernameLength); err != nil {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidUsername, err)
	}

	if err := goframework.CheckRegexp(form.Slug, slugRegexp); err != nil {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidSlug, err)
	}
	if form.Username != "" {
		if err := goframework.CheckRegexp(form.Username, usernameRegexp); err != nil {
			return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidUsername, err)
		}
	}

	slugExists, err := s.profileDAO.SlugExists(ctx, form.Slug)
	if err != nil {
		return goerrors.Join(ErrSlugExists, err)
	}
	if slugExists {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidSlug, ErrTaken)
	}

	if _, err := s.profileDAO.Update(ctx, &dao.ProfileModelCore{
		Slug:     form.Slug,
		Username: form.Username,
	}, token.Token.Payload.ID, now); err != nil {
		return goerrors.Join(ErrUpdateProfile, err)
	}

	return nil
}
