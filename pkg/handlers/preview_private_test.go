package handlers_test

import (
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

func TestPreviewPrivateHandler(t *testing.T) {
	data := []struct {
		name string

		authorization string

		serviceResp *models.UserPreviewPrivate
		serviceErr  error

		expect       interface{}
		expectStatus int
	}{
		{
			name:          "Success",
			authorization: "Bearer token",
			serviceResp: &models.UserPreviewPrivate{
				ID:        goframework.NumberUUID(1),
				Email:     "email",
				NewEmail:  "new-email",
				Validated: true,
				UserPreview: models.UserPreview{
					FirstName: "name",
					LastName:  "last-name",
					Username:  "username",
					Slug:      "slug",
					CreatedAt: baseTime,
				},
			},
			expect: map[string]interface{}{
				"id":        goframework.NumberUUID(1).String(),
				"email":     "email",
				"newEmail":  "new-email",
				"validated": true,
				"firstName": "name",
				"lastName":  "last-name",
				"username":  "username",
				"slug":      "slug",
				"createdAt": baseTime.Format(time.RFC3339),
			},
			expectStatus: http.StatusOK,
		},
		{
			name:          "Error/Forbidden",
			authorization: "Bearer token",
			serviceErr:    goframework.ErrInvalidCredentials,
			expectStatus:  http.StatusForbidden,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewPreviewPrivateService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Request.Header.Set("Authorization", d.authorization)

			service.On("Preview", c, d.authorization, mock.Anything).Return(d.serviceResp, d.serviceErr)

			handler := handlers.NewPreviewPrivateHandler(service)
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
