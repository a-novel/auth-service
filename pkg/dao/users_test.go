package dao_test

import (
	"context"
	"github.com/a-novel/auth-service/migrations"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/bunovel"
	goframework "github.com/a-novel/go-framework"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"io/fs"
	"sort"
	"testing"
	"time"
)

func TestUserRepository_Create(t *testing.T) {
	db, sqlDB := bunovel.GetTestPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []interface{}{
		&dao.CredentialsModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user2@domain.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		&dao.ProfileModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime, &baseTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug: "slug-2",
			},
		},
	}

	data := []struct {
		name string

		data *dao.UserModelCore
		id   uuid.UUID
		now  time.Time

		expect    *dao.UserModel
		expectErr error
	}{
		{
			name: "Success",
			data: &dao.UserModelCore{
				Credentials: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
				Identity: dao.IdentityModelCore{
					FirstName: "name-1",
					LastName:  "last-name-1",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: dao.ProfileModelCore{
					Slug: "slug-1",
				},
			},
			id:  goframework.NumberUUID(1),
			now: baseTime,
			expect: &dao.UserModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, nil),
				UserModelCore: dao.UserModelCore{
					Credentials: dao.CredentialsModelCore{
						Email:    MustParseEmail("user1@domain.com"),
						Password: dao.Password{Hashed: "password-hashed"},
					},
					Identity: dao.IdentityModelCore{
						FirstName: "name-1",
						LastName:  "last-name-1",
						Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						Sex:       models.SexMale,
					},
					Profile: dao.ProfileModelCore{
						Slug: "slug-1",
					},
				},
			},
		},
		{
			name: "Success/WithEmailValidation",
			data: &dao.UserModelCore{
				Credentials: dao.CredentialsModelCore{
					Email:    MustParseEmailWithValidation("user1@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
				Identity: dao.IdentityModelCore{
					FirstName: "name-1",
					LastName:  "last-name-1",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: dao.ProfileModelCore{
					Slug: "slug-1",
				},
			},
			id:  goframework.NumberUUID(1),
			now: baseTime,
			expect: &dao.UserModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, nil),
				UserModelCore: dao.UserModelCore{
					Credentials: dao.CredentialsModelCore{
						Email:    MustParseEmailWithValidation("user1@domain.com", "validation-code"),
						Password: dao.Password{Hashed: "password-hashed"},
					},
					Identity: dao.IdentityModelCore{
						FirstName: "name-1",
						LastName:  "last-name-1",
						Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						Sex:       models.SexMale,
					},
					Profile: dao.ProfileModelCore{
						Slug: "slug-1",
					},
				},
			},
		},
		{
			name: "Success/WithPasswordReset",
			data: &dao.UserModelCore{
				Credentials: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					Password: dao.Password{Validation: "validation-code"},
				},
				Identity: dao.IdentityModelCore{
					FirstName: "name-1",
					LastName:  "last-name-1",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: dao.ProfileModelCore{
					Slug: "slug-1",
				},
			},
			id:  goframework.NumberUUID(1),
			now: baseTime,
			expect: &dao.UserModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, nil),
				UserModelCore: dao.UserModelCore{
					Credentials: dao.CredentialsModelCore{
						Email:    MustParseEmail("user1@domain.com"),
						Password: dao.Password{Validation: "validation-code"},
					},
					Identity: dao.IdentityModelCore{
						FirstName: "name-1",
						LastName:  "last-name-1",
						Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						Sex:       models.SexMale,
					},
					Profile: dao.ProfileModelCore{
						Slug: "slug-1",
					},
				},
			},
		},
		{
			name: "Success/WithUsername",
			data: &dao.UserModelCore{
				Credentials: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
				Identity: dao.IdentityModelCore{
					FirstName: "name-1",
					LastName:  "last-name-1",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: dao.ProfileModelCore{
					Slug:     "slug-1",
					Username: "username-1",
				},
			},
			id:  goframework.NumberUUID(1),
			now: baseTime,
			expect: &dao.UserModel{
				Metadata: bunovel.NewMetadata(goframework.NumberUUID(1), baseTime, nil),
				UserModelCore: dao.UserModelCore{
					Credentials: dao.CredentialsModelCore{
						Email:    MustParseEmail("user1@domain.com"),
						Password: dao.Password{Hashed: "password-hashed"},
					},
					Identity: dao.IdentityModelCore{
						FirstName: "name-1",
						LastName:  "last-name-1",
						Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						Sex:       models.SexMale,
					},
					Profile: dao.ProfileModelCore{
						Slug:     "slug-1",
						Username: "username-1",
					},
				},
			},
		},
		{
			name: "Error/EmailTaken",
			data: &dao.UserModelCore{
				Credentials: dao.CredentialsModelCore{
					Email:    MustParseEmail("user2@domain.com"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
				Identity: dao.IdentityModelCore{
					FirstName: "name-1",
					LastName:  "last-name-1",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: dao.ProfileModelCore{
					Slug: "slug-1",
				},
			},
			id:        goframework.NumberUUID(1),
			now:       baseTime,
			expectErr: bunovel.ErrUniqConstraintViolation,
		},
		{
			name: "Error/EmailUserMissing",
			data: &dao.UserModelCore{
				Credentials: dao.CredentialsModelCore{
					Email:    dao.Email{Domain: "domain.com"},
					Password: dao.Password{Hashed: "password-hashed"},
				},
				Identity: dao.IdentityModelCore{
					FirstName: "name-1",
					LastName:  "last-name-1",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: dao.ProfileModelCore{
					Slug: "slug-1",
				},
			},
			id:        goframework.NumberUUID(1),
			now:       baseTime,
			expectErr: bunovel.ErrConstraintViolation,
		},
		{
			name: "Error/EmailDomainMissing",
			data: &dao.UserModelCore{
				Credentials: dao.CredentialsModelCore{
					Email:    dao.Email{User: "user1"},
					Password: dao.Password{Hashed: "password-hashed"},
				},
				Identity: dao.IdentityModelCore{
					FirstName: "name-1",
					LastName:  "last-name-1",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: dao.ProfileModelCore{
					Slug: "slug-1",
				},
			},
			id:        goframework.NumberUUID(1),
			now:       baseTime,
			expectErr: bunovel.ErrConstraintViolation,
		},
		{
			name: "Error/SlugTaken",
			data: &dao.UserModelCore{
				Credentials: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
				Identity: dao.IdentityModelCore{
					FirstName: "name-1",
					LastName:  "last-name-1",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: dao.ProfileModelCore{
					Slug: "slug-2",
				},
			},
			id:        goframework.NumberUUID(1),
			now:       baseTime,
			expectErr: bunovel.ErrUniqConstraintViolation,
		},
		{
			name: "Error/SlugMissing",
			data: &dao.UserModelCore{
				Credentials: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
				Identity: dao.IdentityModelCore{
					FirstName: "name-1",
					LastName:  "last-name-1",
					Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					Sex:       models.SexMale,
				},
				Profile: dao.ProfileModelCore{},
			},
			id:        goframework.NumberUUID(1),
			now:       baseTime,
			expectErr: bunovel.ErrConstraintViolation,
		},
	}

	err := bunovel.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := dao.NewUserRepository(stx).Create(ctx, d.data, d.id, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestUserRepository_Search(t *testing.T) {
	db, sqlDB := bunovel.GetTestPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []interface{}{
		// User 1
		&dao.CredentialsModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("elon.bezos@space-origin.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		&dao.IdentityModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "Elon",
				LastName:  "Bezos",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
		&dao.ProfileModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug: "space-origin",
			},
		},

		// User 2
		&dao.CredentialsModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("jeff.musk@amazon.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		&dao.IdentityModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "Jeff",
				LastName:  "Bezos",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
		&dao.ProfileModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug: "big-brother",
			},
		},

		// User 3
		&dao.CredentialsModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1002), baseTime.Add(2*time.Hour), &updateTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("sarasaka@corpo.plaza"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		&dao.IdentityModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1002), baseTime.Add(2*time.Hour), &updateTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "Saburo",
				LastName:  "Arasaka",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
		},
		&dao.ProfileModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1002), baseTime.Add(2*time.Hour), &updateTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug:     "i-dont-have-any-ideas-anymore",
				Username: "Mikoshi",
			},
		},

		// User 4
		&dao.CredentialsModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("real-elon-musk@tesla.tx"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		&dao.IdentityModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "Elon",
				LastName:  "Musk",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
		&dao.ProfileModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug: "real-elon",
			},
		},

		// User 5
		&dao.CredentialsModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1004), baseTime.Add(6*time.Hour), &updateTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("arararasaka@corpo.neechan"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		&dao.IdentityModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1004), baseTime.Add(6*time.Hour), &updateTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "Saburo",
				LastName:  "Arasaka",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
		},
		&dao.ProfileModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1004), baseTime.Add(6*time.Hour), &updateTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug: "i-dont-have-any-ideas-anymore-alt",
			},
		},
	}

	data := []struct {
		name string

		fixtures []interface{}

		query  string
		offset int
		limit  int

		expect      []*dao.UserModel
		expectCount int
		expectErr   error
	}{
		{
			name:        "Success",
			query:       "Ele Be",
			limit:       10,
			expectCount: 3,
			expect: []*dao.UserModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("elon.bezos@space-origin.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "space-origin",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("real-elon-musk@tesla.tx"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Musk",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "real-elon",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("jeff.musk@amazon.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Jeff",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "big-brother",
						},
					},
				},
			},
		},
		{
			name:        "Success/AccurateQuery",
			query:       "Elon Bezos",
			limit:       10,
			expectCount: 3,
			expect: []*dao.UserModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("elon.bezos@space-origin.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "space-origin",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("jeff.musk@amazon.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Jeff",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "big-brother",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("real-elon-musk@tesla.tx"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Musk",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "real-elon",
						},
					},
				},
			},
		},
		{
			name:        "Success/AccurateQueryReverseOrder",
			query:       "Bezos Elon",
			limit:       10,
			expectCount: 3,
			expect: []*dao.UserModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("elon.bezos@space-origin.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "space-origin",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("jeff.musk@amazon.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Jeff",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "big-brother",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("real-elon-musk@tesla.tx"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Musk",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "real-elon",
						},
					},
				},
			},
		},
		{
			name:        "Success/NoQuery",
			limit:       10,
			expectCount: 5,
			expect: []*dao.UserModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1004), baseTime.Add(6*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("arararasaka@corpo.neechan"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Saburo",
							LastName:  "Arasaka",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexFemale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "i-dont-have-any-ideas-anymore-alt",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("elon.bezos@space-origin.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "space-origin",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1002), baseTime.Add(2*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("sarasaka@corpo.plaza"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Saburo",
							LastName:  "Arasaka",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexFemale,
						},
						Profile: dao.ProfileModelCore{
							Slug:     "i-dont-have-any-ideas-anymore",
							Username: "Mikoshi",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("jeff.musk@amazon.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Jeff",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "big-brother",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("real-elon-musk@tesla.tx"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Musk",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "real-elon",
						},
					},
				},
			},
		},
		{
			name:        "Success/LookUsernameOnlyIfGiven",
			query:       "Saburo",
			limit:       10,
			expectCount: 1,
			expect: []*dao.UserModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1004), baseTime.Add(6*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("arararasaka@corpo.neechan"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Saburo",
							LastName:  "Arasaka",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexFemale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "i-dont-have-any-ideas-anymore-alt",
						},
					},
				},
			},
		},
		{
			name:        "Success/IgnoreAccents",
			query:       "Élöñ Bēzos",
			limit:       10,
			expectCount: 3,
			expect: []*dao.UserModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("elon.bezos@space-origin.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "space-origin",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("jeff.musk@amazon.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Jeff",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "big-brother",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("real-elon-musk@tesla.tx"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Musk",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "real-elon",
						},
					},
				},
			},
		},
		{
			name:   "Success/NoRelevantResults",
			query:  "Abracadabra",
			limit:  10,
			expect: []*dao.UserModel(nil),
		},
		{
			name:        "Success/PaginationLimit",
			query:       "Elon Bezos",
			limit:       2,
			expectCount: 3,
			expect: []*dao.UserModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("elon.bezos@space-origin.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "space-origin",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("jeff.musk@amazon.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Jeff",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "big-brother",
						},
					},
				},
			},
		},
		{
			name:        "Success/PaginationOffset",
			query:       "Elon Bezos",
			offset:      2,
			limit:       10,
			expectCount: 3,

			expect: []*dao.UserModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("real-elon-musk@tesla.tx"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Musk",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "real-elon",
						},
					},
				},
			},
		},
		{
			name:        "Success/OnSlug",
			query:       "i-dont-have-any-ideas-anymore",
			limit:       10,
			expectCount: 2,
			expect: []*dao.UserModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1002), baseTime.Add(2*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("sarasaka@corpo.plaza"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Saburo",
							LastName:  "Arasaka",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexFemale,
						},
						Profile: dao.ProfileModelCore{
							Slug:     "i-dont-have-any-ideas-anymore",
							Username: "Mikoshi",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1004), baseTime.Add(6*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("arararasaka@corpo.neechan"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Saburo",
							LastName:  "Arasaka",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexFemale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "i-dont-have-any-ideas-anymore-alt",
						},
					},
				},
			},
		},
	}

	err := bunovel.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, count, err := dao.NewUserRepository(tx).Search(ctx, d.query, d.limit, d.offset)
				require.ErrorIs(t, err, d.expectErr)

				require.Empty(t, cmp.Diff(d.expect, res, cmpopts.IgnoreUnexported(time.Time{})))
				require.Equal(t, d.expectCount, count)
			})
		}
	})
	require.NoError(t, err)
}

func TestUserRepository_List(t *testing.T) {
	db, sqlDB := bunovel.GetTestPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []interface{}{
		// User 1
		&dao.CredentialsModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("elon.bezos@space-origin.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		&dao.IdentityModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "Elon",
				LastName:  "Bezos",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
		&dao.ProfileModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug: "space-origin",
			},
		},

		// User 2
		&dao.CredentialsModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("jeff.musk@amazon.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		&dao.IdentityModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "Jeff",
				LastName:  "Bezos",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
		&dao.ProfileModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1001), baseTime.Add(time.Hour), &updateTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug: "big-brother",
			},
		},

		// User 3
		&dao.CredentialsModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1002), baseTime.Add(2*time.Hour), &updateTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("sarasaka@corpo.plaza"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		&dao.IdentityModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1002), baseTime.Add(2*time.Hour), &updateTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "Saburo",
				LastName:  "Arasaka",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
		},
		&dao.ProfileModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1002), baseTime.Add(2*time.Hour), &updateTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug:     "i-dont-have-any-ideas-anymore",
				Username: "Mikoshi",
			},
		},

		// User 4
		&dao.CredentialsModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("real-elon-musk@tesla.tx"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		&dao.IdentityModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "Elon",
				LastName:  "Musk",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
		},
		&dao.ProfileModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1003), baseTime.Add(30*time.Minute), &updateTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug: "real-elon",
			},
		},

		// User 5
		&dao.CredentialsModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1004), baseTime.Add(6*time.Hour), &updateTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("arararasaka@corpo.neechan"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		&dao.IdentityModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1004), baseTime.Add(6*time.Hour), &updateTime),
			IdentityModelCore: dao.IdentityModelCore{
				FirstName: "Saburo",
				LastName:  "Arasaka",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexFemale,
			},
		},
		&dao.ProfileModel{
			Metadata: bunovel.NewMetadata(goframework.NumberUUID(1004), baseTime.Add(6*time.Hour), &updateTime),
			ProfileModelCore: dao.ProfileModelCore{
				Slug: "i-dont-have-any-ideas-anymore-alt",
			},
		},
	}

	data := []struct {
		name      string
		ids       []uuid.UUID
		fixtures  []interface{}
		expect    []*dao.UserModel
		expectErr error
	}{
		{
			name: "Success",
			ids: []uuid.UUID{
				goframework.NumberUUID(1000),
				goframework.NumberUUID(1002),
				goframework.NumberUUID(1004),
				// Don't exist.
				goframework.NumberUUID(15),
			},
			expect: []*dao.UserModel{
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1000), baseTime.Add(4*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("elon.bezos@space-origin.com"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Elon",
							LastName:  "Bezos",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexMale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "space-origin",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1002), baseTime.Add(2*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("sarasaka@corpo.plaza"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Saburo",
							LastName:  "Arasaka",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexFemale,
						},
						Profile: dao.ProfileModelCore{
							Slug:     "i-dont-have-any-ideas-anymore",
							Username: "Mikoshi",
						},
					},
				},
				{
					Metadata: bunovel.NewMetadata(goframework.NumberUUID(1004), baseTime.Add(6*time.Hour), &updateTime),
					UserModelCore: dao.UserModelCore{
						Credentials: dao.CredentialsModelCore{
							Email: MustParseEmail("arararasaka@corpo.neechan"),
						},
						Identity: dao.IdentityModelCore{
							FirstName: "Saburo",
							LastName:  "Arasaka",
							Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							Sex:       models.SexFemale,
						},
						Profile: dao.ProfileModelCore{
							Slug: "i-dont-have-any-ideas-anymore-alt",
						},
					},
				},
			},
		},
		{
			name: "Success/NoResults",
			ids: []uuid.UUID{
				// Don't exist.
				goframework.NumberUUID(15),
			},
			expect: nil,
		},
	}

	err := bunovel.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := dao.NewUserRepository(tx).List(ctx, d.ids)
				require.ErrorIs(t, err, d.expectErr)
				sort.Slice(res, func(i, j int) bool {
					return res[i].Profile.Slug < res[j].Profile.Slug
				})
				sort.Slice(d.expect, func(i, j int) bool {
					return d.expect[i].Profile.Slug < d.expect[j].Profile.Slug
				})

				require.Empty(t, cmp.Diff(d.expect, res, cmpopts.IgnoreUnexported(time.Time{})))
			})
		}
	})
	require.NoError(t, err)
}
