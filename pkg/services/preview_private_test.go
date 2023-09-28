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

func TestPreviewPrivate(t *testing.T) {
	data := []struct {
		name string

		tokenRaw string
		now      time.Time

		introspectToken    *models.UserTokenStatus
		introspectTokenErr error

		shouldCallCredentialsDAO bool
		credentialsDAO           *dao.CredentialsModel
		credentialsDAOErr        error

		shouldCallProfileDAO bool
		profileDAO           *dao.ProfileModel
		profileDAOErr        error

		shouldCallIdentityDAO bool
		identityDAO           *dao.IdentityModel
		identityDAOErr        error

		expect    *models.UserPreviewPrivate
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
			shouldCallProfileDAO: true,
			profileDAO: &dao.ProfileModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				ProfileModelCore: dao.ProfileModelCore{
					Slug: "slug",
				},
			},
			shouldCallIdentityDAO: true,
			identityDAO: &dao.IdentityModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name",
					LastName:  "last-name",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
			},
			expect: &models.UserPreviewPrivate{
				Email:     "user@domain.com",
				Validated: true,
				UserPreview: models.UserPreview{
					ID:        goframework.NumberUUID(1),
					FirstName: "name",
					LastName:  "last-name",
					Slug:      "slug",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name:     "Success/WithUsername",
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
			shouldCallProfileDAO: true,
			profileDAO: &dao.ProfileModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				ProfileModelCore: dao.ProfileModelCore{
					Username: "username",
					Slug:     "slug",
				},
			},
			shouldCallIdentityDAO: true,
			identityDAO: &dao.IdentityModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name",
					LastName:  "last-name",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
			},
			expect: &models.UserPreviewPrivate{
				Email:     "user@domain.com",
				Validated: true,
				UserPreview: models.UserPreview{
					ID:        goframework.NumberUUID(1),
					FirstName: "name",
					LastName:  "last-name",
					Username:  "username",
					Slug:      "slug",
					CreatedAt: baseTime,
				},
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
			shouldCallProfileDAO: true,
			profileDAO: &dao.ProfileModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				ProfileModelCore: dao.ProfileModelCore{
					Username: "username",
					Slug:     "slug",
				},
			},
			shouldCallIdentityDAO: true,
			identityDAO: &dao.IdentityModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name",
					LastName:  "last-name",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
			},
			expect: &models.UserPreviewPrivate{
				Email:     "user@domain.com",
				NewEmail:  "new-user@domain.com",
				Validated: true,
				UserPreview: models.UserPreview{
					ID:        goframework.NumberUUID(1),
					FirstName: "name",
					LastName:  "last-name",
					Username:  "username",
					Slug:      "slug",
					CreatedAt: baseTime,
				},
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
			shouldCallProfileDAO: true,
			profileDAO: &dao.ProfileModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				ProfileModelCore: dao.ProfileModelCore{
					Slug: "slug",
				},
			},
			shouldCallIdentityDAO: true,
			identityDAO: &dao.IdentityModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name",
					LastName:  "last-name",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
			},
			expect: &models.UserPreviewPrivate{
				Email:     "user@domain.com",
				Validated: false,
				UserPreview: models.UserPreview{
					ID:        goframework.NumberUUID(1),
					FirstName: "name",
					LastName:  "last-name",
					Slug:      "slug",
					CreatedAt: baseTime,
				},
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
			shouldCallCredentialsDAO: true,
			credentialsDAO: &dao.CredentialsModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email: dao.Email{User: "user", Domain: "domain.com"},
				},
			},
			shouldCallProfileDAO: true,
			profileDAO: &dao.ProfileModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				ProfileModelCore: dao.ProfileModelCore{
					Slug: "slug",
				},
			},
			shouldCallIdentityDAO: true,
			identityDAOErr:        fooErr,
			expectErr:             fooErr,
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
			shouldCallCredentialsDAO: true,
			credentialsDAO: &dao.CredentialsModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email: dao.Email{User: "user", Domain: "domain.com"},
				},
			},
			shouldCallProfileDAO: true,
			profileDAOErr:        fooErr,
			expectErr:            fooErr,
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
			profileDAO := daomocks.NewProfileRepository(t)
			identityDAO := daomocks.NewIdentityRepository(t)

			tokenService.
				On("IntrospectToken", context.Background(), d.tokenRaw, d.now, false).
				Return(d.introspectToken, d.introspectTokenErr)

			if d.shouldCallCredentialsDAO {
				credentialsDAO.
					On("GetCredentials", context.Background(), d.introspectToken.Token.Payload.ID).
					Return(d.credentialsDAO, d.credentialsDAOErr)
			}

			if d.shouldCallProfileDAO {
				profileDAO.
					On("GetProfile", context.Background(), d.introspectToken.Token.Payload.ID).
					Return(d.profileDAO, d.profileDAOErr)
			}

			if d.shouldCallIdentityDAO {
				identityDAO.
					On("GetIdentity", context.Background(), d.introspectToken.Token.Payload.ID).
					Return(d.identityDAO, d.identityDAOErr)
			}

			service := services.NewPreviewPrivateService(credentialsDAO, profileDAO, identityDAO, tokenService)
			user, err := service.Preview(context.Background(), d.tokenRaw, d.now)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, user)

			tokenService.AssertExpectations(t)
			credentialsDAO.AssertExpectations(t)
			profileDAO.AssertExpectations(t)
			identityDAO.AssertExpectations(t)
		})
	}
}
