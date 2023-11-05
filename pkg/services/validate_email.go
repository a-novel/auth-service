package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	apiclients "github.com/a-novel/go-apis/clients"
	goframework "github.com/a-novel/go-framework"
	"github.com/google/uuid"
	"time"
)

type ValidateEmailService interface {
	ValidateEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) error
}

func NewValidateEmailService(
	credentialsDAO dao.CredentialsRepository,
	permissionsClient apiclients.PermissionsClient,
) ValidateEmailService {
	return &validateEmailServiceImpl{
		credentialsDAO:    credentialsDAO,
		permissionsClient: permissionsClient,
	}
}

type validateEmailServiceImpl struct {
	credentialsDAO    dao.CredentialsRepository
	permissionsClient apiclients.PermissionsClient
}

func (s *validateEmailServiceImpl) ValidateEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) error {
	credentials, err := s.credentialsDAO.GetCredentials(ctx, id)
	if err != nil {
		return goerrors.Join(ErrGetCredentials, err)
	}

	// Email already validated.
	if credentials.Email.Validation == "" {
		return goerrors.Join(goframework.ErrInvalidCredentials, ErrMissingPendingValidation)
	}
	ok, err := goframework.VerifyCode(code, credentials.Email.Validation)
	if err != nil {
		return goerrors.Join(ErrVerifyValidationCode, err)
	}
	if !ok {
		return goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidValidationCode)
	}

	err = s.credentialsDAO.RunInTx(ctx, func(ctx context.Context, txClient dao.CredentialsRepository) error {
		_, err = txClient.ValidateEmail(ctx, id, now)
		if err != nil {
			return goerrors.Join(ErrValidateEmail, err)
		}

		err = s.permissionsClient.SetUserPermissions(ctx, apiclients.SetUserPermissionsForm{
			UserID:    id,
			SetFields: []string{apiclients.FieldValidatedAccount},
		})
		if err != nil {
			return goerrors.Join(ErrUpdateUserPermissions, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
