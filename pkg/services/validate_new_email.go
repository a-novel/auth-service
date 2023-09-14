package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	goframework "github.com/a-novel/go-framework"
	"github.com/google/uuid"
	"time"
)

type ValidateNewEmailService interface {
	ValidateNewEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) error
}

func NewValidateNewEmailService(
	credentialsDAO dao.CredentialsRepository,
) ValidateNewEmailService {
	return &validateNewEmailServiceImpl{
		credentialsDAO: credentialsDAO,
	}
}

type validateNewEmailServiceImpl struct {
	credentialsDAO dao.CredentialsRepository
}

func (s *validateNewEmailServiceImpl) ValidateNewEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) error {
	credentials, err := s.credentialsDAO.GetCredentials(ctx, id)
	if err != nil {
		return goerrors.Join(ErrGetCredentials, err)
	}

	// Email already validated.
	if credentials.NewEmail.Validation == "" {
		return goframework.ErrInvalidCredentials
	}
	ok, err := goframework.VerifyCode(code, credentials.NewEmail.Validation)
	if err != nil {
		return goerrors.Join(ErrVerifyValidationCode, err)
	}
	if !ok {
		return goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidValidationCode)
	}

	_, err = s.credentialsDAO.ValidateNewEmail(ctx, id, now)
	if err != nil {
		return goerrors.Join(ErrValidateEmail, err)
	}

	return nil
}
