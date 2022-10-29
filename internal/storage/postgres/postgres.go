package postgres

import (
	"context"
	"fmt"
	"github.com/GeorgeShibanin/avito_test/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"log"
)

const (
	GetBalanceById = `SELECT id, balance FROM users WHERE id = $1`
	PatchBalance   = `UPDATE users SET balance = $1 WHERE id = $2
						RETURNING id, balance`
	InsertNewUser = `INSERT INTO users (id, balance) values ($1, $2)
						RETURNING id, balance`

	dsnTemplate = "postgres://%s:%s@%s:%v/%s"
)

type StoragePostgres struct {
	conn *pgx.Conn
}

func initConnection(conn *pgx.Conn) *StoragePostgres {
	return &StoragePostgres{conn: conn}
}

func Init(ctx context.Context, host, user, db, password string, port uint16) (*StoragePostgres, error) {
	conn, err := pgx.Connect(ctx, fmt.Sprintf(dsnTemplate, user, password, host, port, db))
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to postgres")
	}

	return initConnection(conn), nil
}

func (s *StoragePostgres) PutBalance(ctx context.Context, id storage.Id, balance storage.Balance) (storage.UserInfo, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return storage.UserInfo{}, errors.Wrap(err, "can't create tx")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	if err != nil {
		return storage.UserInfo{}, errors.Wrap(err, "can't create tx")
	}
	userInfo := &storage.UserInfo{}
	err = tx.QueryRow(ctx, GetBalanceById, id).Scan(&userInfo.Id, &userInfo.Balance)
	if err != nil {
		err = tx.QueryRow(ctx, InsertNewUser, id, balance).Scan(&userInfo.Id, &userInfo.Balance)
		if err != nil {
			return storage.UserInfo{}, nil
		}
		return *userInfo, nil
	}
	newBalance := int64(balance) + userInfo.Balance
	if newBalance < 0 {
		return storage.UserInfo{}, storage.ErrBelowZero
	}
	err = tx.QueryRow(ctx, PatchBalance, newBalance, id).Scan(&userInfo.Id, &userInfo.Balance)
	if err != nil {
		return storage.UserInfo{}, errors.Wrap(err, "Cant Update Balance of User")
	}
	return *userInfo, nil
}

func (s *StoragePostgres) GetBalance(ctx context.Context, id storage.Id) (storage.UserInfo, error) {
	userInfo := &storage.UserInfo{}
	log.Println(id)
	err := s.conn.QueryRow(ctx, GetBalanceById, id).Scan(&userInfo.Id, &userInfo.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.UserInfo{}, fmt.Errorf("User does not exist - %w", storage.StorageError)
		}
		return storage.UserInfo{}, fmt.Errorf("Error while Searching for new Item", err)
	}
	return *userInfo, nil
}
