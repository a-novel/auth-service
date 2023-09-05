package handlers_test

import (
	"encoding/json"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/auth-service/pkg/models"
	servicesmocks "github.com/a-novel/auth-service/pkg/services/mocks"
	"github.com/a-novel/go-framework/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIntrospectTokenHandler(t *testing.T) {
	data := []struct {
		name string

		authorization string

		serviceResp *models.UserTokenStatus
		serviceErr  error

		expect       interface{}
		expectStatus int
	}{
		{
			name:          "Success",
			authorization: "Bearer my-token",
			serviceResp: &models.UserTokenStatus{
				OK:        true,
				Expired:   false,
				NotIssued: false,
				Malformed: false,
				Token: &models.UserToken{
					Header: models.UserTokenHeader{
						IAT: baseTime,
						EXP: baseTime.Add(time.Hour),
						ID:  test.NumberUUID(10),
					},
					Payload: models.UserTokenPayload{ID: test.NumberUUID(1)},
				},
				TokenRaw: "Bearer my-token",
			},
			expect: map[string]interface{}{
				"ok":        true,
				"expired":   false,
				"notIssued": false,
				"malformed": false,
				"token": map[string]interface{}{
					"header": map[string]interface{}{
						"iat": baseTime.Format(time.RFC3339),
						"exp": baseTime.Add(time.Hour).Format(time.RFC3339),
						"id":  test.NumberUUID(10).String(),
					},
					"payload": map[string]interface{}{
						"id": test.NumberUUID(1).String(),
					},
				},
				"tokenRaw": "Bearer my-token",
			},
			expectStatus: http.StatusOK,
		},
		{
			name:          "Error/ServiceFailure",
			authorization: "Bearer my-token",
			serviceErr:    fooErr,
			expectStatus:  http.StatusInternalServerError,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			service := servicesmocks.NewIntrospectTokenService(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Request.Header.Set("Authorization", d.authorization)

			service.On("IntrospectToken", c, d.authorization, mock.Anything, true).Return(d.serviceResp, d.serviceErr)

			handler := handlers.NewIntrospectTokenHandler(service)
			handler.Handle(c)

			require.Equal(t, d.expectStatus, w.Code)
			if d.expect != nil {
				var body interface{}
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
				require.Equal(t, d.expect, body)
			}

			service.AssertExpectations(t)
		})
	}
}
