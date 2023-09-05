package handlers_test

import (
	"encoding/json"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/models"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/go-framework/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPreviewHandler(t *testing.T) {
	data := []struct {
		name string

		slug string

		serviceResp *models.UserPreview
		serviceErr  error

		expect       interface{}
		expectStatus int
	}{
		{
			name: "Success",
			slug: "slug",
			serviceResp: &models.UserPreview{
				FirstName: "name",
				LastName:  "last-name",
				Username:  "username",
				Slug:      "slug",
				CreatedAt: baseTime,
			},
			expect: map[string]interface{}{
				"firstName": "name",
				"lastName":  "last-name",
				"username":  "username",
				"slug":      "slug",
				"createdAt": baseTime.Format(time.RFC3339),
			},
			expectStatus: http.StatusOK,
		},
		{
			name:         "Error/NotFound",
			slug:         "slug",
			serviceErr:   errors.ErrNotFound,
			expectStatus: http.StatusNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewPreviewService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?slug="+d.slug, nil)

			service.On("Preview", c, d.slug).Return(d.serviceResp, d.serviceErr)

			handler := handlers.NewPreviewHandler(service)
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
