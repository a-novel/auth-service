package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/a-novel/go-framework/postgresql"
	"github.com/a-novel/go-framework/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPreview(t *testing.T) {
	data := []struct {
		name string

		slug string

		shouldCallProfileDAO bool
		profileDAO           *dao.ProfileModel
		profileDAOErr        error

		shouldCallIdentityDAO bool
		identityDAO           *dao.IdentityModel
		identityDAOErr        error

		expect    *models.UserPreview
		expectErr error
	}{
		{
			name:                 "Success",
			slug:                 "slug",
			shouldCallProfileDAO: true,
			profileDAO: &dao.ProfileModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, &baseTime),
				ProfileModelCore: dao.ProfileModelCore{
					Slug: "slug",
				},
			},
			shouldCallIdentityDAO: true,
			identityDAO: &dao.IdentityModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, &baseTime),
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name",
					LastName:  "last-name",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
			},
			expect: &models.UserPreview{
				FirstName: "name",
				LastName:  "last-name",
				Slug:      "slug",
				CreatedAt: baseTime,
			},
		},
		{
			name:                 "Success/WithUsername",
			slug:                 "slug",
			shouldCallProfileDAO: true,
			profileDAO: &dao.ProfileModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, &baseTime),
				ProfileModelCore: dao.ProfileModelCore{
					Username: "username",
					Slug:     "slug",
				},
			},
			shouldCallIdentityDAO: true,
			identityDAO: &dao.IdentityModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, &baseTime),
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name",
					LastName:  "last-name",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
			},
			expect: &models.UserPreview{
				Username:  "username",
				Slug:      "slug",
				CreatedAt: baseTime,
			},
		},
		{
			name:                 "Error/IdentityDAOFailure",
			slug:                 "slug",
			shouldCallProfileDAO: true,
			profileDAO: &dao.ProfileModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, &baseTime),
				ProfileModelCore: dao.ProfileModelCore{
					Slug: "slug",
				},
			},
			shouldCallIdentityDAO: true,
			identityDAOErr:        fooErr,
			expectErr:             fooErr,
		},
		{
			name:                 "Error/ProfileDAOFailure",
			slug:                 "slug",
			shouldCallProfileDAO: true,
			profileDAOErr:        fooErr,
			expectErr:            fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			profileDAO := daomocks.NewProfileRepository(t)
			identityDAO := daomocks.NewIdentityRepository(t)

			if d.shouldCallProfileDAO {
				profileDAO.
					On("GetProfileBySlug", context.Background(), d.slug).
					Return(d.profileDAO, d.profileDAOErr)
			}

			if d.shouldCallIdentityDAO {
				identityDAO.
					On("GetIdentity", context.Background(), d.profileDAO.ID).
					Return(d.identityDAO, d.identityDAOErr)
			}

			service := services.NewPreviewService(profileDAO, identityDAO)
			res, err := service.Preview(context.Background(), d.slug)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, res)

			profileDAO.AssertExpectations(t)
			identityDAO.AssertExpectations(t)
		})
	}
}
