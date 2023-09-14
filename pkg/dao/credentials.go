package dao

import (
	"context"
	"github.com/a-novel/bunovel"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type CredentialsRepository interface {
	// GetCredentials reads a credentials object, based on a user id.
	GetCredentials(ctx context.Context, id uuid.UUID) (*CredentialsModel, error)
	// GetCredentialsByEmail reads a credentials object, based on a user email.
	// The match should be exact, thus requiring to pass a full Email object, rather than a string representation.
	// The Email.Validation field is ignored, and the CredentialsModelCore.NewEmail is not used for matching.
	GetCredentialsByEmail(ctx context.Context, email Email) (*CredentialsModel, error)
	// EmailExists looks if a given email is already used by a user as their main email (CredentialsModelCore.Email).
	// The match should be exact, thus requiring to pass a full Email object, rather than a string representation.
	// The Email.Validation field is ignored, and the CredentialsModelCore.NewEmail is not used for matching.
	EmailExists(ctx context.Context, email Email) (bool, error)

	// UpdateEmail updates the email of a user. The new email value is set as CredentialsModelCore.NewEmail.
	// The validation code should not be set on the Email.Validation field, as it will be filtered. The code value
	// MUST be hashed.
	// To make the new email the primary email of the user, you must call ValidateNewEmail.
	UpdateEmail(ctx context.Context, email Email, code string, id uuid.UUID, now time.Time) (*CredentialsModel, error)
	// ValidateEmail nullifies the Email.Validation value of CredentialsModelCore.Email, for the targeted user.
	ValidateEmail(ctx context.Context, id uuid.UUID, now time.Time) (*CredentialsModel, error)
	// ValidateNewEmail sets the email in argument as the primary email (CredentialsModelCore.Email) for the targeted user.
	// The CredentialsModelCore.NewEmail value is nullified in the process, and Email.Validation is filtered.
	ValidateNewEmail(ctx context.Context, id uuid.UUID, now time.Time) (*CredentialsModel, error)

	// UpdateEmailValidation sets a new Email.Validation code for the targeted user CredentialsModelCore.Email.
	// The code value MUST be hashed.
	UpdateEmailValidation(ctx context.Context, code string, id uuid.UUID, now time.Time) (*CredentialsModel, error)
	// UpdateNewEmailValidation sets a new Email.Validation code for the targeted user CredentialsModelCore.NewEmail.
	// The code value MUST be hashed. This method fails with sql.ErrNoRows if CredentialsModelCore.NewEmail contains an empty email
	// value.
	UpdateNewEmailValidation(ctx context.Context, code string, id uuid.UUID, now time.Time) (*CredentialsModel, error)
	// CancelNewEmail nullifies the CredentialsModelCore.NewEmail value for the targeted user. It does not fail if this field was
	// already empty.
	CancelNewEmail(ctx context.Context, id uuid.UUID, now time.Time) (*CredentialsModel, error)

	// UpdatePassword updates the password of the targeted user. The password value MUST be hashed in order to be
	// saved properly.
	UpdatePassword(ctx context.Context, newPassword string, id uuid.UUID, now time.Time) (*CredentialsModel, error)
	// ResetPassword sets Password.Validation field. The code value MUST be hashed. This does not nullify the
	// Password.Hashed field, so authentication can still work while password is being reset.
	ResetPassword(ctx context.Context, code string, email Email, now time.Time) (*CredentialsModel, error)
}

type CredentialsModel struct {
	bun.BaseModel `bun:"table:credentials"`
	bunovel.Metadata
	CredentialsModelCore
}

type CredentialsModelCore struct {
	// Email is the main email address of a user, used to authenticate it and to communicate with.
	Email Email `bun:"embed:email_"`
	// NewEmail is set when user wants to change its email address. Because email is the primary way to
	// authenticate a user, the Email value is not directly updated, but saved here in a pending state, with a
	// Email.Validation code set.
	// Once this email is validated, the Email field is updated, and this one is emptied.
	NewEmail Email `bun:"embed:new_email_"`
	// Password used to authenticate the user.
	Password Password `bun:"embed:password_"`
}

func NewCredentialsRepository(db bun.IDB) CredentialsRepository {
	return &credentialsRepositoryImpl{db: db}
}

type credentialsRepositoryImpl struct {
	db bun.IDB
}

func (repository *credentialsRepositoryImpl) GetCredentials(ctx context.Context, id uuid.UUID) (*CredentialsModel, error) {
	model := &CredentialsModel{Metadata: bunovel.NewMetadata(id, time.Time{}, nil)}

	if err := repository.db.NewSelect().Model(model).WherePK().Scan(ctx); err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	return model, nil
}

func (repository *credentialsRepositoryImpl) GetCredentialsByEmail(ctx context.Context, email Email) (*CredentialsModel, error) {
	model := new(CredentialsModel)

	if err := repository.db.NewSelect().Model(model).Where(WhereEmail("email", email)).Scan(ctx); err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	return model, nil
}

func (repository *credentialsRepositoryImpl) EmailExists(ctx context.Context, email Email) (bool, error) {
	ok, err := repository.db.NewSelect().Model(new(CredentialsModel)).Where(WhereEmail("email", email)).Exists(ctx)
	return ok, bunovel.HandlePGError(err)
}

func (repository *credentialsRepositoryImpl) UpdateEmail(ctx context.Context, email Email, code string, id uuid.UUID, now time.Time) (*CredentialsModel, error) {
	model := &CredentialsModel{
		Metadata: bunovel.NewMetadata(id, time.Time{}, &now),
		// Set new email with the given validation code. The main email remains unchanged until this email is
		// validated.
		CredentialsModelCore: CredentialsModelCore{
			NewEmail: Email{User: email.User, Domain: email.Domain, Validation: code},
		},
	}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		Column("new_email_user", "new_email_domain", "new_email_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	if err = bunovel.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *credentialsRepositoryImpl) ValidateEmail(ctx context.Context, id uuid.UUID, now time.Time) (*CredentialsModel, error) {
	model := &CredentialsModel{Metadata: bunovel.NewMetadata(id, time.Time{}, &now)}
	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		// User must have a pending email.
		Where("email_validation_code != ''").
		Column("email_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	if err = bunovel.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *credentialsRepositoryImpl) ValidateNewEmail(ctx context.Context, id uuid.UUID, now time.Time) (*CredentialsModel, error) {
	model := &CredentialsModel{Metadata: bunovel.NewMetadata(id, time.Time{}, nil)}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		// User must have a pending email update.
		Where("new_email_validation_code != ''").
		// Use the pending update ONLY to update the main email.
		SetColumn("email_user", "new_email_user").
		SetColumn("email_domain", "new_email_domain").
		SetColumn("email_validation_code", "''").
		// Empty the new_email columns, and update timestamps.
		SetColumn("new_email_user", "''").
		SetColumn("new_email_domain", "''").
		SetColumn("new_email_validation_code", "''").
		SetColumn("updated_at", "?", now).
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	if err = bunovel.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *credentialsRepositoryImpl) UpdatePassword(ctx context.Context, newPassword string, id uuid.UUID, now time.Time) (*CredentialsModel, error) {
	model := &CredentialsModel{
		Metadata: bunovel.NewMetadata(id, time.Time{}, &now),
		CredentialsModelCore: CredentialsModelCore{
			Password: Password{Hashed: newPassword},
		},
	}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		// "password_validation_code" is important to invalidate any pending reset, since a new known password is
		// now available.
		Column("password_hashed", "password_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	if err = bunovel.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *credentialsRepositoryImpl) ResetPassword(ctx context.Context, code string, email Email, now time.Time) (*CredentialsModel, error) {
	model := &CredentialsModel{
		Metadata: bunovel.NewMetadata(uuid.Nil, time.Time{}, &now),
		CredentialsModelCore: CredentialsModelCore{
			Password: Password{Validation: code},
		},
	}

	res, err := repository.db.NewUpdate().Model(model).
		Where(WhereEmail("email", email)).
		Column("password_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	if err = bunovel.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *credentialsRepositoryImpl) UpdateEmailValidation(ctx context.Context, code string, id uuid.UUID, now time.Time) (*CredentialsModel, error) {
	model := &CredentialsModel{
		Metadata: bunovel.NewMetadata(id, time.Time{}, &now),
		CredentialsModelCore: CredentialsModelCore{
			Email: Email{Validation: code},
		},
	}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		// User must have a pending validation update.
		Where("email_validation_code != ''").
		Column("email_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	if err = bunovel.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *credentialsRepositoryImpl) UpdateNewEmailValidation(ctx context.Context, code string, id uuid.UUID, now time.Time) (*CredentialsModel, error) {
	model := &CredentialsModel{
		Metadata: bunovel.NewMetadata(id, time.Time{}, &now),
		CredentialsModelCore: CredentialsModelCore{
			NewEmail: Email{Validation: code},
		},
	}

	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		// User must have a pending email update.
		Where("new_email_validation_code != ''").
		Column("new_email_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	if err = bunovel.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}

func (repository *credentialsRepositoryImpl) CancelNewEmail(ctx context.Context, id uuid.UUID, now time.Time) (*CredentialsModel, error) {
	model := &CredentialsModel{Metadata: bunovel.NewMetadata(id, time.Time{}, &now)}
	res, err := repository.db.NewUpdate().Model(model).
		WherePK().
		Column("new_email_user", "new_email_domain", "new_email_validation_code", "updated_at").
		Returning("*").
		Exec(ctx)

	if err != nil {
		return nil, bunovel.HandlePGError(err)
	}

	if err = bunovel.ForceRowsUpdate(res); err != nil {
		return nil, err
	}

	return model, nil
}
