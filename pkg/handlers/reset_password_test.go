package handlers_test

import (
	"github.com/a-novel/auth-service/pkg/handlers"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/bunovel"
	goframework "github.com/a-novel/go-framework"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResetPasswordHandler(t *testing.T) {
	data := []struct {
		name string

		email string

		serviceErr error

		expectStatus int
	}{
		{
			name:         "Success",
			email:        "email",
			expectStatus: http.StatusAccepted,
		},
		{
			name:         "Error/ErrInvalidEntity",
			email:        "email",
			serviceErr:   goframework.ErrInvalidEntity,
			expectStatus: http.StatusBadRequest,
		},
		{
			name:         "Error/NotFound",
			email:        "email",
			serviceErr:   bunovel.ErrNotFound,
			expectStatus: http.StatusNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewResetPasswordService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?email="+d.email, nil)

			service.On("ResetPassword", c, d.email, mock.Anything).Return(nil, d.serviceErr)

			handler := handlers.NewResetPasswordHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code)

			service.AssertExpectations(t)
		})
	}
}
