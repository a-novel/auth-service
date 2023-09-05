package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/go-framework/errors"
	"time"
)

type PreviewPrivateService interface {
	// Preview returns a subset of data for the current user.
	Preview(ctx context.Context, tokenRaw string, now time.Time) (*models.UserPreviewPrivate, error)
}

func NewPreviewPrivateService(
	credentialsDAO dao.CredentialsRepository,
	profileDAO dao.ProfileRepository,
	identityDAO dao.IdentityRepository,
	introspectTokenService IntrospectTokenService,
) PreviewPrivateService {
	return &previewPrivateServiceImpl{
		credentialsDAO:         credentialsDAO,
		profileDAO:             profileDAO,
		identityDAO:            identityDAO,
		IntrospectTokenService: introspectTokenService,
	}
}

type previewPrivateServiceImpl struct {
	credentialsDAO dao.CredentialsRepository
	profileDAO     dao.ProfileRepository
	identityDAO    dao.IdentityRepository

	IntrospectTokenService
}

func (s *previewPrivateServiceImpl) Preview(ctx context.Context, tokenRaw string, now time.Time) (*models.UserPreviewPrivate, error) {
	token, err := s.IntrospectToken(ctx, tokenRaw, now, false)
	if err != nil {
		return nil, goerrors.Join(ErrIntrospectToken, err)
	}
	if !token.OK {
		return nil, goerrors.Join(errors.ErrInvalidCredentials, ErrInvalidToken)
	}

	credentials, err := s.credentialsDAO.GetCredentials(ctx, token.Token.Payload.ID)
	if err != nil {
		return nil, goerrors.Join(ErrGetCredentials, err)
	}

	profile, err := s.profileDAO.GetProfile(ctx, token.Token.Payload.ID)
	if err != nil {
		return nil, goerrors.Join(ErrGetProfile, err)
	}

	identity, err := s.identityDAO.GetIdentity(ctx, token.Token.Payload.ID)
	if err != nil {
		return nil, goerrors.Join(ErrGetIdentity, err)
	}

	return &models.UserPreviewPrivate{
		ID:        token.Token.Payload.ID,
		Email:     credentials.Email.String(),
		NewEmail:  credentials.NewEmail.String(),
		Validated: credentials.Email.Validation == "",
		UserPreview: models.UserPreview{
			FirstName: identity.FirstName,
			LastName:  identity.LastName,
			Username:  profile.Username,
			Slug:      profile.Slug,
			CreatedAt: profile.CreatedAt,
		},
	}, nil
}
