package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	goframework "github.com/a-novel/go-framework"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UpdatePasswordService interface {
	UpdatePassword(ctx context.Context, form models.UpdatePasswordForm, now time.Time) error
}

func NewUpdatePasswordService(credentialsDAO dao.CredentialsRepository) UpdatePasswordService {
	return &updatePasswordServiceImpl{
		credentialsDAO: credentialsDAO,
	}
}

type updatePasswordServiceImpl struct {
	credentialsDAO dao.CredentialsRepository
}

func (s *updatePasswordServiceImpl) UpdatePassword(ctx context.Context, form models.UpdatePasswordForm, now time.Time) error {
	if err := goframework.CheckMinMax(form.NewPassword, MinPasswordLength, MaxPasswordLength); err != nil {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidPassword, err)
	}
	if form.Code == "" && form.OldPassword == "" {
		return goerrors.Join(goframework.ErrInvalidEntity, ErrMissingPasswordValidation)
	}

	credentials, err := s.credentialsDAO.GetCredentials(ctx, form.ID)
	if err != nil {
		return goerrors.Join(ErrGetCredentials, err)
	}

	if form.Code != "" {
		if credentials.Password.Validation == "" {
			return goerrors.Join(goframework.ErrInvalidCredentials, ErrMissingPendingValidation)
		}

		ok, err := goframework.VerifyCode(form.Code, credentials.Password.Validation)
		if err != nil {
			return goerrors.Join(ErrVerifyValidationCode, err)
		}
		if !ok {
			return goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidValidationCode)
		}
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(credentials.Password.Hashed), []byte(form.OldPassword))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				return goerrors.Join(goframework.ErrInvalidCredentials, ErrWrongPassword)
			}

			return goerrors.Join(ErrCheckPassword, err)
		}
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(form.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return goerrors.Join(ErrHashPassword, err)
	}

	if _, err := s.credentialsDAO.UpdatePassword(ctx, string(passwordHashed), form.ID, now); err != nil {
		return goerrors.Join(ErrUpdatePassword, err)
	}

	return nil
}
