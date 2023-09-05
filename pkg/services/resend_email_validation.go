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

type ResendEmailValidationService interface {
	ResendEmailValidation(ctx context.Context, tokenRaw string, now time.Time) (func() error, error)
}

func NewResendEmailValidationService(
	credentialsDAO dao.CredentialsRepository,
	identityDAO dao.IdentityRepository,
	mailer mailer.Mailer,
	generateValidationLink func() (string, string, error),
	introspectTokenService IntrospectTokenService,
	validateEmailLink string,
	validateEmailTemplate string,
) ResendEmailValidationService {
	return &resendEmailValidationServiceImpl{
		credentialsDAO:         credentialsDAO,
		identityDAO:            identityDAO,
		mailer:                 mailer,
		generateValidationLink: generateValidationLink,
		IntrospectTokenService: introspectTokenService,
		validateEmailLink:      validateEmailLink,
		validateEmailTemplate:  validateEmailTemplate,
	}
}

type resendEmailValidationServiceImpl struct {
	credentialsDAO         dao.CredentialsRepository
	identityDAO            dao.IdentityRepository
	mailer                 mailer.Mailer
	generateValidationLink func() (string, string, error)
	IntrospectTokenService

	validateEmailLink     string
	validateEmailTemplate string
}

func (s *resendEmailValidationServiceImpl) ResendEmailValidation(ctx context.Context, tokenRaw string, now time.Time) (func() error, error) {
	token, err := s.IntrospectToken(ctx, tokenRaw, now, false)
	if err != nil {
		return nil, goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return nil, goerrors.Join(errors.ErrInvalidCredentials, ErrInvalidToken)
	}

	publicValidationCode, privateValidationCode, err := s.generateValidationLink()
	if err != nil {
		return nil, goerrors.Join(ErrGenerateValidationCode, err)
	}

	credentials, err := s.credentialsDAO.UpdateEmailValidation(ctx, privateValidationCode, token.Token.Payload.ID, now)
	if err != nil {
		return nil, goerrors.Join(ErrUpdateEmailValidation, err)
	}

	identity, err := s.identityDAO.GetIdentity(ctx, token.Token.Payload.ID)
	if err != nil {
		return nil, goerrors.Join(ErrGetIdentity, err)
	}

	deferred := func() error {
		to := mail.NewEmail(identity.FirstName, credentials.Email.String())
		templateData := map[string]interface{}{
			"name":            identity.FirstName,
			"validation_link": fmt.Sprintf("%s?id=%s&code=%s", s.validateEmailLink, token.Token.Payload.ID, publicValidationCode),
		}

		if err := s.mailer.Send(to, s.validateEmailTemplate, templateData); err != nil {
			return goerrors.Join(ErrSendValidationEmail, err)
		}

		return nil
	}

	return deferred, nil
}
