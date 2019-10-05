package user

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rest_server/pkg/database"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestService_Get(t *testing.T) {
	dbUser := database.Account{
		ID:       1,
		Email:    "test@gmail.com",
		Name:     "test testov",
		Phone:    "791839405750",
		Password: "test",
		RegionID: 1,
		Meta:     []byte(`{"gender":"male"}`),
	}
	logrus.SetLevel(logrus.DebugLevel)
	s := Service{DB: database.ServiceMock{GetUserMock: func(ctx context.Context, userID int64) (user *database.Account, e error) {
		if userID == -1 {
			return nil, sql.ErrConnDone
		}
		if userID != 1 {
			return nil, nil
		}
		return &dbUser, nil
	}}}
	t.Run("OK", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/user?user_id=1", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		s.Get(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		var responseUser database.Account
		if err := json.Unmarshal(w.Body.Bytes(), &responseUser); err != nil {
			t.Fatalf("cant unmarshal: %v", err)
		}

		require.Equal(t, dbUser, responseUser)
	})
	t.Run("FAILED", func(t *testing.T) {
		t.Run("query without user_id key", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user?limit=1", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			s.Get(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "user_id is required")
		})
		t.Run("empty query", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			s.Get(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "user_id is required")
		})
		t.Run("not exists user", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user?user_id=2", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			s.Get(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "not exists")
		})
		t.Run("db error", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user?user_id=-1", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			s.Get(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		})
	})
}

func TestService_Add(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	s := Service{DB: database.ServiceMock{InsertUserMock: func(ctx context.Context, dbUser *database.Account) (i int64, e error) {
		if dbUser.Email == "fail@gmail.com" {
			return 0, sql.ErrTxDone
		}
		if dbUser.Email == "exists@gmail.com" {
			return 0, database.UserAlreadyExistsError{}
		}
		return 1, nil
	}}}

	type response struct {
		ID int64 `json:"user_id"`
	}

	t.Run("OK", func(t *testing.T) {
		b := []byte(`{"email": "test@gmail.com", "password": "pwd", "phone": "7919", "region_id": 1, "name": "test t"}`)
		r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
		require.NoError(t, err)

		w := httptest.NewRecorder()
		s.Add(w, r)

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
			s.Add(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "some required params are empty")
		})
		t.Run("bad json", func(t *testing.T) {
			b := []byte(`{"password": "pwd", "phone": "7919454685", "region_id": 1, "name": "test t"`)
			r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			s.Add(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "Parse request body failed")
		})
		t.Run("bad json", func(t *testing.T) {
			b := []byte(`{"password": "pwd", "phone": {"a": 1}, "region_id": 1, "name": "test t"}`)
			r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			s.Add(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "Parse request body failed")
		})
		t.Run("db error", func(t *testing.T) {
			b := []byte(`{"email": "fail@gmail.com", "password": "pwd", "phone": "8919", "region_id": 1, "name": "test t"}`)
			r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			s.Add(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		})
		t.Run("exists email ", func(t *testing.T) {
			b := []byte(`{"email": "exists@gmail.com", "password": "pwd", "phone": "8919", "region_id": 1, "name": "test t"}`)
			r, err := http.NewRequest(http.MethodPost, "/user/add", bytes.NewReader(b))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			s.Add(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "user already exists")
		})
	})
}
