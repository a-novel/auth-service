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
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetCredentials(t *testing.T) {
	data := []struct {
		name string

		tokenRaw string
		now      time.Time

		introspectToken    *models.UserTokenStatus
		introspectTokenErr error

		shouldCallCredentialsDAO bool
		credentialsDAO           *dao.CredentialsModel
		credentialsDAOErr        error

		expect    *models.Credentials
		expectErr error
	}{
		{
			name:     "Success",
			tokenRaw: "string-token",
			now:      baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallCredentialsDAO: true,
			credentialsDAO: &dao.CredentialsModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email: dao.Email{User: "user", Domain: "domain.com"},
				},
			},
			expect: &models.Credentials{
				Email:     "user@domain.com",
				Validated: true,
			},
		},
		{
			name:     "Success/WithEmailPendingValidation",
			tokenRaw: "string-token",
			now:      baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallCredentialsDAO: true,
			credentialsDAO: &dao.CredentialsModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    dao.Email{User: "user", Domain: "domain.com"},
					NewEmail: dao.Email{User: "new-user", Domain: "domain.com"},
				},
			},
			expect: &models.Credentials{
				Email:     "user@domain.com",
				NewEmail:  "new-user@domain.com",
				Validated: true,
			},
		},
		{
			name:     "Success/WithEmailNotValidated",
			tokenRaw: "string-token",
			now:      baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallCredentialsDAO: true,
			credentialsDAO: &dao.CredentialsModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email: dao.Email{User: "user", Domain: "domain.com", Validation: "validation-code"},
				},
			},
			expect: &models.Credentials{
				Email:     "user@domain.com",
				Validated: false,
			},
		},
		{
			name:     "Error/CredentialsDAOFailure",
			tokenRaw: "string-token",
			now:      baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallCredentialsDAO: true,
			credentialsDAOErr:        fooErr,
			expectErr:                fooErr,
		},
		{
			name:     "Error/InvalidToken",
			tokenRaw: "string-token",
			now:      baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: false,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			expectErr: goframework.ErrInvalidCredentials,
		},
		{
			name:               "Error/IntrospectTokenFailure",
			tokenRaw:           "string-token",
			now:                baseTime,
			introspectTokenErr: fooErr,
			expectErr:          fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			tokenService := servicesmocks.NewIntrospectTokenService(t)
			credentialsDAO := daomocks.NewCredentialsRepository(t)

			tokenService.
				On("IntrospectToken", context.Background(), d.tokenRaw, d.now, false).
				Return(d.introspectToken, d.introspectTokenErr)

			if d.shouldCallCredentialsDAO {
				credentialsDAO.
					On("GetCredentials", context.Background(), d.introspectToken.Token.Payload.ID).
					Return(d.credentialsDAO, d.credentialsDAOErr)
			}

			service := services.NewGetCredentialsService(credentialsDAO, tokenService)
			user, err := service.Get(context.Background(), d.tokenRaw, d.now)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, user)

			tokenService.AssertExpectations(t)
			credentialsDAO.AssertExpectations(t)
		})
	}
}
