package dao_test

import (
	"context"
	"github.com/a-novel/auth-service/migrations"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/bunovel"
	goframework "github.com/a-novel/go-framework"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"io/fs"
	"testing"
	"time"
)

func TestIdentityRepository_GetIdentity(t *testing.T) {
	db, sqlDB := bunovel.GetTestPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.IdentityModel{
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime, &baseTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "name-1",
				LastName:  "last-name-1",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
	}

	data := []struct {
		name string

		id uuid.UUID

		expect    *dao.IdentityModel
		expectErr error
	}{
		{
			name:   "Success",
			id:     goframework.NumberUUID(1000),
			expect: fixtures[0],
		},
		{
			name:      "Error/NotFound",
			id:        goframework.NumberUUID(1),
			expectErr: bunovel.ErrNotFound,
		},
	}

	err := bunovel.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := dao.NewIdentityRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.GetIdentity(ctx, d.id)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestIdentityRepository_Update(t *testing.T) {
	db, sqlDB := bunovel.GetTestPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.IdentityModel{
		{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime, &baseTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "name-1",
				LastName:  "last-name-1",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
	}

	data := []struct {
		name string

		core *dao.IdentityModelCore
		id   uuid.UUID
		now  time.Time

		expect    *dao.IdentityModel
		expectErr error
	}{
		{
			name: "Success",
			core: &dao.IdentityModelCore{
				FirstName: "name-2",
				LastName:  "last-name-2",
				Birthday:  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			id:  goframework.NumberUUID(1000),
			now: updateTime,
			expect: &dao.IdentityModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime, &updateTime),
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "name-2",
					LastName:  "last-name-2",
					Birthday:  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexFemale,
				},
			},
		},
		{
			name: "Success/NonRomanizedName",
			core: &dao.IdentityModelCore{
				FirstName: "ルイズ フランソワーズ ル ブラン",
				LastName:  "ド ラ ヴァリエール",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
			id:  goframework.NumberUUID(1000),
			now: updateTime,
			expect: &dao.IdentityModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime, &updateTime),
				IdentityModelCore: dao.IdentityModelCore{
					FirstName: "ルイズ フランソワーズ ル ブラン",
					LastName:  "ド ラ ヴァリエール",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
			},
		},
		{
			name: "Error/NotFound",
			core: &dao.IdentityModelCore{
				FirstName: "name-2",
				LastName:  "last-name-2",
				Birthday:  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
			id:        goframework.NumberUUID(1),
			now:       updateTime,
			expectErr: bunovel.ErrNotFound,
		},
	}

	err := bunovel.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := dao.NewIdentityRepository(tx)

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
