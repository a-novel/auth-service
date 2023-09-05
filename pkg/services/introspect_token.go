package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/models"
	"time"
)

type IntrospectTokenService interface {
	// IntrospectToken parses, and verifies the provided token. If the autoRefresh flag is set to true, a new token
	// will automatically be issued when close enough to the expiration date.
	IntrospectToken(ctx context.Context, token string, now time.Time, autoRefresh bool) (*models.UserTokenStatus, error)
}

func NewIntrospectTokenService(
	generateTokenService GenerateTokenService,
	getTokenStatusService GetTokenStatusService,
	tokenRefreshThreshold time.Duration,
) IntrospectTokenService {
	return &introspectTokenServiceImpl{
		GenerateTokenService:  generateTokenService,
		GetTokenStatusService: getTokenStatusService,
		tokenRefreshThreshold: tokenRefreshThreshold,
	}
}

type introspectTokenServiceImpl struct {
	GenerateTokenService
	GetTokenStatusService

	tokenRefreshThreshold time.Duration
}

func (s *introspectTokenServiceImpl) IntrospectToken(ctx context.Context, token string, now time.Time, autoRefresh bool) (*models.UserTokenStatus, error) {
	status, err := s.GetTokenStatus(ctx, token, now)
	if err != nil {
		return nil, goerrors.Join(ErrGetTokenStatus, err)
	}

	if !status.OK {
		return status, nil
	}

	if autoRefresh && status.Token.Header.EXP.Sub(now) <= s.tokenRefreshThreshold {
		status, err = s.GenerateToken(ctx, status.Token.Payload, status.Token.Header.ID, now)
		if err != nil {
			return nil, goerrors.Join(ErrGenerateToken, err)
		}
	}

	return status, err
}
