package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/go-framework/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestIntrospectToken(t *testing.T) {
	data := []struct {
		name string

		tokenRefreshThreshold time.Duration

		token       string
		now         time.Time
		autoRefresh bool

		tokenStatus    *models.UserTokenStatus
		tokenStatusErr error

		shouldCallGenerateToken bool
		generateTokenStatus     *models.UserTokenStatus
		generateTokenErr        error

		expect    *models.UserTokenStatus
		expectErr error
	}{
		{
			name:                  "Success/NoRefresh",
			tokenRefreshThreshold: 15 * time.Minute,
			token:                 "string-token",
			now:                   baseTime,
			autoRefresh:           false,
			tokenStatus: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime.Add(-time.Hour),
						EXP: baseTime.Add(30 * time.Minute),
						ID:  test.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			expect: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime.Add(-time.Hour),
						EXP: baseTime.Add(30 * time.Minute),
						ID:  test.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
		},
		{
			name:                  "Success/Refresh/NotYet",
			tokenRefreshThreshold: 15 * time.Minute,
			token:                 "string-token",
			now:                   baseTime,
			autoRefresh:           true,
			tokenStatus: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime.Add(-time.Hour),
						EXP: baseTime.Add(30 * time.Minute),
						ID:  test.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			expect: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime.Add(-time.Hour),
						EXP: baseTime.Add(30 * time.Minute),
						ID:  test.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
		},
		{
			name:                  "Success/Refresh",
			tokenRefreshThreshold: 15 * time.Minute,
			token:                 "string-token",
			now:                   baseTime.Add(15 * time.Minute),
			autoRefresh:           true,
			tokenStatus: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime.Add(-time.Hour),
						EXP: baseTime.Add(30 * time.Minute),
						ID:  test.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallGenerateToken: true,
			generateTokenStatus: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime.Add(15 * time.Hour),
						EXP: baseTime.Add(75 * time.Minute),
						ID:  test.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			expect: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime.Add(15 * time.Hour),
						EXP: baseTime.Add(75 * time.Minute),
						ID:  test.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
		},
		{
			name:                  "Error/RefreshFailure",
			tokenRefreshThreshold: 15 * time.Minute,
			token:                 "string-token",
			now:                   baseTime.Add(15 * time.Minute),
			autoRefresh:           true,
			tokenStatus: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime.Add(-time.Hour),
						EXP: baseTime.Add(30 * time.Minute),
						ID:  test.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallGenerateToken: true,
			generateTokenErr:        fooErr,
			expectErr:               fooErr,
		},
		{
			name:                  "Error/GetTokenStatusFailure",
			tokenRefreshThreshold: 15 * time.Minute,
			token:                 "string-token",
			now:                   baseTime,
			autoRefresh:           false,
			tokenStatusErr:        fooErr,
			expectErr:             fooErr,
		},
		{
			name:                  "Success/InvalidToken",
			tokenRefreshThreshold: 15 * time.Minute,
			token:                 "string-token",
			now:                   baseTime,
			autoRefresh:           false,
			tokenStatus:           &models.UserTokenStatus{OK: false},
			expect:                &models.UserTokenStatus{OK: false},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			getTokenStatusService := servicesmocks.NewGetTokenStatusService(t)
			generateTokenService := servicesmocks.NewGenerateTokenService(t)

			getTokenStatusService.
				On("GetTokenStatus", context.Background(), d.token, d.now).
				Return(d.tokenStatus, d.tokenStatusErr)

			if d.shouldCallGenerateToken {
				generateTokenService.
					On("GenerateToken", context.Background(), d.tokenStatus.Token.Payload, d.tokenStatus.Token.Header.ID, d.now).
					Return(d.generateTokenStatus, d.generateTokenErr)
			}

			service := services.NewIntrospectTokenService(generateTokenService, getTokenStatusService, d.tokenRefreshThreshold)
			status, err := service.IntrospectToken(context.Background(), d.token, d.now, d.autoRefresh)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, status)

			getTokenStatusService.AssertExpectations(t)
			generateTokenService.AssertExpectations(t)
		})
	}
}
