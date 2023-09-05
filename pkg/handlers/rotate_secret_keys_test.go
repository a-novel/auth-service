package handlers_test

import (
	"github.com/a-novel/auth-service/pkg/handlers"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRotateSecretKeysHandler(t *testing.T) {
	data := []struct {
		name string

		serviceErr error

		expectStatus int
	}{
		{
			name:         "Success",
			expectStatus: http.StatusCreated,
		},
		{
			name:         "Error",
			serviceErr:   fooErr,
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewRotateSecretKeysService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)

			service.On("RotateSecretKeys", c).Return(d.serviceErr)

			handler := handlers.NewRotateSecretKeysHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code)

			service.AssertExpectations(t)
		})
	}
}
