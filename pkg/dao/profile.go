package dao

import (
	"context"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/postgresql"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type ProfileRepository interface {
	// GetProfile reads a profile object, based on a user id.
	GetProfile(ctx context.Context, id uuid.UUID) (*ProfileModel, error)
	// GetProfileBySlug reads a profile object, based on a slug.
	GetProfileBySlug(ctx context.Context, slug string) (*ProfileModel, error)
	// SlugExists looks if a given slug is already used by another profile.
	SlugExists(ctx context.Context, slug string) (bool, error)

	// Update the slug of the targeted user.
	Update(ctx context.Context, data *ProfileModelCore, id uuid.UUID, now time.Time) (*ProfileModel, error)
}

type ProfileModel struct {
	bun.BaseModel `bun:"table:profiles"`

	postgresql.Metadata
	ProfileModelCore
}

type ProfileModelCore struct {
	// Username is a fake name displayed for a user.
	Username string `bun:"username"`
	// Slug is the unique url suffix used to access the current profile.
	Slug string `bun:"slug"`
}

func NewProfileRepository(db bun.IDB) ProfileRepository {
	return &profileRepositoryImpl{db: db}
}

type profileRepositoryImpl struct {
	db bun.IDB
}

func (repository *profileRepositoryImpl) GetProfile(ctx context.Context, id uuid.UUID) (*ProfileModel, error) {
	model := &ProfileModel{Metadata: postgresql.NewMetadata(id, time.Time{}, nil)}

	if err := repository.db.NewSelect().Model(model).WherePK().Scan(ctx); err != nil {
		return nil, errors.HandlePGError(err)
	}

	return model, nil
}

func (repository *profileRepositoryImpl) GetProfileBySlug(ctx context.Context, slug string) (*ProfileModel, error) {
	model := new(ProfileModel)

	if err := repository.db.NewSelect().Model(model).Where("slug = ?", slug).Scan(ctx); err != nil {
		return nil, errors.HandlePGError(err)
	}

	return model, nil
}

func (repository *profileRepositoryImpl) SlugExists(ctx context.Context, slug string) (bool, error) {
	ok, err := repository.db.NewSelect().Model(new(ProfileModel)).Where("slug = ?", slug).Exists(ctx)
	return ok, errors.HandlePGError(err)
}

func (repository *profileRepositoryImpl) Update(ctx context.Context, data *ProfileModelCore, id uuid.UUID, now time.Time) (*ProfileModel, error) {
	model := &ProfileModel{Metadata: postgresql.NewMetadata(id, time.Time{}, &now), ProfileModelCore: *data}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		Column(
			"slug",
			"username",
			"updated_at",
		).
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, errors.HandlePGError(err)
	}

	if err := errors.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}
