package services_test

import (
	"context"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/mailer"
	"github.com/a-novel/go-framework/postgresql"
	"github.com/a-novel/go-framework/test"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	data := []struct {
		name string

		form                  models.RegisterForm
		now                   time.Time
		validateEmailTemplate string
		validateEmailLink     string

		publicValidationCode      string
		privateValidationCode     string
		generateValidationCodeErr error

		shouldCallEmailExists bool
		emailExists           bool
		emailExistsErr        error

		shouldCallSlugExists bool
		slugExists           bool
		slugExistsErr        error

		shouldCallGenerateToken bool
		generateTokenStatus     *models.UserTokenStatus
		generateTokenErr        error

		shouldCallCreateUser bool
		createUser           *dao.UserModel
		createUserErr        error

		shouldCallMailer          bool
		shouldCallMailerWithEmail *mail.Email
		shouldCallMailerWithData  map[string]interface{}
		mailerErr                 error

		expect            *models.UserTokenStatus
		expectErr         error
		expectDeferred    bool
		expectDeferredErr error
	}{
		{
			name: "Success",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear), // 20 Yo
				Slug:      "slug",
			},
			now:                     baseTime,
			validateEmailTemplate:   "validate-email-template",
			validateEmailLink:       "validate-email-link",
			publicValidationCode:    "public-validation-code",
			privateValidationCode:   "private-validation-code",
			shouldCallEmailExists:   true,
			emailExists:             false,
			shouldCallSlugExists:    true,
			slugExists:              false,
			shouldCallGenerateToken: true,
			generateTokenStatus: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallCreateUser: true,
			createUser: &dao.UserModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, &baseTime),
				UserModelCore: dao.UserModelCore{
					Credentials: dao.CredentialsModelCore{
						Email:    dao.Email{User: "user", Domain: "domain.com", Validation: "private-validation-code"},
						Password: dao.Password{Hashed: "password"},
					},
					Identity: dao.IdentityModelCore{
						FirstName: "name",
						LastName:  "last-name",
						Sex:       models.SexMale,
						Birthday:  baseTime.Add(-20 * timeYear),
					},
					Profile: dao.ProfileModelCore{
						Slug: "slug",
					},
				},
			},
			shouldCallMailer:          true,
			shouldCallMailerWithEmail: mail.NewEmail("name", "user@domain.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "name",
				"validation_link": "validate-email-link?id=01010101-0101-0101-0101-010101010101&code=public-validation-code",
			},
			expect: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			expectDeferred: true,
		},
		{
			name: "Success/WithUsername",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexFemale,
				Birthday:  baseTime.Add(-20 * timeYear), // 20 Yo
				Slug:      "slug",
				Username:  "my username",
			},
			now:                     baseTime,
			validateEmailTemplate:   "validate-email-template",
			validateEmailLink:       "validate-email-link",
			publicValidationCode:    "public-validation-code",
			privateValidationCode:   "private-validation-code",
			shouldCallEmailExists:   true,
			emailExists:             false,
			shouldCallSlugExists:    true,
			slugExists:              false,
			shouldCallGenerateToken: true,
			generateTokenStatus: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallCreateUser: true,
			createUser: &dao.UserModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, &baseTime),
				UserModelCore: dao.UserModelCore{
					Credentials: dao.CredentialsModelCore{
						Email:    dao.Email{User: "user", Domain: "domain.com", Validation: "private-validation-code"},
						Password: dao.Password{Hashed: "password"},
					},
					Identity: dao.IdentityModelCore{
						FirstName: "name",
						LastName:  "last-name",
						Sex:       models.SexMale,
						Birthday:  baseTime.Add(-20 * timeYear),
					},
					Profile: dao.ProfileModelCore{
						Slug:     "slug",
						Username: "username",
					},
				},
			},
			shouldCallMailer:          true,
			shouldCallMailerWithEmail: mail.NewEmail("name", "user@domain.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "name",
				"validation_link": "validate-email-link?id=01010101-0101-0101-0101-010101010101&code=public-validation-code",
			},
			expect: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			expectDeferred: true,
		},
		{
			name: "Error/EmailSendingFailure",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear), // 20 Yo
				Slug:      "slug",
			},
			now:                     baseTime,
			validateEmailTemplate:   "validate-email-template",
			validateEmailLink:       "validate-email-link",
			publicValidationCode:    "public-validation-code",
			privateValidationCode:   "private-validation-code",
			shouldCallEmailExists:   true,
			emailExists:             false,
			shouldCallSlugExists:    true,
			slugExists:              false,
			shouldCallGenerateToken: true,
			generateTokenStatus: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallCreateUser: true,
			createUser: &dao.UserModel{
				Metadata: postgresql.NewMetadata(test.NumberUUID(1), baseTime, &baseTime),
				UserModelCore: dao.UserModelCore{
					Credentials: dao.CredentialsModelCore{
						Email:    dao.Email{User: "user", Domain: "domain.com", Validation: "private-validation-code"},
						Password: dao.Password{Hashed: "password"},
					},
					Identity: dao.IdentityModelCore{
						FirstName: "name",
						LastName:  "last-name",
						Sex:       models.SexMale,
						Birthday:  baseTime.Add(-20 * timeYear),
					},
					Profile: dao.ProfileModelCore{
						Slug: "slug",
					},
				},
			},
			shouldCallMailer:          true,
			shouldCallMailerWithEmail: mail.NewEmail("name", "user@domain.com"),
			shouldCallMailerWithData: map[string]interface{}{
				"name":            "name",
				"validation_link": "validate-email-link?id=01010101-0101-0101-0101-010101010101&code=public-validation-code",
			},
			mailerErr: fooErr,
			expect: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			expectDeferred:    true,
			expectDeferredErr: fooErr,
		},
		{
			name: "Error/CreateUserFailure",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear), // 20 Yo
				Slug:      "slug",
			},
			now:                     baseTime,
			validateEmailTemplate:   "validate-email-template",
			validateEmailLink:       "validate-email-link",
			publicValidationCode:    "public-validation-code",
			privateValidationCode:   "private-validation-code",
			shouldCallEmailExists:   true,
			emailExists:             false,
			shouldCallSlugExists:    true,
			slugExists:              false,
			shouldCallGenerateToken: true,
			generateTokenStatus: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
			},
			shouldCallCreateUser: true,
			createUserErr:        fooErr,
			expectErr:            fooErr,
		},
		{
			name: "Error/GenerateTokenFailure",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear), // 20 Yo
				Slug:      "slug",
			},
			now:                     baseTime,
			validateEmailTemplate:   "validate-email-template",
			validateEmailLink:       "validate-email-link",
			publicValidationCode:    "public-validation-code",
			privateValidationCode:   "private-validation-code",
			shouldCallEmailExists:   true,
			emailExists:             false,
			shouldCallSlugExists:    true,
			slugExists:              false,
			shouldCallGenerateToken: true,
			generateTokenErr:        fooErr,
			expectErr:               fooErr,
		},
		{
			name: "Error/GenerateValidationCodeFailure",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear), // 20 Yo
				Slug:      "slug",
			},
			now:                       baseTime,
			validateEmailTemplate:     "validate-email-template",
			validateEmailLink:         "validate-email-link",
			generateValidationCodeErr: fooErr,
			shouldCallEmailExists:     true,
			emailExists:               false,
			shouldCallSlugExists:      true,
			slugExists:                false,
			expectErr:                 fooErr,
		},
		{
			name: "Error/SlugExists",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear), // 20 Yo
				Slug:      "slug",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			shouldCallEmailExists: true,
			emailExists:           false,
			shouldCallSlugExists:  true,
			slugExists:            true,
			expectErr:             services.ErrTaken,
		},
		{
			name: "Error/SlugCheckFailure",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear), // 20 Yo
				Slug:      "slug",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			shouldCallEmailExists: true,
			emailExists:           false,
			shouldCallSlugExists:  true,
			slugExistsErr:         fooErr,
			expectErr:             fooErr,
		},
		{
			name: "Error/EmailExists",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear), // 20 Yo
				Slug:      "slug",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			shouldCallEmailExists: true,
			emailExists:           true,
			expectErr:             services.ErrTaken,
		},
		{
			name: "Error/EmailCheckFailure",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear), // 20 Yo
				Slug:      "slug",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			shouldCallEmailExists: true,
			emailExistsErr:        fooErr,
			expectErr:             fooErr,
		},
		{
			name: "Error/UserTooYoung",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add((services.MinAge - 1) * timeYear),
				Slug:      "slug",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/UserNotBorn(Seriously)",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(2 * timeYear),
				Slug:      "slug",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/UserTooOld",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add((services.MaxAge + 1) * timeYear),
				Slug:      "slug",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/EmailInvalid",
			form: models.RegisterForm{
				Email:     "userdomain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/UsernameInvalid",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "😊😊😊",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/UsernameTooLong",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  strings.Repeat("a", services.MaxUsernameLength+1),
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidLastName",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "$%&/()=?",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/LastNameTooLong",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  strings.Repeat("a", services.MaxNameLength+1),
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidFirstName",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "$%&/()=?",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/FirstNameTooLong",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: strings.Repeat("a", services.MaxNameLength+1),
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidSlug",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "Sl#ug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/SlugTooLong",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      strings.Repeat("a", services.MaxSlugLength+1),
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/PasswordTooShort",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "p",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/PasswordTooLong",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  strings.Repeat("a", services.MaxPasswordLength+1),
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/EmailTooLong",
			form: models.RegisterForm{
				Email:     strings.Repeat("a", services.MaxEmailLength+1) + "@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/InvalidSex",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.Sex("non-binary"),
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/NoEmail",
			form: models.RegisterForm{
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/NoPassword",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/NoFirstName",
			form: models.RegisterForm{
				Email:    "user@domain.com",
				Password: "password",
				LastName: "last-name",
				Sex:      models.SexMale,
				Birthday: baseTime.Add(-20 * timeYear),
				Slug:     "slug",
				Username: "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/NoLastName",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Slug:      "slug",
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
		{
			name: "Error/NoSlug",
			form: models.RegisterForm{
				Email:     "user@domain.com",
				Password:  "password",
				FirstName: "name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime.Add(-20 * timeYear),
				Username:  "username",
			},
			now:                   baseTime,
			validateEmailTemplate: "validate-email-template",
			validateEmailLink:     "validate-email-link",
			expectErr:             errors.ErrInvalidEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			credentialsDAO := daomocks.NewCredentialsRepository(t)
			profileDAO := daomocks.NewProfileRepository(t)
			userDAO := daomocks.NewUserRepository(t)
			mailerService := mailer.NewMockMailer(t)
			generateTokenService := servicesmocks.NewGenerateTokenService(t)

			generateLink := func() (string, string, error) {
				return d.publicValidationCode, d.privateValidationCode, d.generateValidationCodeErr
			}

			if d.shouldCallMailer {
				mailerService.
					On("Send", d.shouldCallMailerWithEmail, d.validateEmailTemplate, d.shouldCallMailerWithData).
					Return(d.mailerErr)
			}

			if d.shouldCallEmailExists {
				credentialsDAO.
					On("EmailExists", context.Background(), mock.Anything).
					Return(d.emailExists, d.emailExistsErr)
			}

			if d.shouldCallSlugExists {
				profileDAO.
					On("SlugExists", context.Background(), d.form.Slug).
					Return(d.slugExists, d.slugExistsErr)
			}

			if d.shouldCallGenerateToken {
				generateTokenService.
					On("GenerateToken", context.Background(), mock.Anything, mock.Anything, d.now).
					Return(d.generateTokenStatus, d.generateTokenErr)
			}

			if d.shouldCallCreateUser {
				userDAO.
					On("Create", context.Background(), mock.Anything, mock.Anything, d.now).
					Return(d.createUser, d.createUserErr)
			}

			service := services.NewRegisterService(credentialsDAO, profileDAO, userDAO, mailerService, generateLink, generateTokenService, d.validateEmailTemplate, d.validateEmailLink)
			res, deferred, err := service.Register(context.Background(), d.form, d.now)

			require.ErrorIs(t, err, d.expectErr)
			require.Equal(t, d.expect, res)

			if d.expectDeferred {
				require.NotNil(t, deferred)
				require.ErrorIs(t, deferred(), d.expectDeferredErr)
			} else {
				require.Nil(t, deferred)
			}

			credentialsDAO.AssertExpectations(t)
			profileDAO.AssertExpectations(t)
			userDAO.AssertExpectations(t)
			mailerService.AssertExpectations(t)
			generateTokenService.AssertExpectations(t)
		})
	}
}
