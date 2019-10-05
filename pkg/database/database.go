package database

import (
	"context"
	"encoding/json"
	"fmt"
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
	ID         int             `json:"user_id,omitempty"`
	Email      string          `json:"email" validate:"required,email"`
	Password   string          `json:"password" validate:"required"`
	Name       string          `json:"name" validate:"required"`
	Phone      string          `json:"phone" validate:"required"`
	RegionID   int             `json:"region_id" validate:"required"`
	Meta       json.RawMessage `json:"meta,omitempty"`
	Version    int             `json:"version,omitempty"`
	Created    *time.Time      `json:"created,omitempty"`
	Updated    *time.Time      `json:"updated,omitempty"`
	LastLogin  *time.Time      `json:"last_login,omitempty"`
	LastAction *time.Time      `json:"last_action,omitempty"`
	IsBlocked  bool            `json:"is_blocked,omitempty"`
	IsDeleted  bool            `json:"is_deleted,omitempty"`
}

type UserAlreadyExistsError struct{}

func (UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user already exists")
}
