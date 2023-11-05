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

type ValidateNewEmailService interface {
	ValidateNewEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) error
}

func NewValidateNewEmailService(
	credentialsDAO dao.CredentialsRepository,
	permissionsClient apiclients.PermissionsClient,
) ValidateNewEmailService {
	return &validateNewEmailServiceImpl{
		credentialsDAO:    credentialsDAO,
		permissionsClient: permissionsClient,
	}
}

type validateNewEmailServiceImpl struct {
	credentialsDAO    dao.CredentialsRepository
	permissionsClient apiclients.PermissionsClient
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

	err = s.credentialsDAO.RunInTx(ctx, func(ctx context.Context, txClient dao.CredentialsRepository) error {
		_, err = txClient.ValidateNewEmail(ctx, id, now)
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
