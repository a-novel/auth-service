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

func TestSlugExistsHandler(t *testing.T) {
	data := []struct {
		name string

		slug string

		serviceResp bool
		serviceErr  error

		expectStatus int
	}{
		{
			name:         "Success",
			slug:         "slug",
			serviceResp:  true,
			expectStatus: http.StatusNoContent,
		},
		{
			name:         "Success/NotFound",
			slug:         "slug",
			serviceResp:  false,
			expectStatus: http.StatusNotFound,
		},
		{
			name:         "Error",
			slug:         "slug",
			serviceErr:   fooErr,
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewSlugExistsService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?slug="+d.slug, nil)

			service.On("SlugExists", c, d.slug).Return(d.serviceResp, d.serviceErr)

			handler := handlers.NewSlugExistsHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code)

			service.AssertExpectations(t)
		})
	}
}
