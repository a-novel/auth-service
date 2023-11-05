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

func TestGetIdentityHandler(t *testing.T) {
	data := []struct {
		name string

		authorization string

		serviceResp *models.Identity
		serviceErr  error

		expect       interface{}
		expectStatus int
	}{
		{
			name:          "Success",
			authorization: "Bearer token",
			serviceResp: &models.Identity{
				FirstName: "name",
				LastName:  "last-name",
				Birthday:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Sex:       models.SexMale,
			},
			expect: map[string]interface{}{
				"firstName": "name",
				"lastName":  "last-name",
				"birthday":  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
				"sex":       "male",
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
			service := servicesmocks.NewGetIdentityService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Request.Header.Set("Authorization", d.authorization)

			service.On("Get", c, d.authorization, mock.Anything).Return(d.serviceResp, d.serviceErr)

			handler := handlers.NewGetIdentityHandler(service)
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
