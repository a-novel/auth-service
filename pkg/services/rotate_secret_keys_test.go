package services_test

import (
	"context"
	"crypto/ed25519"
	"github.com/a-novel/auth-service/pkg/dao"
	daomocks "github.com/a-novel/auth-service/pkg/dao/mocks"
	"github.com/a-novel/auth-service/pkg/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRotateSecretKeys(t *testing.T) {
	type deleteCall struct {
		key string
		err error
	}

	data := []struct {
		name string

		maxBackups int

		keyGen    ed25519.PrivateKey
		keyGenErr error

		shouldCallWrite bool
		writeErr        error

		shouldCallList bool
		list           []*dao.SecretKeyModel
		listErr        error

		shouldCallDelete bool
		deleteCalls      []deleteCall

		expectErr error
	}{
		{
			name:            "Success/NoInitialKeys",
			maxBackups:      3,
			keyGen:          MockedSecretKeys[0],
			shouldCallWrite: true,
			shouldCallList:  true,
			list: []*dao.SecretKeyModel{
				{
					Name: "key-0",
					Key:  MockedSecretKeys[0],
				},
			},
		},
		{
			name:            "Success",
			maxBackups:      3,
			keyGen:          MockedSecretKeys[0],
			shouldCallWrite: true,
			shouldCallList:  true,
			list: []*dao.SecretKeyModel{
				{
					Name: "key-0",
					Key:  MockedSecretKeys[0],
				},
				{
					Name: "key-1",
					Key:  MockedSecretKeys[1],
				},
				{
					Name: "key-2",
					Key:  MockedSecretKeys[2],
				},
			},
		},
		{
			name:            "Success/TooMuchKeys",
			maxBackups:      2,
			keyGen:          MockedSecretKeys[0],
			shouldCallWrite: true,
			shouldCallList:  true,
			list: []*dao.SecretKeyModel{
				{
					Name: "key-0",
					Key:  MockedSecretKeys[0],
				},
				{
					Name: "key-1",
					Key:  MockedSecretKeys[1],
				},
				{
					Name: "key-2",
					Key:  MockedSecretKeys[2],
				},
				{
					Name: "key-3",
					Key:  MockedSecretKeys[3],
				},
			},
			shouldCallDelete: true,
			deleteCalls: []deleteCall{
				{key: "key-2"},
				{key: "key-3"},
			},
		},
		{
			name:            "Error/DeleteKeyFailure",
			maxBackups:      2,
			keyGen:          MockedSecretKeys[0],
			shouldCallWrite: true,
			shouldCallList:  true,
			list: []*dao.SecretKeyModel{
				{
					Name: "key-0",
					Key:  MockedSecretKeys[0],
				},
				{
					Name: "key-1",
					Key:  MockedSecretKeys[1],
				},
				{
					Name: "key-2",
					Key:  MockedSecretKeys[2],
				},
				{
					Name: "key-3",
					Key:  MockedSecretKeys[3],
				},
			},
			shouldCallDelete: true,
			deleteCalls: []deleteCall{
				{key: "key-2"},
				{key: "key-3", err: fooErr},
			},
			expectErr: fooErr,
		},
		{
			name:            "Error/ListFailure",
			maxBackups:      3,
			keyGen:          MockedSecretKeys[0],
			shouldCallWrite: true,
			shouldCallList:  true,
			listErr:         fooErr,
			expectErr:       fooErr,
		},
		{
			name:            "Error/WriteFailure",
			maxBackups:      3,
			keyGen:          MockedSecretKeys[0],
			shouldCallWrite: true,
			writeErr:        fooErr,
			expectErr:       fooErr,
		},
		{
			name:       "Error/KeyGenFailure",
			maxBackups: 3,
			keyGenErr:  fooErr,
			expectErr:  fooErr,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			secretKeysDAO := daomocks.NewSecretKeysRepository(t)

			keyGen := func() (ed25519.PrivateKey, error) {
				return d.keyGen, d.keyGenErr
			}

			if d.shouldCallWrite {
				secretKeysDAO.On("Write", context.Background(), d.keyGen, mock.Anything).Return(nil, d.writeErr)
			}

			if d.shouldCallList {
				secretKeysDAO.On("List", context.Background()).Return(d.list, d.listErr)
			}

			if d.shouldCallDelete {
				for _, call := range d.deleteCalls {
					secretKeysDAO.On("Delete", context.Background(), call.key).Return(call.err)
				}
			}

			service := services.NewRotateSecretKeysService(secretKeysDAO, keyGen, d.maxBackups)
			err := service.RotateSecretKeys(context.Background())

			require.ErrorIs(t, err, d.expectErr)

			secretKeysDAO.AssertExpectations(t)
		})
	}
}
