package user

import "rest_server/pkg/database"

type Module struct {
	s service
	H *Handler
}

func New(db database.DataStore) *Module {
	s := Service{db: db}
	return &Module{s: &s, H: &Handler{user: &s}}
}
