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

func TestGetProfile(t *testing.T) {
	data := []struct {
		name string

		tokenRaw string
		now      time.Time

		introspectToken    *models.UserTokenStatus
		introspectTokenErr error

		shouldCallProfileDAO bool
		profileDAO           *dao.ProfileModel
		profileDAOErr        error

		expect    *models.Profile
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
			shouldCallProfileDAO: true,
			profileDAO: &dao.ProfileModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				ProfileModelCore: dao.ProfileModelCore{
					Username: "username-1",
					Slug:     "slug-1",
				},
			},
			expect: &models.Profile{
				Username: "username-1",
				Slug:     "slug-1",
			},
		},
		{
			name:     "Error/ProfileDAOFailure",
			tokenRaw: "string-token",
			now:      baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallProfileDAO: true,
			profileDAOErr:        fooErr,
			expectErr:            fooErr,
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
			profileDAO := daomocks.NewProfileRepository(t)

			tokenService.
				On("IntrospectToken", context.Background(), d.tokenRaw, d.now, false).
				Return(d.introspectToken, d.introspectTokenErr)

			if d.shouldCallProfileDAO {
				profileDAO.
					On("GetProfile", context.Background(), d.introspectToken.Token.Payload.ID).
					Return(d.profileDAO, d.profileDAOErr)
			}

			service := services.NewGetProfileService(profileDAO, tokenService)
			user, err := service.Get(context.Background(), d.tokenRaw, d.now)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, user)

			tokenService.AssertExpectations(t)
			profileDAO.AssertExpectations(t)
		})
	}
}
