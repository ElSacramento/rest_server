package user

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"rest_server/pkg/errors"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestHandler_Get(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	usr := serviceUser{
		ID:       1,
		Email:    "test@gmail.com",
		Name:     "test testov",
		Phone:    "791839405750",
		Password: "test",
		RegionID: 1,
		Meta:     []byte(`{"gender":"male"}`),
	}

	s := serviceMock{get: func(ctx context.Context, userID int64) (user *serviceUser, e error) {
		if userID == -1 {
			return nil, sql.ErrConnDone
		}
		if userID != 1 {
			return nil, errors.UserNotExistsError{}
		}
		return &usr, nil
	}}
	h := Handler{user: &s}

	t.Run("OK", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/user?user_id=1", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		h.Get(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		var resp responseUser
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("cant unmarshal: %v", err)
		}

		expected := responseUser{
			ID:       usr.ID,
			Email:    usr.Email,
			Name:     usr.Name,
			Phone:    usr.Phone,
			RegionID: usr.RegionID,
			Meta:     usr.Meta,
		}

		require.Equal(t, expected, resp)
	})
	t.Run("FAILED", func(t *testing.T) {
		t.Run("query without user_id key", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user?limit=1", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			h.Get(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "user_id is required")
		})
		t.Run("empty query", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			h.Get(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "user_id is required")
		})
		t.Run("not exists user", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user?user_id=2", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			h.Get(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "not exists")
		})
		t.Run("bad user_id", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user?user_id=test", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			h.Get(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "bad user_id")
		})
		t.Run("db error", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user?user_id=-1", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			h.Get(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		})
	})
}

func TestHandler_Add(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	s := serviceMock{insert: func(ctx context.Context, usr *serviceUser) (i int64, e error) {
		if usr.Email == "fail@gmail.com" {
			return 0, sql.ErrTxDone
		}
		if usr.Email == "exists@gmail.com" {
			return 0, errors.UserAlreadyExistsError{}
		}
		return 1, nil
	}}
	h := Handler{user: &s}

	type response struct {
		ID int64 `json:"user_id"`
	}

	t.Run("OK", func(t *testing.T) {
		b := []byte(`{"email": "test@gmail.com", "password": "pwd", "phone": "7919", "region_id": 1, "name": "test t"}`)
		r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
		require.NoError(t, err)

		w := httptest.NewRecorder()
		h.Add(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		var resp response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("cant unmarshal: %v", err)
		}

		require.Equal(t, response{ID: 1}, resp)
	})
	t.Run("FAIL", func(t *testing.T) {
		t.Run("empty email", func(t *testing.T) {
			b := []byte(`{"password": "pwd", "phone": "7919454685", "region_id": 1, "name": "test t"}`)
			r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			h.Add(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "some required params are empty")
		})
		t.Run("bad json", func(t *testing.T) {
			b := []byte(`{"password": "pwd", "phone": "7919454685", "region_id": 1, "name": "test t"`)
			r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			h.Add(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "Parse request body failed")
		})
		t.Run("bad json", func(t *testing.T) {
			b := []byte(`{"password": "pwd", "phone": {"a": 1}, "region_id": 1, "name": "test t"}`)
			r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			h.Add(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "Parse request body failed")
		})
		t.Run("db error", func(t *testing.T) {
			b := []byte(`{"email": "fail@gmail.com", "password": "pwd", "phone": "8919", "region_id": 1, "name": "test t"}`)
			r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			h.Add(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		})
		t.Run("exists email ", func(t *testing.T) {
			b := []byte(`{"email": "exists@gmail.com", "password": "pwd", "phone": "8919", "region_id": 1, "name": "test t"}`)
			r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			h.Add(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "user already exists")
		})
	})
}
