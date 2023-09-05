package services_test

import (
	"context"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSlugExists(t *testing.T) {
	data := []struct {
		name string

		slug string

		slugDAO    bool
		slugDAOErr error

		expect    bool
		expectErr error
	}{
		{
			name:    "Success",
			slug:    "slug",
			slugDAO: true,
			expect:  true,
		},
		{
			name:    "Success/NotFound",
			slug:    "slug",
			slugDAO: false,
			expect:  false,
		},
		{
			name:       "Error/DAOFailure",
			slug:       "slug",
			slugDAOErr: fooErr,
			expectErr:  fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			profileDAO := daomocks.NewProfileRepository(t)

			profileDAO.
				On("SlugExists", context.Background(), d.slug).
				Return(d.slugDAO, d.slugDAOErr)

			service := services.NewSlugExistsService(profileDAO)
			exists, err := service.SlugExists(context.Background(), d.slug)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, exists)

			profileDAO.AssertExpectations(t)
		})
	}
}
