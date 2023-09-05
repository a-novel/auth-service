package services

import (
	"context"
	goerrors "errors"
	"fmt"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/mailer"
	"github.com/a-novel/go-framework/validation"
	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type RegisterService interface {
	// Register creates a new user in the database, and return an initial token for it.
	Register(ctx context.Context, form models.RegisterForm, now time.Time) (*models.UserTokenStatus, func() error, error)
}

func NewRegisterService(
	credentialsDAO dao.CredentialsRepository,
	profileDAO dao.ProfileRepository,
	userDAO dao.UserRepository,
	mailer mailer.Mailer,
	generateValidationCode func() (string, string, error),
	generateTokenService GenerateTokenService,
	validateEmailTemplate string,
	validateEmailLink string,
) RegisterService {
	return &registerServiceImpl{
		credentialsDAO:         credentialsDAO,
		profileDAO:             profileDAO,
		userDAO:                userDAO,
		mailer:                 mailer,
		generateValidationCode: generateValidationCode,
		GenerateTokenService:   generateTokenService,
		validateEmailTemplate:  validateEmailTemplate,
		validateEmailLink:      validateEmailLink,
	}
}

type registerServiceImpl struct {
	credentialsDAO         dao.CredentialsRepository
	profileDAO             dao.ProfileRepository
	userDAO                dao.UserRepository
	mailer                 mailer.Mailer
	generateValidationCode func() (string, string, error)
	GenerateTokenService

	validateEmailTemplate string
	validateEmailLink     string
}

func (s *registerServiceImpl) Register(ctx context.Context, form models.RegisterForm, now time.Time) (*models.UserTokenStatus, func() error, error) {
	if err := validation.CheckMinMax(form.Email, MinEmailLength, MaxEmailLength); err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidEmail, err)
	}
	if err := validation.CheckMinMax(form.Password, MinPasswordLength, MaxPasswordLength); err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidPassword, err)
	}
	if err := validation.CheckMinMax(form.FirstName, 1, MaxNameLength); err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidFirstName, err)
	}
	if err := validation.CheckMinMax(form.LastName, 1, MaxNameLength); err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidLastName, err)
	}
	if err := validation.CheckMinMax(form.Slug, 1, MaxSlugLength); err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidSlug, err)
	}
	if err := validation.CheckMinMax(form.Username, -1, MaxUsernameLength); err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidUsername, err)
	}

	if err := validation.CheckRestricted(form.Sex, models.SexMale, models.SexFemale); err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidSex, err)
	}
	if err := validation.CheckRegexp(form.Slug, slugRegexp); err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidSlug, err)
	}
	if err := validation.CheckRegexp(form.FirstName, nameRegexp); err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidFirstName, err)
	}
	if err := validation.CheckRegexp(form.LastName, nameRegexp); err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidLastName, err)
	}
	if form.Username != "" {
		if err := validation.CheckRegexp(form.Username, usernameRegexp); err != nil {
			return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidUsername, err)
		}
	}

	daoEmail, err := dao.ParseEmail(form.Email)
	if err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidEmail, err)
	}

	age := getUserAge(form.Birthday, now)
	if err := validation.CheckMinMax(age, MinAge, MaxAge); err != nil {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidAge, err)
	}

	emailExists, err := s.credentialsDAO.EmailExists(ctx, daoEmail)
	if err != nil {
		return nil, nil, goerrors.Join(ErrEmailExists, err)
	}
	if emailExists {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidEmail, ErrTaken)
	}

	slugExists, err := s.profileDAO.SlugExists(ctx, form.Slug)
	if err != nil {
		return nil, nil, goerrors.Join(ErrSlugExists, err)
	}
	if slugExists {
		return nil, nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidSlug, ErrTaken)
	}

	// Generate the code to validate user email. The private (hashed) code goes in the database. The public code will
	// be sent to the user address, to ensure it is valid.
	publicValidationCode, privateValidationCode, err := s.generateValidationCode()
	if err != nil {
		return nil, nil, goerrors.Join(ErrGenerateValidationCode, err)
	}
	daoEmail.Validation = privateValidationCode

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, goerrors.Join(ErrHashPassword, err)
	}

	userID := uuid.New()
	token, err := s.GenerateToken(ctx, models.UserTokenPayload{ID: userID}, uuid.New(), now)
	if err != nil {
		return nil, nil, goerrors.Join(ErrGenerateToken, err)
	}

	user, err := s.userDAO.Create(ctx, &dao.UserModelCore{
		Credentials: dao.CredentialsModelCore{
			Email:    daoEmail,
			Password: dao.Password{Hashed: string(passwordHashed)},
		},
		Identity: dao.IdentityModelCore{
			FirstName: form.FirstName,
			LastName:  form.LastName,
			Birthday:  form.Birthday,
			Sex:       form.Sex,
		},
		Profile: dao.ProfileModelCore{
			Username: form.Username,
			Slug:     form.Slug,
		},
	}, userID, now)
	if err != nil {
		return nil, nil, goerrors.Join(ErrCreateUser, err)
	}

	// Perform heavy load, post registration tasks in the background, after response has been sent back to the user.
	deferred := func() error {
		to := mail.NewEmail(user.Identity.FirstName, form.Email)
		templateData := map[string]interface{}{
			"name":            user.Identity.FirstName,
			"validation_link": fmt.Sprintf("%s?id=%s&code=%s", s.validateEmailLink, user.ID, publicValidationCode),
		}

		if err := s.mailer.Send(to, s.validateEmailTemplate, templateData); err != nil {
			return goerrors.Join(ErrSendValidationEmail, err)
		}

		return nil
	}

	return token, deferred, nil
}
