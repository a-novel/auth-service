package handlers_test

import (
	"encoding/json"
	"fmt"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/models"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	goframework "github.com/a-novel/go-framework"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSearchHandler(t *testing.T) {
	data := []struct {
		name string

		query  string
		limit  int
		offset int

		shouldCallService bool
		serviceResp       []*models.UserPreview
		serviceTotal      int
		serviceErr        error

		expect       interface{}
		expectStatus int
	}{
		{
			name:              "Success",
			query:             "elon-bezos",
			limit:             10,
			offset:            20,
			shouldCallService: true,
			serviceResp: []*models.UserPreview{
				{
					ID:        goframework.NumberUUID(1),
					FirstName: "name-1",
					LastName:  "surname-1",
					Username:  "username-1",
					Slug:      "slug-1",
					CreatedAt: baseTime,
				},
				{
					ID:        goframework.NumberUUID(2),
					FirstName: "name-2",
					LastName:  "surname-2",
					Username:  "username-2",
					Slug:      "slug-2",
					CreatedAt: baseTime,
				},
			},
			serviceTotal: 200,
			expect: map[string]interface{}{
				"total": float64(200),
				"res": []interface{}{
					map[string]interface{}{
						"id":        goframework.NumberUUID(1).String(),
						"firstName": "name-1",
						"lastName":  "surname-1",
						"username":  "username-1",
						"slug":      "slug-1",
						"createdAt": baseTime.Format(time.RFC3339),
					},
					map[string]interface{}{
						"id":        goframework.NumberUUID(2).String(),
						"firstName": "name-2",
						"lastName":  "surname-2",
						"username":  "username-2",
						"slug":      "slug-2",
						"createdAt": baseTime.Format(time.RFC3339),
					},
				},
			},
			expectStatus: http.StatusOK,
		},
		{
			name:              "Error/InvalidEntity",
			query:             "elon-bezos",
			limit:             10,
			offset:            20,
			shouldCallService: true,
			serviceErr:        goframework.ErrInvalidEntity,
			expectStatus:      http.StatusBadRequest,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewSearchService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/?query=%s&limit=%d&offset=%d", d.query, d.limit, d.offset), nil)

			if d.shouldCallService {
				service.On("Search", c, d.query, d.limit, d.offset).Return(d.serviceResp, d.serviceTotal, d.serviceErr)
			}

			handler := handlers.NewSearchHandler(service)
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
