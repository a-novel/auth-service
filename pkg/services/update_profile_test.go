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
	"strings"
	"testing"
	"time"
)

func TestUpdateProfile(t *testing.T) {
	data := []struct {
		name string

		tokenRaw string
		now      time.Time
		form     models.UpdateProfileForm

		introspectToken    *models.UserTokenStatus
		introspectTokenErr error

		shouldCallSlugExists bool
		slugExists           *dao.ProfileModel
		slugExistsErr        error

		shouldCallDAO bool
		daoErr        error

		expectErr error
	}{
		{
			name:     "Success",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Slug: "slug",
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallSlugExists: true,
			slugExists:           nil,
			shouldCallDAO:        true,
		},
		{
			name:     "Success/WithUsername",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Username: "username",
				Slug:     "slug",
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallSlugExists: true,
			slugExists:           nil,
			shouldCallDAO:        true,
		},
		{
			name:     "Error/UsernameInvalid",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Username: "😊😊😊",
				Slug:     "slug",
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			expectErr: goframework.ErrInvalidEntity,
		},
		{
			name:     "Error/UsernameTooLong",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Username: strings.Repeat("a", services.MaxUsernameLength+1),
				Slug:     "slug",
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			expectErr: goframework.ErrInvalidEntity,
		},
		{
			name:     "Error/InvalidSlug",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Slug: "Sl#ug",
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			expectErr: goframework.ErrInvalidEntity,
		},
		{
			name:     "Error/SlugTooLong",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Slug: strings.Repeat("a", services.MaxSlugLength+1),
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			expectErr: goframework.ErrInvalidEntity,
		},
		{
			name:     "Error/NoSlug",
			tokenRaw: "string-token",
			now:      baseTime,
			form:     models.UpdateProfileForm{},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			expectErr: goframework.ErrInvalidEntity,
		},
		{
			name:     "Error/UpdateFailure",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Slug: "slug",
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallSlugExists: true,
			slugExists:           nil,
			shouldCallDAO:        true,
			daoErr:               fooErr,
			expectErr:            fooErr,
		},
		{
			name:     "Error/SlugExists",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Slug: "slug",
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallSlugExists: true,
			slugExists: &dao.ProfileModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime, &baseTime),
			},
			expectErr: services.ErrTaken,
		},
		{
			name:     "Success/SlugExists/ButSameUser",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Slug: "slug",
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallSlugExists: true,
			slugExists: &dao.ProfileModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
			},
			shouldCallDAO: true,
		},
		{
			name:     "Error/SlugCheckFailure",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Slug: "slug",
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallSlugExists: true,
			slugExistsErr:        fooErr,
			expectErr:            fooErr,
		},
		{
			name:     "Error/TokenInvalid",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Slug: "slug",
			},
			introspectToken: &models.UserTokenStatus{
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			expectErr: goframework.ErrInvalidCredentials,
		},
		{
			name:     "Error/IntrospectTokenFailure",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateProfileForm{
				Slug: "slug",
			},
			introspectTokenErr: fooErr,
			expectErr:          fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			profileDAO := daomocks.NewProfileRepository(t)
			introspectTokenService := servicesmocks.NewIntrospectTokenService(t)

			introspectTokenService.
				On("IntrospectToken", context.Background(), d.tokenRaw, d.now, false).
				Return(d.introspectToken, d.introspectTokenErr)

			if d.shouldCallSlugExists {
				profileDAO.
					On("GetProfileBySlug", context.Background(), d.form.Slug).
					Return(d.slugExists, d.slugExistsErr)
			}

			if d.shouldCallDAO {
				profileDAO.
					On("Update", context.Background(), &dao.ProfileModelCore{
						Username: d.form.Username,
						Slug:     d.form.Slug,
					}, d.introspectToken.Token.Payload.ID, d.now).
					Return(nil, d.daoErr)
			}

			service := services.NewUpdateProfileService(profileDAO, introspectTokenService)
			err := service.UpdateProfile(context.Background(), d.tokenRaw, d.now, d.form)

			require.ErrorIs(t, err, d.expectErr)

			profileDAO.AssertExpectations(t)
		})
	}
}
