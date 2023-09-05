package handlers_test

import (
	"encoding/json"
	"github.com/a-novel/auth-service/pkg/handlers"
	"github.com/a-novel/go-framework/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	db, sqlDB := test.GetPostgres(t, []fs.FS{})
	defer db.Close()
	defer sqlDB.Close()

	data := []struct {
		name string

		expect       interface{}
		expectStatus int
	}{
		{
			name:         "Success",
			expectStatus: http.StatusOK,
			expect: map[string]interface{}{
				"database": map[string]interface{}{
					"available": true,
				},
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler := handlers.NewHealthCheckHandler(db)
			handler.Handle(c)

			var body interface{}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

			require.Equal(t, d.expectStatus, w.Code)
			require.Equal(t, d.expect, body)
		})
	}
}
