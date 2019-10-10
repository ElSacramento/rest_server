package user

import (
	"context"
	"encoding/json"
	"rest_server/pkg/database"
)

type service interface {
	Get(ctx context.Context, userID int64) (*serviceUser, error)
	Insert(ctx context.Context, usr *serviceUser) (int64, error)
}

type Service struct {
	db database.DataStore
}

type serviceUser struct {
	ID       int64
	Email    string
	Password string
	Name     string
	Phone    string
	RegionID int64
	Meta     json.RawMessage
}

func (s *Service) Get(ctx context.Context, userID int64) (*serviceUser, error) {
	acc, err := s.db.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	usr := serviceUser{
		ID:       acc.ID,
		Email:    acc.Email,
		Name:     acc.Name,
		Phone:    acc.Phone,
		RegionID: acc.RegionID,
		Meta:     acc.Meta,
	}
	return &usr, nil
}

func (s *Service) securePassword(pwd string) string {
	// todo
	return pwd + "hash"
}

func (s *Service) Insert(ctx context.Context, usr *serviceUser) (int64, error) {
	acc := database.Account{
		Email:    usr.Email,
		Password: s.securePassword(usr.Password),
		Name:     usr.Name,
		Phone:    usr.Phone,
		RegionID: usr.RegionID,
		Meta:     usr.Meta,
	}
	userID, err := s.db.InsertUser(ctx, &acc)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
