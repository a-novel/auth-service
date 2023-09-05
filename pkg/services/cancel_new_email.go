package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/go-framework/errors"
	"time"
)

type CancelNewEmailService interface {
	CancelNewEmail(ctx context.Context, tokenRaw string, now time.Time) error
}

func NewCancelNewEmailService(
	credentialsDAO dao.CredentialsRepository,
	introspectTokenService IntrospectTokenService,
) CancelNewEmailService {
	return &cancelNewEmailServiceImpl{
		credentialsDAO:         credentialsDAO,
		IntrospectTokenService: introspectTokenService,
	}
}

type cancelNewEmailServiceImpl struct {
	credentialsDAO dao.CredentialsRepository
	IntrospectTokenService
}

func (s *cancelNewEmailServiceImpl) CancelNewEmail(ctx context.Context, tokenRaw string, now time.Time) error {
	token, err := s.IntrospectToken(ctx, tokenRaw, now, false)
	if err != nil {
		return goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return goerrors.Join(errors.ErrInvalidCredentials, ErrInvalidToken)
	}

	_, err = s.credentialsDAO.CancelNewEmail(ctx, token.Token.Payload.ID, now)
	if err != nil {
		return goerrors.Join(ErrCancelNewEmail, err)
	}

	return nil
}
