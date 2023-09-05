package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/mailer"
	"github.com/a-novel/go-framework/test"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUpdateEmail(t *testing.T) {
	data := []struct {
		name string

		validateEmailLink     string
		validateEmailTemplate string

		tokenRaw string
		newEmail string
		now      time.Time

		introspectToken    *models.UserTokenStatus
		introspectTokenErr error

		shouldCallEmailExists bool
		emailExists           bool
		emailExistsErr        error

		publicValidationCode      string
		privateValidationCode     string
		generateValidationCodeErr error

		shouldCallUpdateEmail bool
		updateEmailErr        error

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
			name:                  "Success",
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			tokenRaw:              "string-token",
			newEmail:              "new-user@domain.com",
			now:                   baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallEmailExists: true,
			emailExists:           false,
			publicValidationCode:  "public-validation-code",
			privateValidationCode: "private-validation-code",
			shouldCallUpdateEmail: true,
			shouldCallIdentityDAO: true,
			identityDAO: &dao.IdentityModel{
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name",
				},
			},
			shouldCallMailer:          true,
			shouldCallMailerWithEmail: mail.NewEmail("name", "new-user@domain.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "name",
				"validation_link": "validate-email-link?id=01010101-0101-0101-0101-010101010101&code=public-validation-code",
			},
			expectDeferred: true,
		},
		{
			name:                  "Error/SendingEmailFailure",
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			tokenRaw:              "string-token",
			newEmail:              "new-user@domain.com",
			now:                   baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallEmailExists: true,
			emailExists:           false,
			publicValidationCode:  "public-validation-code",
			privateValidationCode: "private-validation-code",
			shouldCallUpdateEmail: true,
			shouldCallIdentityDAO: true,
			identityDAO: &dao.IdentityModel{
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name",
				},
			},
			shouldCallMailer:          true,
			shouldCallMailerWithEmail: mail.NewEmail("name", "new-user@domain.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "name",
				"validation_link": "validate-email-link?id=01010101-0101-0101-0101-010101010101&code=public-validation-code",
			},
			mailerErr:         fooErr,
			expectDeferred:    true,
			expectDeferredErr: fooErr,
		},
		{
			name:                  "Error/IdentityDAOFailure",
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			tokenRaw:              "string-token",
			newEmail:              "new-user@domain.com",
			now:                   baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallEmailExists: true,
			emailExists:           false,
			publicValidationCode:  "public-validation-code",
			privateValidationCode: "private-validation-code",
			shouldCallUpdateEmail: true,
			shouldCallIdentityDAO: true,
			identityDAOErr:        fooErr,
			expectErr:             fooErr,
		},
		{
			name:                  "Error/UpdateEmailFailure",
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			tokenRaw:              "string-token",
			newEmail:              "new-user@domain.com",
			now:                   baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallEmailExists: true,
			emailExists:           false,
			publicValidationCode:  "public-validation-code",
			privateValidationCode: "private-validation-code",
			shouldCallUpdateEmail: true,
			updateEmailErr:        fooErr,
			expectErr:             fooErr,
		},
		{
			name:                  "Error/GenerateValidationCodeFailure",
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			tokenRaw:              "string-token",
			newEmail:              "new-user@domain.com",
			now:                   baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallEmailExists:     true,
			emailExists:               false,
			generateValidationCodeErr: fooErr,
			expectErr:                 fooErr,
		},
		{
			name:                  "Error/EmailAlreadyTaken",
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			tokenRaw:              "string-token",
			newEmail:              "new-user@domain.com",
			now:                   baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallEmailExists: true,
			emailExists:           true,
			expectErr:             services.ErrTaken,
		},
		{
			name:                  "Error/EmailExistsFailure",
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			tokenRaw:              "string-token",
			newEmail:              "new-user@domain.com",
			now:                   baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallEmailExists: true,
			emailExistsErr:        fooErr,
			expectErr:             fooErr,
		},
		{
			name:                  "Error/InvalidEmail",
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			tokenRaw:              "string-token",
			newEmail:              "new-userdomain.com",
			now:                   baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			expectErr: errors.ErrInvalidEntity,
		},
		{
			name:                  "Error/NoEmail",
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			tokenRaw:              "string-token",
			now:                   baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			expectErr: errors.ErrInvalidEntity,
		},
		{
			name:                  "Error/InvalidToken",
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			tokenRaw:              "string-token",
			newEmail:              "new-userdomain.com",
			now:                   baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: false,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			expectErr: errors.ErrInvalidCredentials,
		},
		{
			name:                  "Error/IntrospectTokenFailure",
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			tokenRaw:              "string-token",
			newEmail:              "new-userdomain.com",
			now:                   baseTime,
			introspectTokenErr:    fooErr,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsDAO := daomocks.NewCredentialsRepository(t)
			identityDAO := daomocks.NewIdentityRepository(t)
			mailerService := mailer.NewMockMailer(t)
			introspectTokenService := servicesmocks.NewIntrospectTokenService(t)

			generateLink := func() (string, string, error) {
				return d.publicValidationCode, d.privateValidationCode, d.generateValidationCodeErr
			}

			introspectTokenService.
				On("IntrospectToken", context.Background(), d.tokenRaw, d.now, false).
				Return(d.introspectToken, d.introspectTokenErr)

			if d.shouldCallEmailExists {
				credentialsDAO.
					On("EmailExists", context.Background(), mock.Anything).
					Return(d.emailExists, d.emailExistsErr)
			}

			if d.shouldCallUpdateEmail {
				credentialsDAO.
					On("UpdateEmail", context.Background(), mock.Anything, d.privateValidationCode, d.introspectToken.Token.Payload.ID, d.now).
					Return(nil, d.updateEmailErr)
			}

			if d.shouldCallIdentityDAO {
				identityDAO.
					On("GetIdentity", context.Background(), d.introspectToken.Token.Payload.ID).
					Return(d.identityDAO, d.identityDAOErr)
			}

			if d.shouldCallMailer {
				mailerService.
					On("Send", d.shouldCallMailerWithEmail, d.validateEmailTemplate, d.shouldCallMailerWithData).
					Return(d.mailerErr)
			}

			service := services.NewUpdateEmailService(
				credentialsDAO,
				identityDAO,
				mailerService,
				generateLink,
				introspectTokenService,
				d.validateEmailLink,
				d.validateEmailTemplate,
			)
			deferred, err := service.UpdateEmail(context.Background(), d.tokenRaw, d.newEmail, d.now)

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
			introspectTokenService.AssertExpectations(t)
		})
	}
}
