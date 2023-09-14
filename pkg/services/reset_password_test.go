package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/bunovel"
	goframework "github.com/a-novel/go-framework"
	sendgridproxy "github.com/a-novel/sendgrid-proxy"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestResetPassword(t *testing.T) {
	data := []struct {
		name string

		email string
		now   time.Time

		passwordResetLink     string
		passwordResetTemplate string

		updatePasswordLink     string
		updatePasswordTemplate string

		publicValidationCode      string
		privateValidationCode     string
		generateValidationCodeErr error

		shouldCallCredentialsDAO bool
		credentialsDAO           *dao.CredentialsModel
		credentialsDAOErr        error

		shouldCallIdentityDAO bool
		identityDAO           *dao.IdentityModel
		identityDAOErr        error

		shouldCallMailer          bool
		shouldCallMailerWithEmail *mail.Email
		shouldCallMailerWithData  map[string]interface{}
		mailerErr                 error

		expectErr         error
		expectDeferred    bool
		expectDeferredErr error
	}{
		{
			name:                     "Success",
			email:                    "user@domain.com",
			now:                      baseTime,
			passwordResetLink:        "password-reset-link",
			passwordResetTemplate:    "password-reset-template",
			updatePasswordLink:       "update-password-link",
			updatePasswordTemplate:   "update-password-template",
			publicValidationCode:     "public-validation-code",
			privateValidationCode:    "private-validation-code",
			shouldCallCredentialsDAO: true,
			credentialsDAO: &dao.CredentialsModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email: dao.Email{User: "user", Domain: "domain.com"},
				},
			},
			shouldCallIdentityDAO: true,
			identityDAO: &dao.IdentityModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name",
				},
			},
			shouldCallMailer:          true,
			shouldCallMailerWithEmail: mail.NewEmail("name", "user@domain.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "name",
				"validation_link": "update-password-link?id=01010101-0101-0101-0101-010101010101&code=public-validation-code",
			},
			expectDeferred: true,
		},
		{
			name:                     "Error/SendingEmailFailure",
			email:                    "user@domain.com",
			now:                      baseTime,
			passwordResetLink:        "password-reset-link",
			passwordResetTemplate:    "password-reset-template",
			updatePasswordLink:       "update-password-link",
			updatePasswordTemplate:   "update-password-template",
			publicValidationCode:     "public-validation-code",
			privateValidationCode:    "private-validation-code",
			shouldCallCredentialsDAO: true,
			credentialsDAO: &dao.CredentialsModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email: dao.Email{User: "user", Domain: "domain.com"},
				},
			},
			shouldCallIdentityDAO: true,
			identityDAO: &dao.IdentityModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name",
				},
			},
			shouldCallMailer:          true,
			shouldCallMailerWithEmail: mail.NewEmail("name", "user@domain.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "name",
				"validation_link": "update-password-link?id=01010101-0101-0101-0101-010101010101&code=public-validation-code",
			},
			mailerErr:         fooErr,
			expectDeferred:    true,
			expectDeferredErr: fooErr,
		},
		{
			name:                     "Error/IdentityDAOFailure",
			email:                    "user@domain.com",
			now:                      baseTime,
			passwordResetLink:        "password-reset-link",
			passwordResetTemplate:    "password-reset-template",
			updatePasswordLink:       "update-password-link",
			updatePasswordTemplate:   "update-password-template",
			publicValidationCode:     "public-validation-code",
			privateValidationCode:    "private-validation-code",
			shouldCallCredentialsDAO: true,
			credentialsDAO: &dao.CredentialsModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email: dao.Email{User: "user", Domain: "domain.com"},
				},
			},
			shouldCallIdentityDAO: true,
			identityDAOErr:        fooErr,
			expectErr:             fooErr,
		},
		{
			name:                     "Error/CredentialsDAOFailure",
			email:                    "user@domain.com",
			now:                      baseTime,
			passwordResetLink:        "password-reset-link",
			passwordResetTemplate:    "password-reset-template",
			updatePasswordLink:       "update-password-link",
			updatePasswordTemplate:   "update-password-template",
			publicValidationCode:     "public-validation-code",
			privateValidationCode:    "private-validation-code",
			shouldCallCredentialsDAO: true,
			credentialsDAOErr:        fooErr,
			expectErr:                fooErr,
		},
		{
			name:                      "Error/GenerateValidationCodeFailure",
			email:                     "user@domain.com",
			now:                       baseTime,
			passwordResetLink:         "password-reset-link",
			passwordResetTemplate:     "password-reset-template",
			updatePasswordLink:        "update-password-link",
			updatePasswordTemplate:    "update-password-template",
			generateValidationCodeErr: fooErr,
			expectErr:                 fooErr,
		},
		{
			name:                  "Error/InvalidEmail",
			email:                 "userdomain.com",
			now:                   baseTime,
			passwordResetLink:     "password-reset-link",
			passwordResetTemplate: "password-reset-template",
			expectErr:             goframework.ErrInvalidEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsDAO := daomocks.NewCredentialsRepository(t)
			identityDAO := daomocks.NewIdentityRepository(t)
			mailerService := sendgridproxy.NewMockMailer(t)

			generateLink := func() (string, string, error) {
				return d.publicValidationCode, d.privateValidationCode, d.generateValidationCodeErr
			}

			if d.shouldCallCredentialsDAO {
				credentialsDAO.
					On("ResetPassword", context.Background(), d.privateValidationCode, mock.Anything, d.now).
					Return(d.credentialsDAO, d.credentialsDAOErr)
			}

			if d.shouldCallIdentityDAO {
				identityDAO.
					On("GetIdentity", context.Background(), d.credentialsDAO.ID).
					Return(d.identityDAO, d.identityDAOErr)
			}

			if d.shouldCallMailer {
				mailerService.
					On("Send", context.Background(), d.shouldCallMailerWithEmail, d.updatePasswordTemplate, d.shouldCallMailerWithData).
					Return(d.mailerErr)
			}

			service := services.NewResetPasswordService(credentialsDAO, identityDAO, mailerService, generateLink, d.updatePasswordLink, d.updatePasswordTemplate)
			deferred, err := service.ResetPassword(context.Background(), d.email, d.now)

			require.ErrorIs(t, err, d.expectErr)

			if d.expectDeferred {
				require.NotNil(t, deferred)
				require.ErrorIs(t, deferred(), d.expectDeferredErr)
			} else {
				require.Nil(t, deferred)
			}

			credentialsDAO.AssertExpectations(t)
			identityDAO.AssertExpectations(t)
			mailerService.AssertExpectations(t)
		})
	}
}
