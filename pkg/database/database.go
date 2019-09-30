package database

import (
	"context"
	"encoding/json"
	"time"
)

type DataStore interface {
	GetUser(ctx context.Context, userID int) (*User, error)
	InsertUser(ctx context.Context, dbUser *User) (int64, error)
	CreateConnection() error
	CloseConnection() error
	//DeleteUser(context.Context, int) (interface{}, error)
}

type User struct {
	ID         int             `json:"user_id"`
	Email      string          `json:"email"`
	Password   string          `json:"password"`
	Name       string          `json:"name"`
	Phone      string          `json:"phone"`
	RegionID   int             `json:"region_id"`
	Meta       json.RawMessage `json:"meta"`
	Version    int             `json:"version"`
	Created    time.Time       `json:"created"`
	Updated    time.Time       `json:"updated"`
	LastLogin  time.Time       `json:"last_login"`
	LastAction time.Time       `json:"last_action"`
	IsBlocked  bool            `json:"is_blocked"`
	IsDeleted  bool            `json:"is_deleted"`
}
