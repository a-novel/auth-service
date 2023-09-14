package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/services"
	goframework "github.com/a-novel/go-framework"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmailExists(t *testing.T) {
	data := []struct {
		name string

		email    string
		daoEmail dao.Email

		shouldCallCredentialsDAO bool

		emailExists    bool
		emailExistsErr error

		expect    bool
		expectErr error
	}{
		{
			name:  "Success",
			email: "user@domain.com",
			daoEmail: dao.Email{
				User:   "user",
				Domain: "domain.com",
			},
			shouldCallCredentialsDAO: true,
			emailExists:              true,
			expect:                   true,
		},
		{
			name:  "Success/NotFound",
			email: "user@domain.com",
			daoEmail: dao.Email{
				User:   "user",
				Domain: "domain.com",
			},
			shouldCallCredentialsDAO: true,
			emailExists:              false,
			expect:                   false,
		},
		{
			name:  "Error/DAOFailure",
			email: "user@domain.com",
			daoEmail: dao.Email{
				User:   "user",
				Domain: "domain.com",
			},
			shouldCallCredentialsDAO: true,
			emailExistsErr:           fooErr,
			expectErr:                fooErr,
		},
		{
			name:      "Error/BadEmailFormat",
			email:     "userdomain.com",
			expectErr: goframework.ErrInvalidEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsDAO := daomocks.NewCredentialsRepository(t)

			if d.shouldCallCredentialsDAO {
				credentialsDAO.
					On("EmailExists", context.Background(), d.daoEmail).
					Return(d.emailExists, d.emailExistsErr)
			}

			service := services.NewEmailExistsService(credentialsDAO)
			ok, err := service.EmailExists(context.Background(), d.email)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, ok)

			credentialsDAO.AssertExpectations(t)
		})
	}
}
