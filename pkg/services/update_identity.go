package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	goframework "github.com/a-novel/go-framework"
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
		return goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidToken)
	}

	if err := goframework.CheckMinMax(form.FirstName, 1, MaxNameLength); err != nil {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidFirstName, err)
	}
	if err := goframework.CheckMinMax(form.LastName, 1, MaxNameLength); err != nil {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidLastName, err)
	}

	if err := goframework.CheckRestricted(form.Sex, models.SexMale, models.SexFemale); err != nil {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidSex, err)
	}
	if err := goframework.CheckRegexp(form.FirstName, nameRegexp); err != nil {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidFirstName, err)
	}
	if err := goframework.CheckRegexp(form.LastName, nameRegexp); err != nil {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidLastName, err)
	}

	age := getUserAge(form.Birthday, now)
	if err := goframework.CheckMinMax(age, MinAge, MaxAge); err != nil {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidAge, err)
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
