package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type PostgresDB struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string

	ConnectionPool *sql.DB
}

func (db *PostgresDB) CreateConnection() error {
	logrus.Infof("Create connection to %s", db.Name)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s", db.Host, db.Port, db.User, db.Password)
	pool, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return xerrors.Errorf("Failed to initialize database: %v", err)
	}

	if err := pool.Ping(); err != nil {
		return xerrors.Errorf("Database ping failed: %v", err)
	}

	db.ConnectionPool = pool
	return nil
}

func (db *PostgresDB) CloseConnection() error {
	logrus.Infof("Close connection to %s", db.Name)

	if err := db.ConnectionPool.Close(); err != nil {
		return xerrors.Errorf("Unable to close connect to database: %v", err)
	}
	return nil
}

func (db *PostgresDB) GetUser(ctx context.Context, userID int64) (*Account, error) {
	dbUser := Account{}
	err := db.ConnectionPool.QueryRowContext(
		ctx, `select user_id, email, password, name, phone, region_id, meta from account where user_id = $1;`, userID).
		Scan(dbUser.ID, dbUser.Email, dbUser.Password, dbUser.Name, dbUser.Phone, dbUser.RegionID, dbUser.Meta)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dbUser, nil
}

func (db *PostgresDB) InsertUser(ctx context.Context, dbUser *Account) (int64, error) {
	tx, err := db.ConnectionPool.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return 0, err
	}

	// todo pwd hash
	pwd := dbUser.Password
	createdAt := time.Now()
	var userID int64
	err = tx.QueryRowContext(ctx,
		"insert into account (email, password, name, phone, region_id, version, created, updated, last_login, last_action) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning user_id",
		dbUser.Email, pwd, dbUser.Name, dbUser.Phone, dbUser.RegionID, 1, createdAt, createdAt, time.Time{}, time.Time{}).Scan(&userID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				return 0, UserAlreadyExistsError{}
			}
			return 0, xerrors.Errorf("pq error: %s, code: %s", err.Code.Name(), err.Code)
		}
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}
