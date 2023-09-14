package services

import (
	"context"
	goerrors "errors"
	"fmt"
	"github.com/a-novel/auth-service/pkg/dao"
	goframework "github.com/a-novel/go-framework"
	sendgridproxy "github.com/a-novel/sendgrid-proxy"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"time"
)

type ResendNewEmailValidationService interface {
	ResendNewEmailValidation(ctx context.Context, tokenRaw string, now time.Time) (func() error, error)
}

func NewResendNewEmailValidationService(
	credentialsDAO dao.CredentialsRepository,
	identityDAO dao.IdentityRepository,
	mailer sendgridproxy.Mailer,
	generateValidationLink func() (string, string, error),
	introspectTokenService IntrospectTokenService,
	validateNewEmailLink string,
	validateNewEmailTemplate string,
) ResendNewEmailValidationService {
	return &resendNewEmailValidationServiceImpl{
		credentialsDAO:           credentialsDAO,
		identityDAO:              identityDAO,
		mailer:                   mailer,
		generateValidationLink:   generateValidationLink,
		IntrospectTokenService:   introspectTokenService,
		validateNewEmailLink:     validateNewEmailLink,
		validateNewEmailTemplate: validateNewEmailTemplate,
	}
}

type resendNewEmailValidationServiceImpl struct {
	credentialsDAO         dao.CredentialsRepository
	identityDAO            dao.IdentityRepository
	mailer                 sendgridproxy.Mailer
	generateValidationLink func() (string, string, error)
	IntrospectTokenService

	validateNewEmailLink     string
	validateNewEmailTemplate string
}

func (s *resendNewEmailValidationServiceImpl) ResendNewEmailValidation(ctx context.Context, tokenRaw string, now time.Time) (func() error, error) {
	token, err := s.IntrospectToken(ctx, tokenRaw, now, false)
	if err != nil {
		return nil, goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return nil, goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidToken)
	}

	publicValidationCode, privateValidationCode, err := s.generateValidationLink()
	if err != nil {
		return nil, goerrors.Join(ErrGenerateValidationCode, err)
	}

	credentials, err := s.credentialsDAO.UpdateNewEmailValidation(ctx, privateValidationCode, token.Token.Payload.ID, now)
	if err != nil {
		return nil, goerrors.Join(ErrUpdateNewEmailValidation, err)
	}

	identity, err := s.identityDAO.GetIdentity(ctx, token.Token.Payload.ID)
	if err != nil {
		return nil, goerrors.Join(ErrGetIdentity, err)
	}

	deferred := func() error {
		to := mail.NewEmail(identity.FirstName, credentials.NewEmail.String())
		templateData := map[string]interface{}{
			"name":            identity.FirstName,
			"validation_link": fmt.Sprintf("%s?id=%s&code=%s", s.validateNewEmailLink, token.Token.Payload.ID, publicValidationCode),
		}

		if err := s.mailer.Send(ctx, to, s.validateNewEmailTemplate, templateData); err != nil {
			return goerrors.Join(ErrSendValidationEmail, err)
		}

		return nil
	}

	return deferred, nil
}
