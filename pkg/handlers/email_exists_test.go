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

func TestEmailExistsHandler(t *testing.T) {
	data := []struct {
		name string

		email string

		serviceResp bool
		serviceErr  error

		expectStatus int
	}{
		{
			name:         "Success",
			email:        "email",
			serviceResp:  true,
			expectStatus: http.StatusNoContent,
		},
		{
			name:         "Success/NotFound",
			email:        "email",
			serviceResp:  false,
			expectStatus: http.StatusNotFound,
		},
		{
			name:         "Error",
			email:        "email",
			serviceErr:   fooErr,
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewEmailExistsService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?email="+d.email, nil)

			service.On("EmailExists", c, d.email).Return(d.serviceResp, d.serviceErr)

			handler := handlers.NewEmailExistsHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code)

			service.AssertExpectations(t)
		})
	}
}
