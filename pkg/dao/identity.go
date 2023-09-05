package dao

import (
	"context"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/postgresql"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type IdentityRepository interface {
	// GetIdentity reads an identity object, based on a user id.
	GetIdentity(ctx context.Context, id uuid.UUID) (*IdentityModel, error)
	// Update the identity of the targeted user.
	Update(ctx context.Context, data *IdentityModelCore, id uuid.UUID, now time.Time) (*IdentityModel, error)
}

type IdentityModel struct {
	bun.BaseModel `bun:"table:identities"`

	postgresql.Metadata
	IdentityModelCore
}

// IdentityModelCore contains the explicitly editable data of the IdentityModel.
type IdentityModelCore struct {
	FirstName string     `json:"firstName" bun:"first_name"`
	LastName  string     `json:"lastName" bun:"last_name"`
	Birthday  time.Time  `json:"birthday" bun:"birthday"`
	Sex       models.Sex `json:"sex" bun:"sex,notnull"`
}

func NewIdentityRepository(db bun.IDB) IdentityRepository {
	return &identityRepositoryImpl{db: db}
}

type identityRepositoryImpl struct {
	db bun.IDB
}

func (repository *identityRepositoryImpl) GetIdentity(ctx context.Context, id uuid.UUID) (*IdentityModel, error) {
	model := &IdentityModel{Metadata: postgresql.NewMetadata(id, time.Time{}, nil)}

	if err := repository.db.NewSelect().Model(model).WherePK().Scan(ctx); err != nil {
		return nil, errors.HandlePGError(err)
	}

	return model, nil
}

func (repository *identityRepositoryImpl) Update(ctx context.Context, data *IdentityModelCore, id uuid.UUID, now time.Time) (*IdentityModel, error) {
	model := &IdentityModel{Metadata: postgresql.NewMetadata(id, time.Time{}, &now), IdentityModelCore: *data}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		Column("first_name", "last_name", "birthday", "sex", "updated_at").
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
