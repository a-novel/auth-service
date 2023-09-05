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

func TestCredentialsRepository_GetCredentials(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.CredentialsModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user@domain.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
	}

	data := []struct {
		name string

		id uuid.UUID

		expect    *dao.CredentialsModel
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
		repository := dao.NewCredentialsRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.GetCredentials(ctx, d.id)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_GetCredentialsByEmail(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.CredentialsModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user@domain.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
	}

	data := []struct {
		name string

		email dao.Email

		expect    *dao.CredentialsModel
		expectErr error
	}{
		{
			name:   "Success",
			email:  MustParseEmail("user@domain.com"),
			expect: fixtures[0],
		},
		{
			name:      "Error/NotFound",
			email:     MustParseEmail("fake-user@domain.com"),
			expectErr: errors.ErrNotFound,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := dao.NewCredentialsRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.GetCredentialsByEmail(ctx, d.email)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_EmailExists(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.CredentialsModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user@domain.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
	}

	data := []struct {
		name string

		email dao.Email

		expect    bool
		expectErr error
	}{
		{
			name:   "Success/Exists",
			email:  MustParseEmail("user@domain.com"),
			expect: true,
		},
		{
			name:   "Success/DoesNotExists",
			email:  MustParseEmail("fake-user@domain.com"),
			expect: false,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		repository := dao.NewCredentialsRepository(tx)

		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				res, err := repository.EmailExists(ctx, d.email)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_UpdateEmail(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.CredentialsModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user1@domain.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmailWithValidation("user2@domain.com", "initial-validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1002), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user3@domain.com"),
				NewEmail: MustParseEmailWithValidation("new-other-user3@domain.com", "other-validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
	}

	data := []struct {
		name string

		email dao.Email
		code  string
		id    uuid.UUID
		now   time.Time

		expect    *dao.CredentialsModel
		expectErr error
	}{
		{
			name:  "Success",
			email: MustParseEmailWithValidation("new-user1@domain.com", "this-code-should-be-ignored"),
			code:  "validation-code",
			id:    test.NumberUUID(1000),
			now:   updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					NewEmail: MustParseEmailWithValidation("new-user1@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name:  "Success/WithValidationOnMainEmail",
			email: MustParseEmailWithValidation("new-user2@domain.com", "this-code-should-be-ignored"),
			code:  "validation-code",
			id:    test.NumberUUID(1001),
			now:   updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmailWithValidation("user2@domain.com", "initial-validation-code"),
					NewEmail: MustParseEmailWithValidation("new-user2@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name:  "Success/WithPreviousNewEmail",
			email: MustParseEmailWithValidation("new-user3@domain.com", "this-code-should-be-ignored"),
			code:  "validation-code",
			id:    test.NumberUUID(1002),
			now:   updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1002), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user3@domain.com"),
					NewEmail: MustParseEmailWithValidation("new-user3@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name:  "Success/WhenAnotherAccountHasTheSameEmailPendingValidation",
			email: MustParseEmailWithValidation("new-other-user3@domain.com", "this-code-should-be-ignored"),
			code:  "validation-code",
			id:    test.NumberUUID(1000),
			now:   updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					NewEmail: MustParseEmailWithValidation("new-other-user3@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		// This is allowed, validating the email will however fail. It should be blocked in the service, by checking
		// if the email is already taken.
		{
			name:  "Success/TakenByAnotherAccount",
			email: MustParseEmail("user2@domain.com"),
			code:  "validation-code",
			id:    test.NumberUUID(1000),
			now:   updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					NewEmail: MustParseEmailWithValidation("user2@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name:      "Error/NotFound",
			email:     MustParseEmailWithValidation("new-user1@domain.com", "this-code-should-be-ignored"),
			code:      "validation-code",
			id:        test.NumberUUID(100),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
		{
			name:      "Error/WithoutValidationCode",
			email:     MustParseEmailWithValidation("new-user1@domain.com", "this-code-should-be-ignored"),
			id:        test.NumberUUID(1000),
			now:       updateTime,
			expectErr: errors.ErrConstraintViolation,
		},
		{
			name: "Error/WithoutUser",
			email: dao.Email{
				Domain: "domain.com",
			},
			code:      "validation-code",
			id:        test.NumberUUID(1000),
			now:       updateTime,
			expectErr: errors.ErrConstraintViolation,
		},
		{
			name: "Error/WithoutDomain",
			email: dao.Email{
				User: "new-user1",
			},
			code:      "validation-code",
			id:        test.NumberUUID(1000),
			now:       updateTime,
			expectErr: errors.ErrConstraintViolation,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := dao.NewCredentialsRepository(stx).UpdateEmail(ctx, d.email, d.code, d.id, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_ValidateEmail(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.CredentialsModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmailWithValidation("user1@domain.com", "initial-validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmailWithValidation("user2@domain.com", "initial-validation-code"),
				NewEmail: MustParseEmailWithValidation("new-user2@domain.com", "validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1002), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user3@domain.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
	}

	data := []struct {
		name string

		id  uuid.UUID
		now time.Time

		expect    *dao.CredentialsModel
		expectErr error
	}{
		{
			name: "Success",
			id:   test.NumberUUID(1000),
			now:  updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name: "Success/WithEmailPendingValidation",
			id:   test.NumberUUID(1001),
			now:  updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user2@domain.com"),
					NewEmail: MustParseEmailWithValidation("new-user2@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name:      "Error/NoPendingValidation",
			id:        test.NumberUUID(1002),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
		{
			name:      "Error/NotFound",
			id:        test.NumberUUID(1),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := dao.NewCredentialsRepository(stx).ValidateEmail(ctx, d.id, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_ValidateNewEmail(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.CredentialsModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user1@domain.com"),
				NewEmail: MustParseEmailWithValidation("new-user1@domain.com", "validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmailWithValidation("user2@domain.com", "initial-validation-code"),
				NewEmail: MustParseEmailWithValidation("new-user2@domain.com", "validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1002), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user3@domain.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1003), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user4@domain.com"),
				NewEmail: MustParseEmailWithValidation("user1@domain.com", "validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
	}

	data := []struct {
		name string

		id  uuid.UUID
		now time.Time

		expect    *dao.CredentialsModel
		expectErr error
	}{
		{
			name: "Success",
			id:   test.NumberUUID(1000),
			now:  updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("new-user1@domain.com"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name: "Success/WithMainEmailPendingValidation",
			id:   test.NumberUUID(1001),
			now:  updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("new-user2@domain.com"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name:      "Error/NoPendingValidation",
			id:        test.NumberUUID(1002),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
		{
			name:      "Error/NotFound",
			id:        test.NumberUUID(1),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
		{
			name:      "Error/AlreadyTaken",
			id:        test.NumberUUID(1003),
			now:       updateTime,
			expectErr: errors.ErrUniqConstraintViolation,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := dao.NewCredentialsRepository(stx).ValidateNewEmail(ctx, d.id, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_UpdatePassword(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.CredentialsModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user1@domain.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user2@domain.com"),
				Password: dao.Password{Hashed: "password-hashed", Validation: "validation-code"},
			},
		},
	}

	data := []struct {
		name string

		newPassword string
		id          uuid.UUID
		now         time.Time

		expect    *dao.CredentialsModel
		expectErr error
	}{
		{
			name:        "Success",
			newPassword: "new-password-hashed",
			id:          test.NumberUUID(1000),
			now:         updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					Password: dao.Password{Hashed: "new-password-hashed"},
				},
			},
		},
		{
			name:        "Success/WithPendingReset",
			newPassword: "new-password-hashed",
			id:          test.NumberUUID(1001),
			now:         updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user2@domain.com"),
					Password: dao.Password{Hashed: "new-password-hashed"},
				},
			},
		},
		{
			name:        "Error/NotFound",
			newPassword: "new-password-hashed",
			id:          test.NumberUUID(100),
			now:         updateTime,
			expectErr:   errors.ErrNotFound,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := dao.NewCredentialsRepository(stx).UpdatePassword(ctx, d.newPassword, d.id, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_ResetPassword(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.CredentialsModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user1@domain.com"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user2@domain.com"),
				NewEmail: MustParseEmailWithValidation("new-user2@domain.com", "validation-code"),
				Password: dao.Password{Hashed: "password-hashed", Validation: "old-validation-code"},
			},
		},
	}

	data := []struct {
		name string

		code  string
		email dao.Email
		now   time.Time

		expect    *dao.CredentialsModel
		expectErr error
	}{
		{
			name:  "Success",
			code:  "validation-code",
			email: MustParseEmail("user1@domain.com"),
			now:   updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					Password: dao.Password{Hashed: "password-hashed", Validation: "validation-code"},
				},
			},
		},
		{
			name:  "Success/WithPreviousResetPending",
			code:  "validation-code",
			email: MustParseEmail("user2@domain.com"),
			now:   updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user2@domain.com"),
					NewEmail: MustParseEmailWithValidation("new-user2@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed", Validation: "validation-code"},
				},
			},
		},
		{
			name:      "Error/NotFound",
			code:      "validation-code",
			email:     MustParseEmail("user3@domain.com"),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
		{
			name:      "Error/NotFoundIfEmailPendingUpdate",
			code:      "validation-code",
			email:     MustParseEmail("new-user2@domain.com"),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := dao.NewCredentialsRepository(stx).ResetPassword(ctx, d.code, d.email, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_UpdateEmailValidation(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.CredentialsModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmailWithValidation("user1@domain.com", "old-validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmailWithValidation("user2@domain.com", "old-validation-code"),
				NewEmail: MustParseEmailWithValidation("new-user2@domain.com", "validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1002), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user3@domain.com"),
				NewEmail: MustParseEmailWithValidation("new-user3@domain.com", "validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
	}

	data := []struct {
		name string

		code string
		id   uuid.UUID
		now  time.Time

		expect    *dao.CredentialsModel
		expectErr error
	}{
		{
			name: "Success",
			code: "validation-code",
			id:   test.NumberUUID(1000),
			now:  updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmailWithValidation("user1@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name: "Success/WithAnEmailPendingUpdate",
			code: "validation-code",
			id:   test.NumberUUID(1001),
			now:  updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmailWithValidation("user2@domain.com", "validation-code"),
					NewEmail: MustParseEmailWithValidation("new-user2@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name:      "Error/NotFound",
			code:      "validation-code",
			id:        test.NumberUUID(100),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
		{
			name:      "Error/AlreadyValidated",
			code:      "validation-code",
			id:        test.NumberUUID(1002),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := dao.NewCredentialsRepository(stx).UpdateEmailValidation(ctx, d.code, d.id, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_UpdateNewEmailValidation(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.CredentialsModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user1@domain.com"),
				NewEmail: MustParseEmailWithValidation("new-user1@domain.com", "old-validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmailWithValidation("user2@domain.com", "initial-validation-code"),
				NewEmail: MustParseEmailWithValidation("new-user2@domain.com", "old-validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1002), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmailWithValidation("user3@domain.com", "initial-validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
	}

	data := []struct {
		name string

		code string
		id   uuid.UUID
		now  time.Time

		expect    *dao.CredentialsModel
		expectErr error
	}{
		{
			name:      "Error/NoValidationCode",
			id:        test.NumberUUID(1000),
			now:       updateTime,
			expectErr: errors.ErrConstraintViolation,
		},
		{
			name: "Success",
			code: "validation-code",
			id:   test.NumberUUID(1000),
			now:  updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					NewEmail: MustParseEmailWithValidation("new-user1@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name: "Success/WithMainEmailPendingValidation",
			code: "validation-code",
			id:   test.NumberUUID(1001),
			now:  updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmailWithValidation("user2@domain.com", "initial-validation-code"),
					NewEmail: MustParseEmailWithValidation("new-user2@domain.com", "validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name:      "Error/NotFound",
			code:      "validation-code",
			id:        test.NumberUUID(100),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
		{
			name:      "Error/NoPendingUpdate",
			code:      "validation-code",
			id:        test.NumberUUID(1002),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := dao.NewCredentialsRepository(stx).UpdateNewEmailValidation(ctx, d.code, d.id, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}

func TestCredentialsRepository_CancelNewEmail(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{migrations.Migrations})
	defer db.Close()
	defer sqlDB.Close()

	fixtures := []*dao.CredentialsModel{
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmail("user1@domain.com"),
				NewEmail: MustParseEmailWithValidation("new-user1@domain.com", "old-validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmailWithValidation("user2@domain.com", "initial-validation-code"),
				NewEmail: MustParseEmailWithValidation("new-user2@domain.com", "old-validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
		{
			Metadata: postgresql.NewMetadata(test.NumberUUID(1002), baseTime, &baseTime),
			CredentialsModelCore: dao.CredentialsModelCore{
				Email:    MustParseEmailWithValidation("user3@domain.com", "initial-validation-code"),
				Password: dao.Password{Hashed: "password-hashed"},
			},
		},
	}

	data := []struct {
		name string

		id  uuid.UUID
		now time.Time

		expect    *dao.CredentialsModel
		expectErr error
	}{
		{
			name: "Success",
			id:   test.NumberUUID(1000),
			now:  updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1000), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmail("user1@domain.com"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name: "Success/WithMainEmailPendingValidation",
			id:   test.NumberUUID(1001),
			now:  updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1001), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmailWithValidation("user2@domain.com", "initial-validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name: "Success/NoPendingUpdate",
			id:   test.NumberUUID(1002),
			now:  updateTime,
			expect: &dao.CredentialsModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1002), baseTime, &updateTime),
				CredentialsModelCore: dao.CredentialsModelCore{
					Email:    MustParseEmailWithValidation("user3@domain.com", "initial-validation-code"),
					Password: dao.Password{Hashed: "password-hashed"},
				},
			},
		},
		{
			name:      "Error/NotFound",
			id:        test.NumberUUID(100),
			now:       updateTime,
			expectErr: errors.ErrNotFound,
		},
	}

	err := test.RunTransactionalTest(db, fixtures, func(ctx context.Context, tx bun.Tx) {
		for _, d := range data {
			t.Run(d.name, func(st *testing.T) {
				stx, err := tx.BeginTx(ctx, nil)
				require.NoError(st, err)
				defer stx.Rollback()

				res, err := dao.NewCredentialsRepository(stx).CancelNewEmail(ctx, d.id, d.now)
				require.ErrorIs(t, err, d.expectErr)
				require.Equal(t, d.expect, res)
			})
		}
	})
	require.NoError(t, err)
}
