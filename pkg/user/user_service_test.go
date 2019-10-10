package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"rest_server/pkg/database"
	"rest_server/pkg/errors"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestService_Get(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	ctx := context.Background()
	date := time.Now()

	sUser := serviceUser{
		ID:       1,
		Email:    "test@gmail.com",
		Name:     "test testov",
		Phone:    "791839405750",
		RegionID: 1,
		Meta:     []byte(`{"gender": "female"}`),
	}
	// check that service return not all db fields
	dbUser := database.Account{
		ID:         1,
		Email:      "test@gmail.com",
		Name:       "test testov",
		Phone:      "791839405750",
		RegionID:   1,
		Meta:       []byte(`{"gender": "female"}`),
		Created:    &date,
		Updated:    &date,
		Password:   "test",
		LastLogin:  &date,
		LastAction: &date,
	}
	s := Service{db: database.ServiceMock{
		GetUserMock: func(ctx context.Context, userID int64) (user *database.Account, e error) {
			if userID == -1 {
				return nil, sql.ErrConnDone
			}
			if userID != 1 {
				return nil, errors.UserNotExistsError{}
			}
			return &dbUser, nil
		}}}

	t.Run("OK", func(t *testing.T) {
		usr, err := s.Get(ctx, 1)
		require.NoError(t, err)
		require.Equal(t, sUser, *usr)
	})
	t.Run("FAIL", func(t *testing.T) {
		t.Run("not exists user", func(t *testing.T) {
			usr, err := s.Get(ctx, 2)
			require.Nil(t, usr)
			require.Error(t, err)
			require.Equal(t, errors.UserNotExistsError{}, err)
		})
		t.Run("db error", func(t *testing.T) {
			usr, err := s.Get(ctx, -1)
			require.Nil(t, usr)
			require.Error(t, err)
			require.Equal(t, sql.ErrConnDone, err)
		})
	})
}

func TestService_Insert(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	ctx := context.Background()

	s := Service{db: database.ServiceMock{InsertUserMock: func(ctx context.Context, dbUser *database.Account) (i int64, e error) {
		if dbUser.Email == "fail@gmail.com" {
			return 0, sql.ErrTxDone
		}
		if dbUser.Email == "exists@gmail.com" {
			return 0, errors.UserAlreadyExistsError{}
		}
		return 1, nil
	}}}

	t.Run("OK", func(t *testing.T) {
		sUser := serviceUser{
			Email:    "test@gmail.com",
			Name:     "test testov",
			Phone:    "791839405750",
			RegionID: 1,
			Meta:     []byte(`{"gender": "female"}`),
			Password: "test",
		}
		userID, err := s.Insert(ctx, &sUser)

		require.NoError(t, err)
		require.Equal(t, int64(1), userID)
	})
	t.Run("FAIL", func(t *testing.T) {
		t.Run("already exists user", func(t *testing.T) {
			sUser := serviceUser{
				Email:    "exists@gmail.com",
				Name:     "test testov",
				Phone:    "791839405750",
				RegionID: 1,
				Meta:     []byte(`{"gender": "female"}`),
				Password: "test",
			}
			userID, err := s.Insert(ctx, &sUser)

			require.Error(t, err)
			require.Equal(t, errors.UserAlreadyExistsError{}, err)
			require.Equal(t, int64(0), userID)
		})
		t.Run("db error", func(t *testing.T) {
			sUser := serviceUser{
				Email:    "fail@gmail.com",
				Name:     "test testov",
				Phone:    "791839405750",
				RegionID: 1,
				Meta:     []byte(`{"gender": "female"}`),
				Password: "test",
			}
			userID, err := s.Insert(ctx, &sUser)

			require.Error(t, err)
			require.Equal(t, sql.ErrTxDone, err)
			require.Equal(t, int64(0), userID)
		})
	})
}
