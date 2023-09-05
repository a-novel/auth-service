package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/models"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdatePasswordHandler(t *testing.T) {
	data := []struct {
		name string

		body interface{}

		shouldCallService     bool
		shouldCallServiceWith models.UpdatePasswordForm

		serviceErr error

		expectStatus int
	}{
		{
			name: "Success",
			body: map[string]interface{}{
				"id":          test.NumberUUID(1).String(),
				"code":        "validation-code",
				"oldPassword": "old-password",
				"newPassword": "new-password",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				Code:        "validation-code",
				OldPassword: "old-password",
				NewPassword: "new-password",
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Error/ErrInvalidCredentials",
			body: map[string]interface{}{
				"id":          test.NumberUUID(1).String(),
				"code":        "validation-code",
				"oldPassword": "old-password",
				"newPassword": "new-password",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				Code:        "validation-code",
				OldPassword: "old-password",
				NewPassword: "new-password",
			},
			serviceErr:   errors.ErrInvalidCredentials,
			expectStatus: http.StatusForbidden,
		},
		{
			name: "Error/ErrInvalidEntity",
			body: map[string]interface{}{
				"id":          test.NumberUUID(1).String(),
				"code":        "validation-code",
				"oldPassword": "old-password",
				"newPassword": "new-password",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				Code:        "validation-code",
				OldPassword: "old-password",
				NewPassword: "new-password",
			},
			serviceErr:   errors.ErrInvalidEntity,
			expectStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Error/ErrNotFound",
			body: map[string]interface{}{
				"id":          test.NumberUUID(1).String(),
				"code":        "validation-code",
				"oldPassword": "old-password",
				"newPassword": "new-password",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.UpdatePasswordForm{
				ID:          test.NumberUUID(1),
				Code:        "validation-code",
				OldPassword: "old-password",
				NewPassword: "new-password",
			},
			serviceErr:   errors.ErrNotFound,
			expectStatus: http.StatusForbidden,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewUpdatePasswordService(t)

			mrshBody, err := json.Marshal(d.body)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(mrshBody))

			if d.shouldCallService {
				service.
					On("UpdatePassword", c, d.shouldCallServiceWith, mock.Anything).
					Return(d.serviceErr)
			}

			handler := handlers.NewUpdatePasswordHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code, c.Errors.String())

			service.AssertExpectations(t)
		})
	}
}
