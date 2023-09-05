package services

import (
	"context"
	"crypto/ed25519"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	"github.com/google/uuid"
)

type RotateSecretKeysService interface {
	RotateSecretKeys(ctx context.Context) error
}

func NewRotateSecretKeysService(
	secretKeysDAO dao.SecretKeysRepository,
	keyGen func() (ed25519.PrivateKey, error),
	maxBackups int,
) RotateSecretKeysService {
	return &rotateSecretKeysServiceImpl{
		secretKeysDAO: secretKeysDAO,
		keyGen:        keyGen,
		maxBackups:    maxBackups,
	}
}

type rotateSecretKeysServiceImpl struct {
	secretKeysDAO dao.SecretKeysRepository
	keyGen        func() (ed25519.PrivateKey, error)
	maxBackups    int
}

func (s *rotateSecretKeysServiceImpl) RotateSecretKeys(ctx context.Context) error {
	newKey, err := s.keyGen()
	if err != nil {
		return goerrors.Join(ErrGenerateSignatureKey, err)
	}

	if _, err := s.secretKeysDAO.Write(ctx, newKey, uuid.NewString()); err != nil {
		return goerrors.Join(ErrWriteSignatureKey, err)
	}

	keys, err := s.secretKeysDAO.List(ctx)
	if err != nil {
		return goerrors.Join(ErrListSignatureKeys, err)
	}

	if len(keys) > s.maxBackups {
		for _, extraKey := range keys[s.maxBackups:] {
			if err = s.secretKeysDAO.Delete(ctx, extraKey.Name); err != nil {
				return goerrors.Join(ErrDeleteSignatureKey, err)
			}
		}
	}

	return nil
}
