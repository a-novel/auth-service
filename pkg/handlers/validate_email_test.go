package handlers_test

import (
	"fmt"
	"github.com/a-novel/auth-service/pkg/handlers"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/go-framework/errors"
	"github.com/a-novel/go-framework/test"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateEmailHandler(t *testing.T) {
	data := []struct {
		name string

		id   string
		code string

		shouldCallService bool
		serviceErr        error

		expectStatus int
	}{
		{
			name:              "Success",
			id:                test.NumberUUID(1).String(),
			code:              "validation-code",
			shouldCallService: true,
			expectStatus:      http.StatusNoContent,
		},
		{
			name:              "Error/ErrInvalidCredentials",
			id:                test.NumberUUID(1).String(),
			code:              "validation-code",
			shouldCallService: true,
			serviceErr:        errors.ErrInvalidCredentials,
			expectStatus:      http.StatusForbidden,
		},
		{
			name:              "Error/ErrInvalidEntity",
			id:                test.NumberUUID(1).String(),
			code:              "validation-code",
			shouldCallService: true,
			serviceErr:        errors.ErrInvalidEntity,
			expectStatus:      http.StatusForbidden,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewValidateEmailService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/?id=%s&code=%s", d.id, d.code), nil)

			if d.shouldCallService {
				service.On("ValidateEmail", c, uuid.MustParse(d.id), d.code, mock.Anything).Return(d.serviceErr)
			}

			handler := handlers.NewValidateEmailHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code, c.Errors.String())

			service.AssertExpectations(t)
		})
	}
}
