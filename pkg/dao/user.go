package dao

import (
	"context"
	"github.com/a-novel/bunovel"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type UserRepository interface {
	// Create creates a new user. The credentials, identity and profile objects will share the same ID and create time.
	// If any error occurs, no data is created.
	Create(ctx context.Context, data *UserModelCore, id uuid.UUID, now time.Time) (*UserModel, error)
	// Search performs a cross-table search query over the user repository.
	Search(ctx context.Context, query string, limit, offset int) ([]*UserModel, int, error)
	// List returns a list of users
	List(ctx context.Context, ids []uuid.UUID) ([]*UserModel, error)
}

type UserModel struct {
	bun.BaseModel `bun:"table:users_view"`
	bunovel.Metadata
	UserModelCore
}

type UserModelCore struct {
	Credentials CredentialsModelCore `bun:"credentials"`
	Identity    IdentityModelCore    `bun:"identity"`
	Profile     ProfileModelCore     `bun:"profile"`
}

func NewUserRepository(db bun.IDB) UserRepository {
	return &userRepositoryImpl{db: db}
}

type userRepositoryImpl struct {
	db bun.IDB
}

func (repository *userRepositoryImpl) Create(ctx context.Context, data *UserModelCore, id uuid.UUID, now time.Time) (*UserModel, error) {
	model := new(UserModel)

	// Create all in a transaction, to avoid partially created users if any part of the operation fails.
	err := repository.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		credentialsModel := &CredentialsModel{Metadata: bunovel.NewMetadata(id, now, nil), CredentialsModelCore: data.Credentials}
		_, err := repository.db.NewInsert().Model(credentialsModel).Exec(ctx)
		if err != nil {
			return err
		}

		identityModel := &IdentityModel{Metadata: bunovel.NewMetadata(id, now, nil), IdentityModelCore: data.Identity}
		_, err = repository.db.NewInsert().Model(identityModel).Exec(ctx)
		if err != nil {
			return err
		}

		profileModel := &ProfileModel{Metadata: bunovel.NewMetadata(id, now, nil), ProfileModelCore: data.Profile}
		_, err = repository.db.NewInsert().Model(profileModel).Exec(ctx)
		if err != nil {
			return err
		}

		// Assign common metadata.
		model.ID = id
		model.CreatedAt = now

		model.Credentials = credentialsModel.CredentialsModelCore
		model.Identity = identityModel.IdentityModelCore
		model.Profile = profileModel.ProfileModelCore
		return nil
	})

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	return model, nil
}

func (repository *userRepositoryImpl) Search(ctx context.Context, query string, limit, offset int) ([]*UserModel, int, error) {
	var results []*UserModel

	count, err := repository.db.NewSelect().Model(&results).
		Join(`
LEFT JOIN LATERAL (
	SELECT GREATEST(
		search_field(?0, COALESCE(NULLIF(profile ->> 'username', ''), (identity ->> 'firstName') || ' ' || (identity ->> 'lastName'))),
		search_field(?0, profile ->> 'slug')
	) AS score
) AS proximity ON TRUE`, query).
		Where("proximity.score > 0.1").
		Order("proximity.score DESC", "created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx, &results)

	if err != nil {
		return nil, 0, bunovel.HandlePGError(err)
	}

	return results, count, nil
}

func (repository *userRepositoryImpl) List(ctx context.Context, ids []uuid.UUID) ([]*UserModel, error) {
	var results []*UserModel

	err := repository.db.NewSelect().Model(&results).Where("id IN (?)", bun.In(ids)).Scan(ctx)
	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	return results, nil
}
