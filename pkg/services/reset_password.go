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

type ResetPasswordService interface {
	ResetPassword(ctx context.Context, email string, now time.Time) (func() error, error)
}

func NewResetPasswordService(
	credentialsDAO dao.CredentialsRepository,
	identityDAO dao.IdentityRepository,
	mailer sendgridproxy.Mailer,
	generateValidationLink func() (string, string, error),
	passwordResetLink string,
	passwordResetTemplate string,
) ResetPasswordService {
	return &resetPasswordServiceImpl{
		credentialsDAO:         credentialsDAO,
		identityDAO:            identityDAO,
		mailer:                 mailer,
		generateValidationLink: generateValidationLink,
		passwordResetLink:      passwordResetLink,
		passwordResetTemplate:  passwordResetTemplate,
	}
}

type resetPasswordServiceImpl struct {
	credentialsDAO         dao.CredentialsRepository
	identityDAO            dao.IdentityRepository
	mailer                 sendgridproxy.Mailer
	generateValidationLink func() (string, string, error)

	passwordResetLink     string
	passwordResetTemplate string
}

func (s *resetPasswordServiceImpl) ResetPassword(ctx context.Context, email string, now time.Time) (func() error, error) {
	daoEmail, err := dao.ParseEmail(email)
	if err != nil {
		return nil, goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidEmail, err)
	}

	publicValidationCode, privateValidationCode, err := s.generateValidationLink()
	if err != nil {
		return nil, goerrors.Join(ErrGenerateValidationCode, err)
	}

	credentials, err := s.credentialsDAO.ResetPassword(ctx, privateValidationCode, daoEmail, now)
	if err != nil {
		return nil, goerrors.Join(ErrResetPassword, err)
	}

	identity, err := s.identityDAO.GetIdentity(ctx, credentials.ID)
	if err != nil {
		return nil, goerrors.Join(ErrGetIdentity, err)
	}

	deferred := func() error {
		to := mail.NewEmail(identity.FirstName, credentials.CredentialsModelCore.Email.String())
		templateData := map[string]interface{}{
			"name":            identity.FirstName,
			"validation_link": fmt.Sprintf("%s?id=%s&code=%s", s.passwordResetLink, credentials.ID, publicValidationCode),
		}

		if err := s.mailer.Send(ctx, to, s.passwordResetTemplate, templateData); err != nil {
			return goerrors.Join(ErrSendValidationEmail, err)
		}

		return nil
	}

	return deferred, nil
}
