package middleware

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"sync"
)

type database struct {
	user string
	password string
	baseURL string

	conn *pgx.Conn
	mut sync.Mutex
}

func (db *database) createConnection() error {
	connString := "postgres://" + db.user + ":" + db.password + "@"+ db.baseURL
	logrus.Infof("Open connection to database: %v", db.baseURL)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return xerrors.Errorf("Unable to connect to database: %v", err)
	}
	db.conn = conn
	return nil
}

func (db *database) closeConnection() error {
	logrus.Infof("Close connection to database: %v", db.baseURL)
	if err := db.conn.Close(context.Background()); err != nil {
		return xerrors.Errorf("Unable to close connect to database: %v", err)
	}
	return nil
}