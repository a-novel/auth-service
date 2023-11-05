package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	goframework "github.com/a-novel/go-framework"
	"time"
)

type GetCredentialsService interface {
	Get(ctx context.Context, tokenRaw string, now time.Time) (*models.Credentials, error)
}

func NewGetCredentialsService(
	credentialsDAO dao.CredentialsRepository,
	introspectTokenService IntrospectTokenService,
) GetCredentialsService {
	return &getCredentialsServiceImpl{
		credentialsDAO:         credentialsDAO,
		IntrospectTokenService: introspectTokenService,
	}
}

type getCredentialsServiceImpl struct {
	credentialsDAO dao.CredentialsRepository
	IntrospectTokenService
}

func (s *getCredentialsServiceImpl) Get(ctx context.Context, tokenRaw string, now time.Time) (*models.Credentials, error) {
	token, err := s.IntrospectToken(ctx, tokenRaw, now, false)
	if err != nil {
		return nil, goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return nil, goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidToken)
	}

	credentials, err := s.credentialsDAO.GetCredentials(ctx, token.Token.Payload.ID)
	if err != nil {
		return nil, goerrors.Join(ErrGetCredentials, err)
	}

	return &models.Credentials{
		Email:     credentials.Email.String(),
		NewEmail:  credentials.NewEmail.String(),
		Validated: credentials.Email.Validation == "",
	}, nil
}
