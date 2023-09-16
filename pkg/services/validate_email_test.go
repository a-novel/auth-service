package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/services"
	apiclients "github.com/a-novel/go-api-clients"
	apiclientsmocks "github.com/a-novel/go-api-clients/mocks"
	goframework "github.com/a-novel/go-framework"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestValidateEmail(t *testing.T) {
	data := []struct {
		name string

		id   uuid.UUID
		code string
		now  time.Time

		dao    *dao.CredentialsModel
		daoErr error

		shouldCallPermissionsClient bool
		permissionsClientErr        error

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
					Email: dao.Email{User: "user", Domain: "domain", Validation: privateValidationCode},
				},
			},
			shouldCallUpdate:            true,
			shouldCallPermissionsClient: true,
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
					Email: dao.Email{User: "user", Domain: "domain", Validation: privateValidationCode},
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
					Email: dao.Email{User: "user", Domain: "domain", Validation: privateValidationCode},
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
		{
			name: "Error/UpdatePermissionsFailure",
			id:   goframework.NumberUUID(1),
			code: publicValidationCode,
			now:  baseTime,
			dao: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Email: dao.Email{User: "user", Domain: "domain", Validation: privateValidationCode},
				},
			},
			shouldCallUpdate:            true,
			shouldCallPermissionsClient: true,
			permissionsClientErr:        fooErr,
			expectErr:                   fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsDAO := daomocks.NewCredentialsRepository(t)
			permissionsClient := apiclientsmocks.NewPermissionsClient(t)

			credentialsDAO.
				On("GetCredentials", context.Background(), d.id).
				Return(d.dao, d.daoErr)

			if d.shouldCallUpdate {
				credentialsDAO.
					On("ValidateEmail", context.Background(), d.id, d.now).
					Return(nil, d.updateErr)

				// Execute the actual method, but call the mocks inside of it.
				txCall := credentialsDAO.On("RunInTx", context.Background(), mock.Anything)
				txCall.Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context, dao.CredentialsRepository) error)
					txCall.ReturnArguments = []interface{}{fn(context.Background(), credentialsDAO)}
				})
			}

			if d.shouldCallPermissionsClient {
				permissionsClient.
					On("SetUserPermissions", context.Background(), apiclients.SetUserPermissionsForm{
						UserID:    d.id,
						SetFields: []string{apiclients.FieldValidatedAccount},
					}).
					Return(d.permissionsClientErr)
			}

			service := services.NewValidateEmailService(credentialsDAO, permissionsClient)
			err := service.ValidateEmail(context.Background(), d.id, d.code, d.now)

			require.ErrorIs(t, err, d.expectErr)

			credentialsDAO.AssertExpectations(t)
			permissionsClient.AssertExpectations(t)
		})
	}
}
