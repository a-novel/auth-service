package handlers_test

import (
	"github.com/a-novel/auth-service/pkg/handlers"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/go-framework/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResendEmailValidationHandler(t *testing.T) {
	data := []struct {
		name string

		authorization string

		serviceErr error

		expectStatus int
	}{
		{
			name:          "Success",
			authorization: "Bearer token",
			expectStatus:  http.StatusAccepted,
		},
		{
			name:          "Error/InvalidCredentials",
			authorization: "Bearer token",
			serviceErr:    errors.ErrInvalidCredentials,
			expectStatus:  http.StatusForbidden,
		},
		{
			name:          "Error/NotFound",
			authorization: "Bearer token",
			serviceErr:    errors.ErrNotFound,
			expectStatus:  http.StatusNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewResendEmailValidationService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Request.Header.Set("Authorization", d.authorization)

			service.On("ResendEmailValidation", c, d.authorization, mock.Anything).Return(nil, d.serviceErr)

			handler := handlers.NewResendEmailValidationHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code)

			service.AssertExpectations(t)
		})
	}
}
