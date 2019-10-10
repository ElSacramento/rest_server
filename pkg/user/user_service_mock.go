package user

import "context"

type serviceMock struct {
	get    func(ctx context.Context, userID int64) (*serviceUser, error)
	insert func(ctx context.Context, usr *serviceUser) (int64, error)
}

func (s *serviceMock) Get(ctx context.Context, userID int64) (*serviceUser, error) {
	return s.get(ctx, userID)
}

func (s *serviceMock) Insert(ctx context.Context, usr *serviceUser) (int64, error) {
	return s.insert(ctx, usr)
}
