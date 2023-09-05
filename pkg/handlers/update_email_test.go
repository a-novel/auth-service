package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/services"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/go-framework/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateEmailHandler(t *testing.T) {
	data := []struct {
		name string

		authorization string

		body interface{}

		shouldCallService          bool
		shouldCallServiceWithEmail string

		serviceErr error

		expectStatus int
	}{
		{
			name:          "Success",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"newEmail": "new-email",
			},
			shouldCallService:          true,
			shouldCallServiceWithEmail: "new-email",
			expectStatus:               http.StatusAccepted,
		},
		{
			name:          "Error/ErrInvalidCredentials",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"newEmail": "new-email",
			},
			shouldCallService:          true,
			shouldCallServiceWithEmail: "new-email",
			serviceErr:                 errors.ErrInvalidCredentials,
			expectStatus:               http.StatusForbidden,
		},
		{
			name:          "Error/ErrTaken",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"newEmail": "new-email",
			},
			shouldCallService:          true,
			shouldCallServiceWithEmail: "new-email",
			serviceErr:                 services.ErrTaken,
			expectStatus:               http.StatusConflict,
		},
		{
			name:          "Error/ErrInvalidEntity",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"newEmail": "new-email",
			},
			shouldCallService:          true,
			shouldCallServiceWithEmail: "new-email",
			serviceErr:                 errors.ErrInvalidEntity,
			expectStatus:               http.StatusUnprocessableEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewUpdateEmailService(t)

			mrshBody, err := json.Marshal(d.body)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(mrshBody))
			c.Request.Header.Set("Authorization", d.authorization)

			if d.shouldCallService {
				service.
					On("UpdateEmail", c, d.authorization, d.shouldCallServiceWithEmail, mock.Anything).
					Return(nil, d.serviceErr)
			}

			handler := handlers.NewUpdateEmailHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code, c.Errors.String())

			service.AssertExpectations(t)
		})
	}
}
