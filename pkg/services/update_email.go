package services

import (
	"context"
	goerrors "errors"
	"fmt"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/mailer"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"time"
)

type UpdateEmailService interface {
	UpdateEmail(ctx context.Context, tokenRaw, newEmail string, now time.Time) (func() error, error)
}

func NewUpdateEmailService(
	credentialsDAO dao.CredentialsRepository,
	identityDAO dao.IdentityRepository,
	mailer mailer.Mailer,
	generateValidationLink func() (string, string, error),
	introspectTokenService IntrospectTokenService,
	validateNewEmailLink string,
	validateNewEmailTemplate string,
) UpdateEmailService {
	return &updateEmailServiceImpl{
		credentialsDAO:           credentialsDAO,
		identityDAO:              identityDAO,
		mailer:                   mailer,
		generateValidationLink:   generateValidationLink,
		IntrospectTokenService:   introspectTokenService,
		validateNewEmailLink:     validateNewEmailLink,
		validateNewEmailTemplate: validateNewEmailTemplate,
	}
}

type updateEmailServiceImpl struct {
	credentialsDAO         dao.CredentialsRepository
	identityDAO            dao.IdentityRepository
	mailer                 mailer.Mailer
	generateValidationLink func() (string, string, error)
	IntrospectTokenService

	validateNewEmailLink     string
	validateNewEmailTemplate string
}

func (s *updateEmailServiceImpl) UpdateEmail(ctx context.Context, tokenRaw, newEmail string, now time.Time) (func() error, error) {
	token, err := s.IntrospectToken(ctx, tokenRaw, now, false)
	if err != nil {
		return nil, goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return nil, goerrors.Join(errors.ErrInvalidCredentials, ErrInvalidToken)
	}

	newDAOEmail, err := dao.ParseEmail(newEmail)
	if err != nil {
		return nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidEmail, err)
	}

	emailExists, err := s.credentialsDAO.EmailExists(ctx, newDAOEmail)
	if err != nil {
		return nil, goerrors.Join(ErrEmailExists, err)
	}
	if emailExists {
		return nil, goerrors.Join(errors.ErrInvalidEntity, ErrInvalidEmail, ErrTaken)
	}

	publicValidationCode, privateValidationCode, err := s.generateValidationLink()
	if err != nil {
		return nil, goerrors.Join(ErrGenerateValidationCode, err)
	}

	if _, err := s.credentialsDAO.UpdateEmail(ctx, newDAOEmail, privateValidationCode, token.Token.Payload.ID, now); err != nil {
		return nil, goerrors.Join(ErrUpdateEmail, err)
	}

	identity, err := s.identityDAO.GetIdentity(ctx, token.Token.Payload.ID)
	if err != nil {
		return nil, goerrors.Join(ErrGetIdentity, err)
	}

	deferred := func() error {
		to := mail.NewEmail(identity.FirstName, newEmail)
		templateData := map[string]interface{}{
			"name":            identity.FirstName,
			"validation_link": fmt.Sprintf("%s?id=%s&code=%s", s.validateNewEmailLink, token.Token.Payload.ID, publicValidationCode),
		}

		if err := s.mailer.Send(to, s.validateNewEmailTemplate, templateData); err != nil {
			return goerrors.Join(ErrSendValidationEmail, err)
		}

		return nil
	}

	return deferred, nil
}
