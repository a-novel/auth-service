package handlers_test

import (
	"encoding/json"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/models"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/go-framework/test"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListHandler(t *testing.T) {
	data := []struct {
		name string

		query string

		shouldCallService     bool
		shouldCallServiceWith []uuid.UUID
		serviceResp           []*models.UserPreview
		serviceErr            error

		expect       interface{}
		expectStatus int
	}{
		{
			name:                  "Success",
			query:                 "01010101-0101-0101-0101-010101010101,02020202-0202-0202-0202-020202020202",
			shouldCallService:     true,
			shouldCallServiceWith: []uuid.UUID{test.NumberUUID(1), test.NumberUUID(2)},
			serviceResp: []*models.UserPreview{
				{
					Username:  "username 1",
					Slug:      "slug 1",
					CreatedAt: baseTime,
				},
				{
					FirstName: "first name 2",
					LastName:  "last name 2",
					Slug:      "slug 2",
					CreatedAt: baseTime,
				},
			},
			expect: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"firstName": "",
						"lastName":  "",
						"username":  "username 1",
						"slug":      "slug 1",
						"createdAt": baseTime.Format(time.RFC3339),
					},
					map[string]interface{}{
						"firstName": "first name 2",
						"lastName":  "last name 2",
						"username":  "",
						"slug":      "slug 2",
						"createdAt": baseTime.Format(time.RFC3339),
					},
				},
			},
			expectStatus: http.StatusOK,
		},
		{
			name:              "Success/NoID",
			shouldCallService: true,
			expect: map[string]interface{}{
				"users": nil,
			},
			expectStatus: http.StatusOK,
		},
		{
			name:                  "Success/IgnoreInvalidIDs",
			query:                 "fake-id,02020202-0202-0202-0202-020202020202",
			shouldCallService:     true,
			shouldCallServiceWith: []uuid.UUID{test.NumberUUID(2)},
			serviceResp: []*models.UserPreview{
				{
					FirstName: "first name 2",
					LastName:  "last name 2",
					Slug:      "slug 2",
					CreatedAt: baseTime,
				},
			},
			expect: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"firstName": "first name 2",
						"lastName":  "last name 2",
						"username":  "",
						"slug":      "slug 2",
						"createdAt": baseTime.Format(time.RFC3339),
					},
				},
			},
			expectStatus: http.StatusOK,
		},
		{
			name:                  "Error/ServiceFailure",
			query:                 "01010101-0101-0101-0101-010101010101,02020202-0202-0202-0202-020202020202",
			shouldCallService:     true,
			shouldCallServiceWith: []uuid.UUID{test.NumberUUID(1), test.NumberUUID(2)},
			serviceErr:            fooErr,
			expectStatus:          http.StatusInternalServerError,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewListService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?ids="+d.query, nil)

			if d.shouldCallService {
				service.On("List", c, d.shouldCallServiceWith).Return(d.serviceResp, d.serviceErr)
			}

			handler := handlers.NewListHandler(service)
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
