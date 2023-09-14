package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/bunovel"
	goframework "github.com/a-novel/go-framework"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSearch(t *testing.T) {
	data := []struct {
		name string

		query  string
		limit  int
		offset int

		shouldCallUserDAO bool
		userDAO           []*dao.UserModel
		userDAOCount      int
		userDAOErr        error

		expect      []*models.UserPreview
		expectCount int
		expectErr   error
	}{
		{
			name:              "Success",
			query:             "query",
			limit:             10,
			offset:            0,
			shouldCallUserDAO: true,
			userDAO: []*dao.UserModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, &updateTime),
					UserModelCore: dao.UserModelCore{
						Identity: dao.IdentityModelCore{
							FirstName: "name-1",
							LastName:  "surname-1",
						},
						Profile: dao.ProfileModelCore{
							Slug: "slug-1",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(2), baseTime, &updateTime),
					UserModelCore: dao.UserModelCore{
						Identity: dao.IdentityModelCore{
							FirstName: "name-2",
							LastName:  "surname-2",
						},
						Profile: dao.ProfileModelCore{
							Username: "username-2",
							Slug:     "slug-2",
						},
					},
				},
			},
			userDAOCount: 20,
			expect: []*models.UserPreview{
				{
					FirstName: "name-1",
					LastName:  "surname-1",
					Slug:      "slug-1",
					CreatedAt: baseTime,
				},
				{
					Username:  "username-2",
					Slug:      "slug-2",
					CreatedAt: baseTime,
				},
			},
			expectCount: 20,
		},
		{
			name:              "Success/NoResult",
			query:             "query",
			limit:             10,
			offset:            30,
			shouldCallUserDAO: true,
			userDAOCount:      20,
			expect:            []*models.UserPreview{},
			expectCount:       20,
		},
		{
			name:              "Error/DAOFailure",
			query:             "query",
			limit:             10,
			offset:            30,
			shouldCallUserDAO: true,
			userDAOErr:        fooErr,
			expectErr:         fooErr,
		},
		{
			name:      "Error/LimitTooHigh",
			query:     "query",
			limit:     services.MaxUserSearchLimit + 1,
			offset:    30,
			expectErr: goframework.ErrInvalidEntity,
		},
		{
			name:      "Error/NoLimit",
			query:     "query",
			offset:    30,
			expectErr: goframework.ErrInvalidEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			userDAO := daomocks.NewUserRepository(t)

			if d.shouldCallUserDAO {
				userDAO.
					On("Search", context.Background(), d.query, d.limit, d.offset).
					Return(d.userDAO, d.userDAOCount, d.userDAOErr)
			}

			service := services.NewSearchService(userDAO)
			users, total, err := service.Search(context.Background(), d.query, d.limit, d.offset)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expectCount, total)
			require.Equal(t, d.expect, users)

			userDAO.AssertExpectations(t)
		})
	}
}
