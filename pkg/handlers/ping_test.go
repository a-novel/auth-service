package handlers_test

import (
	"encoding/json"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	data := []struct {
		name string

		expect       interface{}
		expectStatus int
	}{
		{
			name:         "Success",
			expect:       map[string]interface{}{"message": "pong"},
			expectStatus: http.StatusOK,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler := handlers.NewPingHandler()
			handler.Handle(c)

			var body interface{}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

			require.Equal(t, d.expectStatus, w.Code)
			require.Equal(t, d.expect, body)
		})
	}
}
