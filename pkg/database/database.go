package database

import (
	"context"
	"encoding/json"
	"time"
)

type DataStore interface {
	CreateConnection() error
	CloseConnection() error
	GetUser(ctx context.Context, userID int64) (*Account, error)
	InsertUser(ctx context.Context, dbUser *Account) (int64, error)
	//DeleteUser(context.Context, int) (interface{}, error)
}

type Account struct {
	ID         int64
	Email      string
	Password   string
	Name       string
	Phone      string
	RegionID   int64
	Meta       json.RawMessage
	Version    int
	Created    *time.Time
	Updated    *time.Time
	LastLogin  *time.Time
	LastAction *time.Time
	IsBlocked  bool
	IsDeleted  bool
}
