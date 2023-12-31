package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/bunovel"
	goframework "github.com/a-novel/go-framework"
	sendgridproxy "github.com/a-novel/sendgrid-proxy"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestResendEmailValidation(t *testing.T) {
	data := []struct {
		name string

		tokenRaw string
		now      time.Time

		validateEmailLink     string
		validateEmailTemplate string

		introspectToken    *models.UserTokenStatus
		introspectTokenErr error

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
			name:                  "Success",
			tokenRaw:              "string-token",
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
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
				"validation_link": "validate-email-link?id=01010101-0101-0101-0101-010101010101&code=public-validation-code",
			},
			expectDeferred: true,
		},
		{
			name:                  "Error/SendingEmailFailure",
			tokenRaw:              "string-token",
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
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
				"validation_link": "validate-email-link?id=01010101-0101-0101-0101-010101010101&code=public-validation-code",
			},
			mailerErr:         fooErr,
			expectDeferred:    true,
			expectDeferredErr: fooErr,
		},
		{
			name:                  "Error/IdentityDAOFailure",
			tokenRaw:              "string-token",
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
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
			name:                  "Error/CredentialsDAOFailure",
			tokenRaw:              "string-token",
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			publicValidationCode:     "public-validation-code",
			privateValidationCode:    "private-validation-code",
			shouldCallCredentialsDAO: true,
			credentialsDAOErr:        fooErr,
			expectErr:                fooErr,
		},
		{
			name:                  "Error/GenerateValidationCodeFailure",
			tokenRaw:              "string-token",
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			generateValidationCodeErr: fooErr,
			expectErr:                 fooErr,
		},
		{
			name:                  "Error/InvalidToken",
			tokenRaw:              "string-token",
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			introspectToken: &models.UserTokenStatus{
				OK: false,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			expectErr: goframework.ErrInvalidCredentials,
		},
		{
			name:                  "Error/IntrospectTokenFailure",
			tokenRaw:              "string-token",
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			introspectTokenErr:    fooErr,
			expectErr:             fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsDAO := daomocks.NewCredentialsRepository(t)
			identityDAO := daomocks.NewIdentityRepository(t)
			mailerService := sendgridproxy.NewMockMailer(t)
			introspectTokenService := servicesmocks.NewIntrospectTokenService(t)

			generateLink := func() (string, string, error) {
				return d.publicValidationCode, d.privateValidationCode, d.generateValidationCodeErr
			}

			introspectTokenService.
				On("IntrospectToken", context.Background(), d.tokenRaw, d.now, false).
				Return(d.introspectToken, d.introspectTokenErr)

			if d.shouldCallCredentialsDAO {
				credentialsDAO.
					On("UpdateEmailValidation", context.Background(), d.privateValidationCode, d.introspectToken.Token.Payload.ID, d.now).
					Return(d.credentialsDAO, d.credentialsDAOErr)
			}

			if d.shouldCallIdentityDAO {
				identityDAO.
					On("GetIdentity", context.Background(), d.introspectToken.Token.Payload.ID).
					Return(d.identityDAO, d.identityDAOErr)
			}

			if d.shouldCallMailer {
				mailerService.
					On("Send", context.Background(), d.shouldCallMailerWithEmail, d.validateEmailTemplate, d.shouldCallMailerWithData).
					Return(d.mailerErr)
			}

			service := services.NewResendEmailValidationService(credentialsDAO, identityDAO, mailerService, generateLink, introspectTokenService, d.validateEmailLink, d.validateEmailTemplate)
			deferred, err := service.ResendEmailValidation(context.Background(), d.tokenRaw, d.now)

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
