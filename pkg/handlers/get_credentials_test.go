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
)

func TestGetCredentialsHandler(t *testing.T) {
	data := []struct {
		name string

		authorization string

		serviceResp *models.Credentials
		serviceErr  error

		expect       interface{}
		expectStatus int
	}{
		{
			name:          "Success",
			authorization: "Bearer token",
			serviceResp: &models.Credentials{
				Email:     "email",
				NewEmail:  "new-email",
				Validated: true,
			},
			expect: map[string]interface{}{
				"email":     "email",
				"newEmail":  "new-email",
				"validated": true,
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
			service := servicesmocks.NewGetCredentialsService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Request.Header.Set("Authorization", d.authorization)

			service.On("Get", c, d.authorization, mock.Anything).Return(d.serviceResp, d.serviceErr)

			handler := handlers.NewGetCredentialsHandler(service)
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
