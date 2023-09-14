package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/services"
	goframework "github.com/a-novel/go-framework"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestValidateNewEmail(t *testing.T) {
	data := []struct {
		name string

		id   uuid.UUID
		code string
		now  time.Time

		dao    *dao.CredentialsModel
		daoErr error

		shouldCallUpdate bool
		updateErr        error

		expectErr error
	}{
		{
			name: "Success",
			id:   goframework.NumberUUID(1),
			code: publicValidationCode,
			now:  baseTime,
			dao: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    dao.Email{User: "user", Domain: "domain"},
					NewEmail: dao.Email{User: "new-user", Domain: "domain", Validation: privateValidationCode},
				},
			},
			shouldCallUpdate: true,
		},
		{
			name: "Error/NoPendingValidation",
			id:   goframework.NumberUUID(1),
			code: publicValidationCode,
			now:  baseTime,
			dao: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Email: dao.Email{User: "user", Domain: "domain"},
				},
			},
			expectErr: goframework.ErrInvalidCredentials,
		},
		{
			name: "Error/WrongValidationCode",
			id:   goframework.NumberUUID(1),
			code: "fakecode",
			now:  baseTime,
			dao: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    dao.Email{User: "user", Domain: "domain"},
					NewEmail: dao.Email{User: "new-user", Domain: "domain", Validation: privateValidationCode},
				},
			},
			expectErr: goframework.ErrInvalidCredentials,
		},
		{
			name: "Error/UpdateFailure",
			id:   goframework.NumberUUID(1),
			code: publicValidationCode,
			now:  baseTime,
			dao: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    dao.Email{User: "user", Domain: "domain"},
					NewEmail: dao.Email{User: "new-user", Domain: "domain", Validation: privateValidationCode},
				},
			},
			shouldCallUpdate: true,
			updateErr:        fooErr,
			expectErr:        fooErr,
		},
		{
			name:      "Error/GetCredentialsFailure",
			id:        goframework.NumberUUID(1),
			code:      publicValidationCode,
			now:       baseTime,
			daoErr:    fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsDAO := daomocks.NewCredentialsRepository(t)

			credentialsDAO.
				On("GetCredentials", context.Background(), d.id).
				Return(d.dao, d.daoErr)

			if d.shouldCallUpdate {
				credentialsDAO.
					On("ValidateNewEmail", context.Background(), d.id, d.now).
					Return(nil, d.updateErr)
			}

			service := services.NewValidateNewEmailService(credentialsDAO)
			err := service.ValidateNewEmail(context.Background(), d.id, d.code, d.now)

			require.ErrorIs(t, err, d.expectErr)

			credentialsDAO.AssertExpectations(t)
		})
	}
}
