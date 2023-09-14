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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	data := []struct {
		name string

		email    string
		password string
		now      time.Time

		shouldCallDAO bool
		daoResponse   *dao.CredentialsModel
		daoErr        error

		shouldCallGenerateToken bool
		generateTokenStatus     *models.UserTokenStatus
		generateTokenErr        error

		expect    *models.UserTokenStatus
		expectErr error
	}{
		{
			name:          "Success",
			email:         "user@domain.com",
			password:      password,
			now:           baseTime,
			shouldCallDAO: true,
			daoResponse: &dao.CredentialsModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    dao.Email{User: "user", Domain: "domain.com"},
					Password: dao.Password{Hashed: passwordEncrypted},
				},
			},
			shouldCallGenerateToken: true,
			generateTokenStatus: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
			expect: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
			},
		},
		{
			name:          "Error/GenerateTokenFailure",
			email:         "user@domain.com",
			password:      password,
			now:           baseTime,
			shouldCallDAO: true,
			daoResponse: &dao.CredentialsModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    dao.Email{User: "user", Domain: "domain.com"},
					Password: dao.Password{Hashed: passwordEncrypted},
				},
			},
			shouldCallGenerateToken: true,
			generateTokenErr:        fooErr,
			expectErr:               fooErr,
		},
		{
			name:          "Error/WrongPassword",
			email:         "user@domain.com",
			password:      "fake-password",
			now:           baseTime,
			shouldCallDAO: true,
			daoResponse: &dao.CredentialsModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &baseTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    dao.Email{User: "user", Domain: "domain.com"},
					Password: dao.Password{Hashed: passwordEncrypted},
				},
			},
			expectErr: goframework.ErrInvalidCredentials,
		},
		{
			name:          "Error/CredentialsDAOFailure",
			email:         "user@domain.com",
			password:      password,
			now:           baseTime,
			shouldCallDAO: true,
			daoErr:        fooErr,
			expectErr:     fooErr,
		},
		{
			name:      "Error/InvalidEmail",
			email:     "userdomain.com",
			password:  password,
			now:       baseTime,
			expectErr: goframework.ErrInvalidEntity,
		},
		{
			name:      "Error/EmailTooLong",
			email:     strings.Repeat("a", services.MaxEmailLength) + "@domain.com",
			password:  password,
			now:       baseTime,
			expectErr: goframework.ErrInvalidEntity,
		},
		{
			name:      "Error/PasswordTooLong",
			email:     strings.Repeat("a", services.MaxPasswordLength) + "x",
			password:  password,
			now:       baseTime,
			expectErr: goframework.ErrInvalidEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsDAO := daomocks.NewCredentialsRepository(t)
			generateTokenService := servicesmocks.NewGenerateTokenService(t)

			if d.shouldCallDAO {
				credentialsDAO.
					On("GetCredentialsByEmail", context.Background(), mock.Anything).
					Return(d.daoResponse, d.daoErr)
			}

			if d.shouldCallGenerateToken {
				generateTokenService.
					On("GenerateToken", context.Background(), models.UserTokenPayload{ID: d.daoResponse.ID}, mock.Anything, d.now).
					Return(d.generateTokenStatus, d.generateTokenErr)
			}

			service := services.NewLoginService(credentialsDAO, generateTokenService)
			res, err := service.Login(context.Background(), d.email, d.password, d.now)

			require.Equal(t, d.expect, res)
			require.ErrorIs(t, err, d.expectErr)

			credentialsDAO.AssertExpectations(t)
			generateTokenService.AssertExpectations(t)
		})
	}
}
