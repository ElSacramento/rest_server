package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"rest_server/pkg/database"
	"testing"
)

func TestService_Get(t *testing.T) {
	dbUser := database.User{
		ID:       1,
		Email:    "test@gmail.com",
		Name:     "test testov",
		Phone:    "791839405750",
		Password: "test",
		RegionID: 1,
		Version:  1,
		Meta:     []byte(`{"gender":"male"}`),
	}
	logrus.SetLevel(logrus.DebugLevel)
	s := Service{DB: database.ServiceMock{GetUserMock: func(ctx context.Context, userID int) (user *database.User, e error) {
		if userID == 0 {
			return nil, sql.ErrConnDone
		}
		if userID != 1 {
			return nil, nil
		}
		return &dbUser, nil
	}}}
	t.Run("OK", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/user?user_id=1", nil)
		require.Nil(t, err)

		w := httptest.NewRecorder()
		s.Get(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		var responseUser database.User
		if err := json.Unmarshal(w.Body.Bytes(), &responseUser); err != nil {
			t.Fatalf("cant unmarshal: %v", err)
		}

		require.Equal(t, dbUser, responseUser)
	})
	t.Run("FAILED", func(t *testing.T) {
		t.Run("query without user_id key", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user?limit=1", nil)
			require.Nil(t, err)

			w := httptest.NewRecorder()
			s.Get(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, string(w.Body.Bytes()), "Bad user id")
		})
		t.Run("empty query", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user", nil)
			require.Nil(t, err)

			w := httptest.NewRecorder()
			s.Get(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, string(w.Body.Bytes()), "Bad user id")
		})
		t.Run("not exists user", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user?user_id=2", nil)
			require.Nil(t, err)

			w := httptest.NewRecorder()
			s.Get(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, string(w.Body.Bytes()), "not exists")
		})
		t.Run("db error", func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/user?user_id=0", nil)
			require.Nil(t, err)

			w := httptest.NewRecorder()
			s.Get(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		})
	})
}
