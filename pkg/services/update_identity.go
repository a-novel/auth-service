package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/validation"
	"time"
)

type UpdateIdentityService interface {
	UpdateIdentity(ctx context.Context, tokenRaw string, now time.Time, form models.UpdateIdentityForm) error
}

func NewUpdateIdentityService(identityDAO dao.IdentityRepository, introspectTokenService IntrospectTokenService) UpdateIdentityService {
	return &updateIdentityServiceImpl{
		identityDAO:            identityDAO,
		IntrospectTokenService: introspectTokenService,
	}
}

type updateIdentityServiceImpl struct {
	identityDAO dao.IdentityRepository
	IntrospectTokenService
}

func (s *updateIdentityServiceImpl) UpdateIdentity(ctx context.Context, tokenRaw string, now time.Time, form models.UpdateIdentityForm) error {
	token, err := s.IntrospectToken(ctx, tokenRaw, now, false)
	if err != nil {
		return goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return goerrors.Join(errors.ErrInvalidCredentials, ErrInvalidToken)
	}

	if err := validation.CheckMinMax(form.FirstName, 1, MaxNameLength); err != nil {
		return goerrors.Join(errors.ErrInvalidEntity, ErrInvalidFirstName, err)
	}
	if err := validation.CheckMinMax(form.LastName, 1, MaxNameLength); err != nil {
		return goerrors.Join(errors.ErrInvalidEntity, ErrInvalidLastName, err)
	}

	if err := validation.CheckRestricted(form.Sex, models.SexMale, models.SexFemale); err != nil {
		return goerrors.Join(errors.ErrInvalidEntity, ErrInvalidSex, err)
	}
	if err := validation.CheckRegexp(form.FirstName, nameRegexp); err != nil {
		return goerrors.Join(errors.ErrInvalidEntity, ErrInvalidFirstName, err)
	}
	if err := validation.CheckRegexp(form.LastName, nameRegexp); err != nil {
		return goerrors.Join(errors.ErrInvalidEntity, ErrInvalidLastName, err)
	}

	age := getUserAge(form.Birthday, now)
	if err := validation.CheckMinMax(age, MinAge, MaxAge); err != nil {
		return goerrors.Join(errors.ErrInvalidEntity, ErrInvalidAge, err)
	}

	if _, err := s.identityDAO.Update(ctx, &dao.IdentityModelCore{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Birthday:  form.Birthday,
		Sex:       form.Sex,
	}, token.Token.Payload.ID, now); err != nil {
		return goerrors.Join(ErrUpdateIdentity, err)
	}

	return nil
}
