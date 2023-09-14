package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	goframework "github.com/a-novel/go-framework"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	TokenKey0 = "eyJpYXQiOiIyMDIwLTA1LTA0VDA4OjAwOjAwWiIsImV4cCI6IjIwMjAtMDUtMDRUMDk6MDA6MDBaIiwiaWQiOiIxMDEwMTAxMC0xMDEwLTEwMTAtMTAxMC0xMDEwMTAxMDEwMTAifQ.eyJpZCI6IjAxMDEwMTAxLTAxMDEtMDEwMS0wMTAxLTAxMDEwMTAxMDEwMSJ9.9Zy2f7aqXxwI3F6SvYTu576NXLmsQZUIsp5V5F9QmvX280PV-ZzlHifnfNKXg7Gb2_xbJahHcvUnP-143Kn1BQ"
)

func TestGenerateToken(t *testing.T) {
	data := []struct {
		name string

		tokenTTL time.Duration

		data models.UserTokenPayload
		id   uuid.UUID
		now  time.Time

		list    []*dao.SecretKeyModel
		listErr error

		expect    *models.UserTokenStatus
		expectErr error
	}{
		{
			name:     "Success",
			tokenTTL: time.Hour,
			data:     models.UserTokenPayload{ID: goframework.NumberUUID(1)},
			id:       goframework.NumberUUID(10),
			now:      baseTime,
			list: []*dao.SecretKeyModel{
				{
					Name: "key-0",
					Key:  MockedSecretKeys[0],
				},
				{
					Name: "key-1",
					Key:  MockedSecretKeys[1],
				},
			},
			expect: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime,
						EXP: baseTime.Add(time.Hour),
						ID:  goframework.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
				TokenRaw: TokenKey0,
			},
		},
		{
			name:      "Error/DAOFailure",
			tokenTTL:  time.Hour,
			data:      models.UserTokenPayload{ID: goframework.NumberUUID(1)},
			id:        goframework.NumberUUID(10),
			now:       baseTime,
			listErr:   fooErr,
			expectErr: fooErr,
		},
		{
			name:      "Error/NoSignatureKeys",
			tokenTTL:  time.Hour,
			data:      models.UserTokenPayload{ID: goframework.NumberUUID(1)},
			id:        goframework.NumberUUID(10),
			now:       baseTime,
			expectErr: services.ErrMissingSignatureKeys,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			secretKeysDAO := daomocks.NewSecretKeysRepository(t)

			secretKeysDAO.On("List", context.Background()).Return(d.list, d.listErr)

			service := services.NewGenerateTokenService(secretKeysDAO, d.tokenTTL)
			token, err := service.GenerateToken(context.Background(), d.data, d.id, d.now)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, token)

			secretKeysDAO.AssertExpectations(t)
		})
	}
}

func TestTokenStatus(t *testing.T) {
	data := []struct {
		name string

		token string
		now   time.Time

		shouldCallList bool
		list           []*dao.SecretKeyModel
		listErr        error

		expect    *models.UserTokenStatus
		expectErr error
	}{
		{
			name:           "Success",
			token:          TokenKey0,
			now:            baseTime,
			shouldCallList: true,
			list: []*dao.SecretKeyModel{
				{
					Name: "key-2",
					Key:  MockedSecretKeys[2],
				},
				{
					Name: "key-0",
					Key:  MockedSecretKeys[0],
				},
				{
					Name: "key-1",
					Key:  MockedSecretKeys[1],
				},
			},
			expect: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime,
						EXP: baseTime.Add(time.Hour),
						ID:  goframework.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
				TokenRaw: TokenKey0,
			},
		},
		{
			name:   "Success/NoToken",
			now:    baseTime,
			expect: &models.UserTokenStatus{},
		},
		{
			name:           "Success/NotIssued",
			token:          TokenKey0,
			shouldCallList: true,
			now:            baseTime.Add(-time.Hour),
			list: []*dao.SecretKeyModel{
				{
					Name: "key-2",
					Key:  MockedSecretKeys[2],
				},
				{
					Name: "key-0",
					Key:  MockedSecretKeys[0],
				},
				{
					Name: "key-1",
					Key:  MockedSecretKeys[1],
				},
			},
			expect: &models.UserTokenStatus{
				NotIssued: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime,
						EXP: baseTime.Add(time.Hour),
						ID:  goframework.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
				TokenRaw: TokenKey0,
			},
		},
		{
			name:           "Success/Expired",
			token:          TokenKey0,
			shouldCallList: true,
			now:            baseTime.Add(2 * time.Hour),
			list: []*dao.SecretKeyModel{
				{
					Name: "key-2",
					Key:  MockedSecretKeys[2],
				},
				{
					Name: "key-0",
					Key:  MockedSecretKeys[0],
				},
				{
					Name: "key-1",
					Key:  MockedSecretKeys[1],
				},
			},
			expect: &models.UserTokenStatus{
				Expired: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime,
						EXP: baseTime.Add(time.Hour),
						ID:  goframework.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
				TokenRaw: TokenKey0,
			},
		},
		{
			name:           "Success/MissingSignatureKey",
			token:          TokenKey0,
			shouldCallList: true,
			now:            baseTime,
			list: []*dao.SecretKeyModel{
				{
					Name: "key-2",
					Key:  MockedSecretKeys[2],
				},
				{
					Name: "key-1",
					Key:  MockedSecretKeys[1],
				},
			},
			expect: &models.UserTokenStatus{
				Expired: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime,
						EXP: baseTime.Add(time.Hour),
						ID:  goframework.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
				TokenRaw: TokenKey0,
			},
		},
		{
			name:           "Error/DAOFailure",
			token:          TokenKey0,
			shouldCallList: true,
			now:            baseTime,
			listErr:        fooErr,
			expectErr:      fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			secretKeysDAO := daomocks.NewSecretKeysRepository(t)

			if d.shouldCallList {
				secretKeysDAO.On("List", context.Background()).Return(d.list, d.listErr)
			}

			service := services.NewGetTokenStatusService(secretKeysDAO)
			token, err := service.GetTokenStatus(context.Background(), d.token, d.now)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, token)

			secretKeysDAO.AssertExpectations(t)
		})
	}
}
