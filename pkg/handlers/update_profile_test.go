package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	goframework "github.com/a-novel/go-framework"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateProfileHandler(t *testing.T) {
	data := []struct {
		name string

		authorization string

		body interface{}

		shouldCallService     bool
		shouldCallServiceWith models.UpdateProfileForm

		serviceErr error

		expectStatus int
	}{
		{
			name:          "Success",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"username": "username",
				"slug":     "slug",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.UpdateProfileForm{
				Username: "username",
				Slug:     "slug",
			},
			expectStatus: http.StatusCreated,
		},
		{
			name:          "Error/ErrInvalidCredentials",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"username": "username",
				"slug":     "slug",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.UpdateProfileForm{
				Username: "username",
				Slug:     "slug",
			},
			serviceErr:   goframework.ErrInvalidCredentials,
			expectStatus: http.StatusForbidden,
		},
		{
			name:          "Error/ErrTaken",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"username": "username",
				"slug":     "slug",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.UpdateProfileForm{
				Username: "username",
				Slug:     "slug",
			},
			serviceErr:   services.ErrTaken,
			expectStatus: http.StatusConflict,
		},
		{
			name:          "Error/ErrInvalidEntity",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"username": "username",
				"slug":     "slug",
			},
			shouldCallService: true,
			shouldCallServiceWith: models.UpdateProfileForm{
				Username: "username",
				Slug:     "slug",
			},
			serviceErr:   goframework.ErrInvalidEntity,
			expectStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewUpdateProfileService(t)

			mrshBody, err := json.Marshal(d.body)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(mrshBody))
			c.Request.Header.Set("Authorization", d.authorization)

			if d.shouldCallService {
				service.
					On("UpdateProfile", c, d.authorization, mock.Anything, d.shouldCallServiceWith).
					Return(d.serviceErr)
			}

			handler := handlers.NewUpdateProfileHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code, c.Errors.String())

			service.AssertExpectations(t)
		})
	}
}
