package services

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	goerrors "errors"
	"fmt"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/a-novel/auth-service/pkg/models"
	goframework "github.com/a-novel/go-framework"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"strings"
	"time"
)

type GenerateTokenService interface {
	GenerateToken(ctx context.Context, data models.UserTokenPayload, id uuid.UUID, now time.Time) (*models.UserTokenStatus, error)
}

func NewGenerateTokenService(secretKeysDAO dao.SecretKeysRepository, tokenTTL time.Duration) GenerateTokenService {
	return &generateTokenServiceImpl{
		secretKeysDAO: secretKeysDAO,
		tokenTTL:      tokenTTL,
	}
}

type generateTokenServiceImpl struct {
	secretKeysDAO dao.SecretKeysRepository
	tokenTTL      time.Duration
}

func (s *generateTokenServiceImpl) GenerateToken(ctx context.Context, data models.UserTokenPayload, id uuid.UUID, now time.Time) (*models.UserTokenStatus, error) {
	signatureKeys, err := s.secretKeysDAO.List(ctx)
	if err != nil {
		return nil, goerrors.Join(ErrListSignatureKeys, err)
	}
	if len(signatureKeys) == 0 {
		return nil, ErrMissingSignatureKeys
	}

	source := models.UserToken{
		Header:  models.UserTokenHeader{IAT: now, EXP: now.Add(s.tokenTTL), ID: id},
		Payload: data,
	}

	mrshHeader, err := json.Marshal(source.Header)
	if err != nil {
		return nil, goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidTokenHeader, err)
	}
	header := base64.RawURLEncoding.EncodeToString(mrshHeader)

	mrshPayload, err := json.Marshal(source.Payload)
	if err != nil {
		return nil, goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidTokenPayload, err)
	}
	payload := base64.RawURLEncoding.EncodeToString(mrshPayload)

	unsigned := fmt.Sprintf("%s.%s", header, payload)
	signature := base64.RawURLEncoding.EncodeToString(ed25519.Sign(signatureKeys[0].Key, []byte(unsigned)))

	return &models.UserTokenStatus{
		OK:       true,
		Token:    &source,
		TokenRaw: fmt.Sprintf("%s.%s", unsigned, signature),
	}, nil
}

type GetTokenStatusService interface {
	GetTokenStatus(ctx context.Context, token string, now time.Time) (*models.UserTokenStatus, error)
}

func NewGetTokenStatusService(secretKeysDAO dao.SecretKeysRepository) GetTokenStatusService {
	return &getTokenStatusServiceImpl{
		secretKeysDAO: secretKeysDAO,
	}
}

type getTokenStatusServiceImpl struct {
	secretKeysDAO dao.SecretKeysRepository
}

func (s *getTokenStatusServiceImpl) splitToken(token string) (string, string, string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", "", "", goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidToken)
	}

	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	return header, payload, signature, nil
}

func (s *getTokenStatusServiceImpl) decodeToken(header, payload, signature string) ([]byte, []byte, []byte, error) {
	decodedSignature, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return nil, nil, nil, goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidTokenSignature, err)
	}

	decodedHeader, err := base64.RawURLEncoding.DecodeString(header)
	if err != nil {
		return nil, nil, nil, goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidTokenHeader, err)
	}

	decodedPayload, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, nil, nil, goerrors.Join(goframework.ErrInvalidEntity, ErrInvalidTokenPayload, err)
	}

	return decodedHeader, decodedPayload, decodedSignature, nil
}

func (s *getTokenStatusServiceImpl) validateToken(ctx context.Context, header, payload string, decodedSignature []byte) error {
	keys, err := s.secretKeysDAO.List(ctx)
	if err != nil {
		return goerrors.Join(ErrListSignatureKeys, err)
	}

	_, ok := lo.Find(keys, func(signatureKey *dao.SecretKeyModel) bool {
		return ed25519.Verify(
			// We know for sure the public type is correct, because we read it from the private key.
			signatureKey.Key.Public().(ed25519.PublicKey),
			[]byte(fmt.Sprintf("%s.%s", header, payload)),
			decodedSignature,
		)
	})

	if !ok {
		return goerrors.Join(goframework.ErrInvalidCredentials, ErrNoSignatureMatch)
	}

	return nil
}

func (s *getTokenStatusServiceImpl) GetTokenStatus(ctx context.Context, token string, now time.Time) (*models.UserTokenStatus, error) {
	status := &models.UserTokenStatus{TokenRaw: token}

	if token == "" {
		return status, nil
	}

	status.TokenRaw = token

	header, payload, signature, err := s.splitToken(token)
	if err != nil {
		if goerrors.Is(err, goframework.ErrInvalidEntity) {
			status.Malformed = true
			return status, nil
		}

		return nil, err
	}

	decodedHeader, decodedPayload, decodedSignature, err := s.decodeToken(header, payload, signature)
	if err != nil {
		if goerrors.Is(err, goframework.ErrInvalidEntity) {
			status.Malformed = true
			return status, nil
		}

		return nil, err
	}

	if err := s.validateToken(ctx, header, payload, decodedSignature); err != nil {
		if goerrors.Is(err, goframework.ErrInvalidCredentials) {
			status.Expired = true
		} else {
			return nil, goerrors.Join(ErrValidateToken, err)
		}
	}

	parsedToken := new(models.UserToken)

	if err := json.Unmarshal(decodedHeader, &parsedToken.Header); err != nil {
		status.Malformed = true
		return status, nil
	}
	if err := json.Unmarshal(decodedPayload, &parsedToken.Payload); err != nil {
		status.Malformed = true
		return status, nil
	}

	status.Token = parsedToken

	if !status.Expired {
		if parsedToken.Header.ID == uuid.Nil {
			status.Malformed = true
			return status, nil
		}
		if parsedToken.Header.IAT.After(now) {
			status.NotIssued = true
			return status, nil
		}
		if parsedToken.Header.EXP.Before(now) {
			status.Expired = true
			return status, nil
		}

		status.OK = true
	}

	return status, nil
}
