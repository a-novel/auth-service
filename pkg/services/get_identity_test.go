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

func TestGetIdentity(t *testing.T) {
	data := []struct {
		name string

		tokenRaw string
		now      time.Time

		introspectToken    *models.UserTokenStatus
		introspectTokenErr error

		shouldCallIdentityDAO bool
		identityDAO           *dao.IdentityModel
		identityDAOErr        error

		expect    *models.Identity
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
			shouldCallIdentityDAO: true,
			identityDAO: &dao.IdentityModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name-1",
					LastName:  "last-name-1",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
			},
			expect: &models.Identity{
				FirstName: "name-1",
				LastName:  "last-name-1",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
		{
			name:     "Error/IdentityDAOFailure",
			tokenRaw: "string-token",
			now:      baseTime,
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallIdentityDAO: true,
			identityDAOErr:        fooErr,
			expectErr:             fooErr,
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
			identityDAO := daomocks.NewIdentityRepository(t)

			tokenService.
				On("IntrospectToken", context.Background(), d.tokenRaw, d.now, false).
				Return(d.introspectToken, d.introspectTokenErr)

			if d.shouldCallIdentityDAO {
				identityDAO.
					On("GetIdentity", context.Background(), d.introspectToken.Token.Payload.ID).
					Return(d.identityDAO, d.identityDAOErr)
			}

			service := services.NewGetIdentityService(identityDAO, tokenService)
			user, err := service.Get(context.Background(), d.tokenRaw, d.now)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, user)

			tokenService.AssertExpectations(t)
			identityDAO.AssertExpectations(t)
		})
	}
}
