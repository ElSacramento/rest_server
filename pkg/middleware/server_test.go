package middleware

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"rest_server/pkg/user"
	"testing"
)

func TestServer_Run(t *testing.T) {
	srv, err := NewServer()
	if err != nil {
		t.Fatalf("failed to initialize server: %v", err)
	}
	srv.routes()
	t.Run("check route user", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/user", nil)
		if err != nil {
			t.Fatalf("wrong response: %v", err)
		}
		w := httptest.NewRecorder()
		srv.router.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
		var userStruct user.User
		if err := json.Unmarshal(w.Body.Bytes(), &userStruct); err != nil {
			t.Fatalf("cant unmarshal: %v", err)
		}
		assert.Equal(t, userStruct, user.User{ID: 1})
	})
	t.Run("check route user/add", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/user/add", nil)
		if err != nil {
			t.Fatalf("wrong response: %v", err)
		}
		w := httptest.NewRecorder()
		srv.router.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
		var userStruct user.User
		if err := json.Unmarshal(w.Body.Bytes(), &userStruct); err != nil {
			t.Fatalf("cant unmarshal: %v", err)
		}
		assert.Equal(t, userStruct, user.User{ID: 2})
	})
	t.Run("check route user/add with GET", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/user/add", nil)
		if err != nil {
			t.Fatalf("wrong response: %v", err)
		}
		w := httptest.NewRecorder()
		srv.router.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
	})
	t.Run("check route 404", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/user/not_exists", nil)
		if err != nil {
			t.Fatalf("wrong response: %v", err)
		}
		w := httptest.NewRecorder()
		srv.router.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusNotFound)
	})
}
