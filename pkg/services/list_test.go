package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/postgresql"
	"github.com/a-novel/go-framework/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestList(t *testing.T) {
	data := []struct {
		name string

		ids []uuid.UUID

		daoResponse []*dao.UserModel
		daoErr      error

		expect    []*models.UserPreview
		expectErr error
	}{
		{
			name: "Success",
			ids:  []uuid.UUID{test.NumberUUID(1), test.NumberUUID(2)},
			daoResponse: []*dao.UserModel{
				{
					Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: dao.Email{User: "user1", Domain: "domain.com"},
						},
						Identity: dao.IdentityModelCore{
							FirstName: "name-1",
							LastName:  "last-name-1",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexFemale,
						},
						Profile: dao.ProfileModelCore{
							Username: "username-1",
							Slug:     "slug-1",
						},
					},
				},
				{
					Metadata: postgresql.NewMetadata(test.NumberUUID(2), baseTime, &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: dao.Email{User: "user2", Domain: "domain.com"},
						},
						Identity: dao.IdentityModelCore{
							FirstName: "name-2",
							LastName:  "last-name-2",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "slug-2",
						},
					},
				},
			},
			expect: []*models.UserPreview{
				{
					Username:  "username-1",
					Slug:      "slug-1",
					CreatedAt: baseTime,
				},
				{
					FirstName: "name-2",
					LastName:  "last-name-2",
					Slug:      "slug-2",
					CreatedAt: baseTime,
				},
			},
		},
		{
			name:        "Success/NoResults",
			ids:         []uuid.UUID{test.NumberUUID(1), test.NumberUUID(2)},
			daoResponse: nil,
			expect:      []*models.UserPreview{},
		},
		{
			name:      "Error/DAOFailure",
			ids:       []uuid.UUID{test.NumberUUID(1), test.NumberUUID(2)},
			daoErr:    fooErr,
			expectErr: fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			userDAO := daomocks.NewUserRepository(t)

			userDAO.On("List", context.Background(), d.ids).Return(d.daoResponse, d.daoErr)

			service := services.NewListService(userDAO)
			users, err := service.List(context.Background(), d.ids)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, users)

			userDAO.AssertExpectations(t)
		})
	}
}
