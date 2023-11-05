package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/bunovel"
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

	// We don't use slugExist here, because the user may want to update other fields and keep its slug. To check if
	// slug is available, we must also validate it is taken by a different user than the one performing the update.
	profileWithSameSlug, err := s.profileDAO.GetProfileBySlug(ctx, form.Slug)
	if err != nil && !goerrors.Is(err, bunovel.ErrNotFound) {
		return goerrors.Join(ErrSlugExists, err)
	}
	if profileWithSameSlug != nil && profileWithSameSlug.ID != token.Token.Payload.ID {
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
