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

type ValidateEmailService interface {
	ValidateEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) error
}

func NewValidateEmailService(
	credentialsDAO dao.CredentialsRepository,
) ValidateEmailService {
	return &validateEmailServiceImpl{
		credentialsDAO: credentialsDAO,
	}
}

type validateEmailServiceImpl struct {
	credentialsDAO dao.CredentialsRepository
}

func (s *validateEmailServiceImpl) ValidateEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) error {
	credentials, err := s.credentialsDAO.GetCredentials(ctx, id)
	if err != nil {
		return goerrors.Join(ErrGetCredentials, err)
	}

	// Email already validated.
	if credentials.Email.Validation == "" {
		return goerrors.Join(errors.ErrInvalidCredentials, ErrMissingPendingValidation)
	}
	ok, err := security.VerifyCode(code, credentials.Email.Validation)
	if err != nil {
		return goerrors.Join(ErrVerifyValidationCode, err)
	}
	if !ok {
		return goerrors.Join(errors.ErrInvalidCredentials, ErrInvalidValidationCode)
	}

	_, err = s.credentialsDAO.ValidateEmail(ctx, id, now)
	if err != nil {
		return goerrors.Join(ErrValidateEmail, err)
	}

	return nil
}
