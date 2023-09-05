package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/test"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestUpdatePassword(t *testing.T) {
	data := []struct {
		name string

		form models.UpdatePasswordForm
		now  time.Time

		shouldCallGetCredentials bool
		getCredentials           *dao.CredentialsModel
		getCredentialsErr        error

		shouldCallUpdateCredentials bool
		updateCredentialsErr        error

		expectErr error
	}{
		{
			name: "Success",
			form: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				NewPassword: "new-secure-password",
				OldPassword: password,
			},
			now:                      baseTime,
			shouldCallGetCredentials: true,
			getCredentials: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Password: dao.Password{Hashed: passwordEncrypted},
				},
			},
			shouldCallUpdateCredentials: true,
		},
		{
			name: "Success/ValidationCode",
			form: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				NewPassword: "new-secure-password",
				Code:        publicValidationCode,
			},
			now:                      baseTime,
			shouldCallGetCredentials: true,
			getCredentials: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Password: dao.Password{Hashed: passwordEncrypted, Validation: privateValidationCode},
				},
			},
			shouldCallUpdateCredentials: true,
		},
		{
			name: "Error/UpdatePasswordFailure",
			form: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				NewPassword: "new-secure-password",
				OldPassword: password,
			},
			now:                      baseTime,
			shouldCallGetCredentials: true,
			getCredentials: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Password: dao.Password{Hashed: passwordEncrypted},
				},
			},
			shouldCallUpdateCredentials: true,
			updateCredentialsErr:        fooErr,
			expectErr:                   fooErr,
		},
		{
			name: "Error/WrongPassword",
			form: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				NewPassword: "new-secure-password",
				OldPassword: "fake-password",
			},
			now:                      baseTime,
			shouldCallGetCredentials: true,
			getCredentials: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Password: dao.Password{Hashed: passwordEncrypted},
				},
			},
			expectErr: errors.ErrInvalidCredentials,
		},
		{
			name: "Error/WrongValidationCode",
			form: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				NewPassword: "new-secure-password",
				Code:        "fakecode",
			},
			now:                      baseTime,
			shouldCallGetCredentials: true,
			getCredentials: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Password: dao.Password{Hashed: passwordEncrypted, Validation: privateValidationCode},
				},
			},
			expectErr: errors.ErrInvalidCredentials,
		},
		{
			name: "Error/ValidationCodeExpired",
			form: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				NewPassword: "new-secure-password",
				Code:        publicValidationCode,
			},
			now:                      baseTime,
			shouldCallGetCredentials: true,
			getCredentials: &dao.CredentialsModel{
				CredentialsModelCore: dao.CredentialsModelCore{
					Password: dao.Password{Hashed: passwordEncrypted},
				},
			},
			expectErr: errors.ErrInvalidCredentials,
		},
		{
			name: "Error/NoNewPassword",
			form: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				OldPassword: password,
			},
			now:       baseTime,
			expectErr: errors.ErrInvalidEntity,
		},
		{
			name: "Error/NoIdentityConfirmation",
			form: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				NewPassword: "new-secure-password",
			},
			now:       baseTime,
			expectErr: errors.ErrInvalidEntity,
		},
		{
			name: "Error/NewPasswordTooShort",
			form: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				NewPassword: "p",
				OldPassword: password,
			},
			now:       baseTime,
			expectErr: errors.ErrInvalidEntity,
		},
		{
			name: "Error/NewPasswordTooLong",
			form: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				NewPassword: strings.Repeat("a", services.MaxPasswordLength+1),
				OldPassword: password,
			},
			now:       baseTime,
			expectErr: errors.ErrInvalidEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsDAO := daomocks.NewCredentialsRepository(t)

			if d.shouldCallGetCredentials {
				credentialsDAO.
					On("GetCredentials", context.Background(), d.form.ID).
					Return(d.getCredentials, d.getCredentialsErr)
			}

			if d.shouldCallUpdateCredentials {
				credentialsDAO.
					On("UpdatePassword", context.Background(), d.form.NewPassword, d.form.ID, d.now).
					Return(nil, d.updateCredentialsErr)
			}

			service := services.NewUpdatePasswordService(credentialsDAO)
			err := service.UpdatePassword(context.Background(), d.form, d.now)

			require.ErrorIs(t, err, d.expectErr)

			credentialsDAO.AssertExpectations(t)
		})
	}
}
