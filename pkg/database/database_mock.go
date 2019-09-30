package database

import (
	"context"
)

type ServiceMock struct {
	GetUserMock    func(ctx context.Context, userID int) (*User, error)
	InsertUserMock func(ctx context.Context, dbUser *User) (int64, error)
}

func (s ServiceMock) GetUser(ctx context.Context, userID int) (*User, error) {
	return s.GetUserMock(ctx, userID)
}

func (s ServiceMock) InsertUser(ctx context.Context, dbUser *User) (int64, error) {
	return s.InsertUserMock(ctx, dbUser)
}

func (s ServiceMock) CreateConnection() error {
	panic("implements me")
}

func (s ServiceMock) CloseConnection() error {
	panic("implements me")
}
