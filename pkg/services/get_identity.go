package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	goframework "github.com/a-novel/go-framework"
	"time"
)

type GetIdentityService interface {
	Get(ctx context.Context, tokenRaw string, now time.Time) (*models.Identity, error)
}

func NewGetIdentityService(
	identityDAO dao.IdentityRepository,
	introspectTokenService IntrospectTokenService,
) GetIdentityService {
	return &getIdentityServiceImpl{
		identityDAO:            identityDAO,
		IntrospectTokenService: introspectTokenService,
	}
}

type getIdentityServiceImpl struct {
	identityDAO dao.IdentityRepository
	IntrospectTokenService
}

func (s *getIdentityServiceImpl) Get(ctx context.Context, tokenRaw string, now time.Time) (*models.Identity, error) {
	token, err := s.IntrospectToken(ctx, tokenRaw, now, false)
	if err != nil {
		return nil, goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return nil, goerrors.Join(goframework.ErrInvalidCredentials, ErrInvalidToken)
	}

	identity, err := s.identityDAO.GetIdentity(ctx, token.Token.Payload.ID)
	if err != nil {
		return nil, goerrors.Join(ErrGetIdentity, err)
	}

	return &models.Identity{
		FirstName: identity.FirstName,
		LastName:  identity.LastName,
		Sex:       identity.Sex,
		Birthday:  identity.Birthday,
	}, nil
}
