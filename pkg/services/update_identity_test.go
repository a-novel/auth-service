package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	goframework "github.com/a-novel/go-framework"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestUpdateIdentity(t *testing.T) {
	data := []struct {
		name string

		tokenRaw string
		now      time.Time
		form     models.UpdateIdentityForm

		introspectToken    *models.UserTokenStatus
		introspectTokenErr error

		shouldCallDAO bool
		daoErr        error

		expectErr error
	}{
		{
			name:     "Success",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallDAO: true,
		},
		{
			name:     "Error/DAOFailure",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			introspectToken: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			shouldCallDAO: true,
			daoErr:        fooErr,
			expectErr:     fooErr,
		},
		{
			name:     "Error/UserTooYoung",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add((services.MinAge - 1) * timeYear),
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
			name:     "Error/UserNotBorn(Seriously)",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(2 * timeYear),
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
			name:     "Error/UserTooOld",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add((services.MaxAge + 1) * timeYear),
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
			name:     "Error/InvalidLastName",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: "name",
				LastName:  "$%&/()=?",
				Sex:       models.SexMale,
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
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
			name:     "Error/LastNameTooLong",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: "name",
				LastName:  strings.Repeat("a", services.MaxNameLength+1),
				Sex:       models.SexMale,
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
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
			name:     "Error/InvalidFirstName",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: "$%&/()=?",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
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
			name:     "Error/FirstNameTooLong",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: strings.Repeat("a", services.MaxNameLength+1),
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
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
			name:     "Error/InvalidSex",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.Sex("invalid sex"),
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
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
			name:     "Error/NoFirstName",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				LastName: "last-name",
				Sex:      models.SexMale,
				Birthday: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
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
			name:     "Error/NoLastName",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: "name",
				Sex:       models.SexMale,
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
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
			name:     "Error/InvalidToken",
			tokenRaw: "string-token",
			now:      baseTime,
			form: models.UpdateIdentityForm{
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
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
			form: models.UpdateIdentityForm{
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			introspectTokenErr: fooErr,
			expectErr:          fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			identityDAO := daomocks.NewIdentityRepository(t)
			introspectTokenService := servicesmocks.NewIntrospectTokenService(t)

			introspectTokenService.
				On("IntrospectToken", context.Background(), d.tokenRaw, d.now, false).
				Return(d.introspectToken, d.introspectTokenErr)

			if d.shouldCallDAO {
				identityDAO.
					On("Update", context.Background(), &dao.IdentityModelCore{
						FirstName: d.form.FirstName,
						LastName:  d.form.LastName,
						Birthday:  d.form.Birthday,
						Sex:       d.form.Sex,
					}, d.introspectToken.Token.Payload.ID, d.now).
					Return(nil, d.daoErr)
			}

			service := services.NewUpdateIdentityService(identityDAO, introspectTokenService)
			err := service.UpdateIdentity(context.Background(), d.tokenRaw, d.now, d.form)

			require.ErrorIs(t, err, d.expectErr)

			identityDAO.AssertExpectations(t)
			introspectTokenService.AssertExpectations(t)
		})
	}
}
