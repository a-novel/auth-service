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
	"time"
)

func TestLoginHandler(t *testing.T) {
	data := []struct {
		name string

		body interface{}

		shouldCallService             bool
		shouldCallServiceWithEmail    string
		shouldCallServiceWithPassword string

		serviceResp *models.UserTokenStatus
		serviceErr  error

		expect       interface{}
		expectStatus int
	}{
		{
			name: "Success",
			body: map[string]interface{}{
				"email":    "email",
				"password": "password",
			},
			shouldCallService:             true,
			shouldCallServiceWithEmail:    "email",
			shouldCallServiceWithPassword: "password",
			serviceResp: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime,
						EXP: baseTime.Add(time.Hour),
						ID:  test.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
				TokenRaw: "Bearer my-token",
			},
			expect:       map[string]interface{}{"token": "Bearer my-token"},
			expectStatus: http.StatusOK,
		},
		{
			name: "Error/BadForm",
			body: map[string]interface{}{
				"email":    123456,
				"password": "password",
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Error/Forbidden",
			body: map[string]interface{}{
				"email":    "email",
				"password": "password",
			},
			shouldCallService:             true,
			shouldCallServiceWithEmail:    "email",
			shouldCallServiceWithPassword: "password",
			serviceErr:                    errors.ErrInvalidCredentials,
			expectStatus:                  http.StatusForbidden,
		},
		{
			name: "Error/NotFound",
			body: map[string]interface{}{
				"email":    "email",
				"password": "password",
			},
			shouldCallService:             true,
			shouldCallServiceWithEmail:    "email",
			shouldCallServiceWithPassword: "password",
			serviceErr:                    errors.ErrNotFound,
			expectStatus:                  http.StatusNotFound,
		},
		{
			name: "Error/InvalidEntity",
			body: map[string]interface{}{
				"email":    "email",
				"password": "password",
			},
			shouldCallService:             true,
			shouldCallServiceWithEmail:    "email",
			shouldCallServiceWithPassword: "password",
			serviceErr:                    errors.ErrInvalidEntity,
			expectStatus:                  http.StatusUnprocessableEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewLoginService(t)

			mrshBody, err := json.Marshal(d.body)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(mrshBody))

			if d.shouldCallService {
				service.
					On("Login", c, d.shouldCallServiceWithEmail, d.shouldCallServiceWithPassword, mock.Anything).
					Return(d.serviceResp, d.serviceErr)
			}

			handler := handlers.NewLoginHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code, c.Errors.String())
			if d.expect != nil {
				var body interface{}
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
				require.Equal(t, d.expect, body)
			}

			service.AssertExpectations(t)
		})
	}
}
