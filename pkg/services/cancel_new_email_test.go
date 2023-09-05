package services_test

import (
	"context"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCancelNewEmail(t *testing.T) {
	data := []struct {
		name string

		tokenRaw string
		now      time.Time

		introspectTokenResp *models.UserTokenStatus
		introspectTokenErr  error

		shouldCallCredentialsDAO bool
		cancelNewEmailErr        error

		expectErr error
	}{
		{
			name:     "Success",
			tokenRaw: "string-token",
			now:      baseTime,
			introspectTokenResp: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
				TokenRaw: "string-token",
			},
			shouldCallCredentialsDAO: true,
		},
		{
			name:     "Error/DAOFailure",
			tokenRaw: "string-token",
			now:      baseTime,
			introspectTokenResp: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
				TokenRaw: "string-token",
			},
			shouldCallCredentialsDAO: true,
			cancelNewEmailErr:        fooErr,
			expectErr:                fooErr,
		},
		{
			name:               "Error/IntrospectTokenFailure",
			tokenRaw:           "string-token",
			now:                baseTime,
			introspectTokenErr: fooErr,
			expectErr:          fooErr,
		},
		{
			name:                "Error/InvalidToken",
			tokenRaw:            "string-token",
			now:                 baseTime,
			introspectTokenResp: &models.UserTokenStatus{OK: false},
			expectErr:           errors.ErrInvalidCredentials,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsDAO := daomocks.NewCredentialsRepository(t)
			introspectTokenService := servicesmocks.NewIntrospectTokenService(t)

			introspectTokenService.
				On("IntrospectToken", context.Background(), d.tokenRaw, d.now, false).
				Return(d.introspectTokenResp, d.introspectTokenErr)

			if d.shouldCallCredentialsDAO {
				credentialsDAO.
					On("CancelNewEmail", context.Background(), d.introspectTokenResp.Token.Payload.ID, d.now).
					Return(nil, d.cancelNewEmailErr)
			}

			service := services.NewCancelNewEmailService(credentialsDAO, introspectTokenService)
			err := service.CancelNewEmail(context.Background(), d.tokenRaw, d.now)

			require.ErrorIs(t, err, d.expectErr)

			credentialsDAO.AssertExpectations(t)
			introspectTokenService.AssertExpectations(t)
		})
	}
}
