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

		shouldCallAuthorizationClient bool
		authorizationClientErr        error

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
			shouldCallUpdate:              true,
			shouldCallAuthorizationClient: true,
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
			name: "Error/updateAuthorizationsFailure",
			id:   goframework.NumberUUID(1),
			code: publicValidationCode,
			now:  baseTime,
			dao: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Email: dao.Email{User: "user", Domain: "domain", Validation: privateValidationCode},
				},
			},
			shouldCallUpdate:              true,
			shouldCallAuthorizationClient: true,
			authorizationClientErr:        fooErr,
			expectErr:                     fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsDAO := daomocks.NewCredentialsRepository(t)
			authorizationsClient := apiclientsmocks.NewAuthorizationsClient(t)

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

			if d.shouldCallAuthorizationClient {
				authorizationsClient.
					On("SetUserAuthorizations", context.Background(), apiclients.SetUserAuthorizationsForm{
						UserID:    d.id,
						SetFields: []string{"validated_account"},
					}).
					Return(d.authorizationClientErr)
			}

			service := services.NewValidateEmailService(credentialsDAO, authorizationsClient)
			err := service.ValidateEmail(context.Background(), d.id, d.code, d.now)

			require.ErrorIs(t, err, d.expectErr)

			credentialsDAO.AssertExpectations(t)
			authorizationsClient.AssertExpectations(t)
		})
	}
}
