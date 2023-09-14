package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/models"
	"github.com/a-novel/auth-service/pkg/services"
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

func TestRegisterHandler(t *testing.T) {
	data := []struct {
		name string

		body interface{}

		shouldCallService     bool
		shouldCallServiceWith models.RegisterForm
		serviceResp           *models.UserTokenStatus
		serviceErr            error

		expect       interface{}
		expectStatus int
	}{
		{
			name: "Success",
			body: map[string]interface{}{
				"email":     "email",
				"password":  "password",
				"slug":      "slug",
				"firstName": "name",
				"lastName":  "surname",
				"sex":       "male",
				"username":  "username",
				"birthday":  baseTime.Format(time.RFC3339),
			},
			shouldCallService: true,
			shouldCallServiceWith: models.RegisterForm{
				Email:     "email",
				Password:  "password",
				FirstName: "name",
				LastName:  "surname",
				Sex:       models.SexMale,
				Birthday:  baseTime,
				Slug:      "slug",
				Username:  "username",
			},
			serviceResp: &models.UserTokenStatus{
				OK: true,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime,
						EXP: baseTime.Add(time.Hour),
						ID:  goframework.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: goframework.NumberUUID(1)},
				},
				TokenRaw: "Bearer my-token",
			},
			expect:       map[string]interface{}{"token": "Bearer my-token"},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Error/BadForm",
			body: map[string]interface{}{
				"email":     "email",
				"password":  "password",
				"slug":      "slug",
				"firstName": "name",
				"lastName":  "surname",
				"sex":       "male",
				"username":  "username",
				"birthday":  "foobarqux",
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Error/Conflict",
			body: map[string]interface{}{
				"email":     "email",
				"password":  "password",
				"slug":      "slug",
				"firstName": "name",
				"lastName":  "surname",
				"sex":       "male",
				"username":  "username",
				"birthday":  baseTime.Format(time.RFC3339),
			},
			shouldCallService: true,
			shouldCallServiceWith: models.RegisterForm{
				Email:     "email",
				Password:  "password",
				FirstName: "name",
				LastName:  "surname",
				Sex:       models.SexMale,
				Birthday:  baseTime,
				Slug:      "slug",
				Username:  "username",
			},
			serviceErr:   services.ErrTaken,
			expectStatus: http.StatusConflict,
		},
		{
			name: "Error/InvalidEntity",
			body: map[string]interface{}{
				"email":     "email",
				"password":  "password",
				"slug":      "slug",
				"firstName": "name",
				"lastName":  "surname",
				"sex":       "male",
				"username":  "username",
				"birthday":  baseTime.Format(time.RFC3339),
			},
			shouldCallService: true,
			shouldCallServiceWith: models.RegisterForm{
				Email:     "email",
				Password:  "password",
				FirstName: "name",
				LastName:  "surname",
				Sex:       models.SexMale,
				Birthday:  baseTime,
				Slug:      "slug",
				Username:  "username",
			},
			serviceErr:   goframework.ErrInvalidEntity,
			expectStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewRegisterService(t)

			mrshBody, err := json.Marshal(d.body)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(mrshBody))

			if d.shouldCallService {
				service.
					On("Register", c, d.shouldCallServiceWith, mock.Anything).
					Return(d.serviceResp, nil, d.serviceErr)
			}

			handler := handlers.NewRegisterHandler(service)
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
