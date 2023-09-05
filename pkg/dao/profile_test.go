package dao_test

import (
	"context"
	"github.com/a-novel/auth-service/migrations"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/postgresql"
	"github.com/a-novel/go-framework/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"io/fs"
	"testing"
	"time"
)

func TestProfileRepository_GetProfile(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.ProfileModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			ProfileModelCore: dao.ProfileModelCore{
				Username: "username-1",
				Slug:     "slug-1",
			},
		},
	}

	data := []struct {
		name string

		id uuid.UUID

		expect    *dao.ProfileModel
		expectErr error
	}{
		{
			name:   "Success",
			id:     test.NumberUUID(1000),
			expect: fixtures[0],
		},
		{
			name:      "Error/NotFound",
			id:        test.NumberUUID(1),
			expectErr: errors.ErrNotFound,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := dao.NewProfileRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.GetProfile(ctx, d.id)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestProfileRepository_GetProfileBySlug(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.ProfileModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			ProfileModelCore: dao.ProfileModelCore{
				Username: "username-1",
				Slug:     "slug-1",
			},
		},
	}

	data := []struct {
		name string

		slug string

		expect    *dao.ProfileModel
		expectErr error
	}{
		{
			name:   "Success",
			slug:   "slug-1",
			expect: fixtures[0],
		},
		{
			name:      "Error/NotFound",
			slug:      "fake-slug",
			expectErr: errors.ErrNotFound,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := dao.NewProfileRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.GetProfileBySlug(ctx, d.slug)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestProfileRepository_SlugExists(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.ProfileModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			ProfileModelCore: dao.ProfileModelCore{
				Username: "username-1",
				Slug:     "slug-1",
			},
		},
	}

	data := []struct {
		name string

		slug string

		expect    bool
		expectErr error
	}{
		{
			name:   "Success",
			slug:   "slug-1",
			expect: true,
		},
		{
			name:   "Success/NotExist",
			slug:   "fake-slug",
			expect: false,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := dao.NewProfileRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.SlugExists(ctx, d.slug)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestProfileRepository_Update(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.ProfileModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			ProfileModelCore: dao.ProfileModelCore{
				Username: "username-1",
				Slug:     "slug-1",
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &baseTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug: "slug-2",
			},
		},
	}

	data := []struct {
		name string

		core *dao.ProfileModelCore
		id   uuid.UUID
		now  time.Time

		expect    *dao.ProfileModel
		expectErr error
	}{
		{
			name: "Success",
			core: &dao.ProfileModelCore{
				Username: "new-username-1",
				Slug:     "new-slug-1",
			},
			id:  test.NumberUUID(1000),
			now: updateTime,
			expect: &dao.ProfileModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				ProfileModelCore: dao.ProfileModelCore{
					Username: "new-username-1",
					Slug:     "new-slug-1",
				},
			},
		},
		{
			name: "Success/RemoveUsername",
			core: &dao.ProfileModelCore{
				Slug: "new-slug-1",
			},
			id:  test.NumberUUID(1000),
			now: updateTime,
			expect: &dao.ProfileModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				ProfileModelCore: dao.ProfileModelCore{
					Slug: "new-slug-1",
				},
			},
		},
		{
			name: "Success/AddUsername",
			core: &dao.ProfileModelCore{
				Username: "new-username-2",
				Slug:     "new-slug-2",
			},
			id:  test.NumberUUID(1001),
			now: updateTime,
			expect: &dao.ProfileModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &updateTime),
				ProfileModelCore: dao.ProfileModelCore{
					Username: "new-username-2",
					Slug:     "new-slug-2",
				},
			},
		},
		{
			name: "Error/NotFound",
			core: &dao.ProfileModelCore{
				Username: "new-username-1",
				Slug:     "new-slug-1",
			},
			id:        test.NumberUUID(1),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
		{
			name: "Error/RemoveSlug",
			core: &dao.ProfileModelCore{
				Username: "new-username-1",
			},
			id:        test.NumberUUID(1000),
			now:       updateTime,
			expectErr: errors.ErrConstraintViolation,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := dao.NewProfileRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.Update(ctx, d.core, d.id, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}
