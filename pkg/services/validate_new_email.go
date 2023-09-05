package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/security"
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
		return errors.ErrInvalidCredentials
	}
	ok, err := security.VerifyCode(code, credentials.NewEmail.Validation)
	if err != nil {
		return goerrors.Join(ErrVerifyValidationCode, err)
	}
	if !ok {
		return goerrors.Join(errors.ErrInvalidCredentials, ErrInvalidValidationCode)
	}

	_, err = s.credentialsDAO.ValidateNewEmail(ctx, id, now)
	if err != nil {
		return goerrors.Join(ErrValidateEmail, err)
	}

	return nil
}
