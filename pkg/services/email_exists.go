package services

import (
	"context"
	goerrors "errors"
	"github.com/a-novel/auth-service/pkg/dao"
	goframework "github.com/a-novel/go-framework"
)

type EmailExistsService interface {
	EmailExists(ctx context.Context, email string) (bool, error)
}

func NewEmailExistsService(
	credentialsDAO dao.CredentialsRepository,
) EmailExistsService {
	return &emailExistsServiceImpl{
		credentialsDAO: credentialsDAO,
	}
}

type emailExistsServiceImpl struct {
	credentialsDAO dao.CredentialsRepository
}

func (s *emailExistsServiceImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	daoEmail, err := dao.ParseEmail(email)
	if err != nil {
		return false, goerrors.Join(goframework.ErrInvalidEntity, err)
	}

	ok, err := s.credentialsDAO.EmailExists(ctx, daoEmail)
	if err != nil {
		return false, goerrors.Join(ErrEmailExists, err)
	}

	return ok, nil
}
