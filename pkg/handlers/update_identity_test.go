package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/models"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	goframework "github.com/a-novel/go-framework"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUpdateIdentityHandler(t *testing.T) {
	data := []struct {
		name string

		authorization string

		body interface{}

		shouldCallService     bool
		shouldCallServiceWith models.UpdateIdentityForm

		serviceErr error

		expectStatus int
	}{
		{
			name:          "Success",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"firstName": "first-name",
				"lastName":  "last-name",
				"sex":       "male",
				"birthday":  baseTime.Format(time.RFC3339),
			},
			shouldCallService: true,
			shouldCallServiceWith: models.UpdateIdentityForm{
				FirstName: "first-name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime,
			},
			expectStatus: http.StatusCreated,
		},
		{
			name:          "Error/ErrInvalidCredentials",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"firstName": "first-name",
				"lastName":  "last-name",
				"sex":       "male",
				"birthday":  baseTime.Format(time.RFC3339),
			},
			shouldCallService: true,
			shouldCallServiceWith: models.UpdateIdentityForm{
				FirstName: "first-name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime,
			},
			serviceErr:   goframework.ErrInvalidCredentials,
			expectStatus: http.StatusForbidden,
		},
		{
			name:          "Error/ErrInvalidEntity",
			authorization: "Bearer my-token",
			body: map[string]interface{}{
				"firstName": "first-name",
				"lastName":  "last-name",
				"sex":       "male",
				"birthday":  baseTime.Format(time.RFC3339),
			},
			shouldCallService: true,
			shouldCallServiceWith: models.UpdateIdentityForm{
				FirstName: "first-name",
				LastName:  "last-name",
				Sex:       models.SexMale,
				Birthday:  baseTime,
			},
			serviceErr:   goframework.ErrInvalidEntity,
			expectStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewUpdateIdentityService(t)

			mrshBody, err := json.Marshal(d.body)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(mrshBody))
			c.Request.Header.Set("Authorization", d.authorization)

			if d.shouldCallService {
				service.
					On("UpdateIdentity", c, d.authorization, mock.Anything, d.shouldCallServiceWith).
					Return(d.serviceErr)
			}

			handler := handlers.NewUpdateIdentityHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code, c.Errors.String())

			service.AssertExpectations(t)
		})
	}
}
