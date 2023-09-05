package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type LoginService interface {
	// Login creates a new session for the given user, given the right credentials.
	Login(ctx context.Context, email string, password string, now time.Time) (*models.UserTokenStatus, error)
}

func NewLoginService(credentialsDAO dao.CredentialsRepository, generateTokenService GenerateTokenService) LoginService {
	return &loginServiceImpl{
		credentialsDAO:       credentialsDAO,
		GenerateTokenService: generateTokenService,
	}
}

type loginServiceImpl struct {
	credentialsDAO dao.CredentialsRepository
	GenerateTokenService
}

func (s *loginServiceImpl) Login(ctx context.Context, email string, password string, now time.Time) (*models.UserTokenStatus, error) {
	if err := validation.CheckMinMax(email, MinEmailLength, MaxEmailLength); err != nil {
		return nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidEmail, err)
	}

	if err := validation.CheckMinMax(password, MinPasswordLength, MaxPasswordLength); err != nil {
		return nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidPassword, err)
	}

	daoEmail, err := dao.ParseEmail(email)
	if err != nil {
		return nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidEmail, err)
	}

	user, err := s.credentialsDAO.GetCredentialsByEmail(ctx, daoEmail)
	if err != nil {
		return nil, goerrors.Join(ErrGetCredentialsByEmail, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.Hashed), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, goerrors.Join(errors.ErrInvalidCredentials, ErrWrongPassword)
		}

		return nil, goerrors.Join(ErrCheckPassword, err)
	}

	status, err := s.GenerateToken(ctx, models.UserTokenPayload{ID: user.ID}, uuid.New(), now)
	if err != nil {
		return nil, goerrors.Join(ErrGenerateToken, err)
	}

	return status, nil
}
