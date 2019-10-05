package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRequestQuery(t *testing.T) {
	t.Run("GET without query", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/user", nil)
		if err != nil {
			t.Fatalf("failed initialize request: %v", err)
		}
		rawJSON, err := GetRequestQuery(r)
		require.JSONEq(t, `{}`, string(rawJSON))
	})
	t.Run("GET with query", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/user?user_id=1", nil)
		if err != nil {
			t.Fatalf("failed initialize request: %v", err)
		}
		rawJSON, err := GetRequestQuery(r)
		require.JSONEq(t, `{"user_id":"1"}`, string(rawJSON))
	})
	t.Run("GET with empty value query", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/user?user_id=", nil)
		if err != nil {
			t.Fatalf("failed initialize request: %v", err)
		}
		rawJSON, err := GetRequestQuery(r)
		require.JSONEq(t, `{"user_id":""}`, string(rawJSON))
	})
	t.Run("GET with multiple query", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/user?user_id=1&limit=1", nil)
		if err != nil {
			t.Fatalf("failed initialize request: %v", err)
		}
		rawJSON, err := GetRequestQuery(r)
		require.JSONEq(t, `{"limit":"1","user_id":"1"}`, string(rawJSON))
	})
	t.Run("GET with multiple values query", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/user?filter=name&filter=region", nil)
		if err != nil {
			t.Fatalf("failed initialize request: %v", err)
		}
		rawJSON, err := GetRequestQuery(r)
		require.JSONEq(t, `{"filter":["name","region"]}`, string(rawJSON))
	})
}

func TestGetRequestBody(t *testing.T) {
	t.Run("POST without body", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodPost, "/user/add", nil)
		if err != nil {
			t.Fatalf("failed initialize request: %v", err)
		}
		rawJSON, err := GetRequestBody(r)
		require.JSONEq(t, `{}`, string(rawJSON))
	})
	t.Run("POST with body", func(t *testing.T) {
		b := []byte(`{"name": "test", "email": "example@gmail.com"}`)
		r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
		if err != nil {
			t.Fatalf("failed initialize request: %v", err)
		}
		rawJSON, err := GetRequestBody(r)
		require.JSONEq(t, `{"name": "test", "email": "example@gmail.com"}`, string(rawJSON))
	})
}

func TestSendResponse(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		w := httptest.NewRecorder()
		b := []byte(`{"name": "test", "phone": "79195432880"}`)
		if err := SendResponse(w, b); err != nil {
			t.Fatalf("failed to write json: %v", err)
		}
		var response map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("cant unmarshal: %v", err)
		}
		expected := map[string]string{
			"name": "test", "phone": "79195432880",
		}
		require.Equal(t, expected, response)
	})
	t.Run("FAIL", func(t *testing.T) {
		w := httptest.NewRecorder()
		b := []byte(`{"name": "test", "phone": "79195432880"`)
		if err := SendResponse(w, b); err != nil {
			t.Error("should be error")
		}
	})
}
