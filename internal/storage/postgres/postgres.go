package postgres

import (
	"context"
	"fmt"
	"githib.com/GeorgeShibanin/avito_test/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const (
	GetIdQuery    = `SELECT id, url FROM links WHERE id = $1`
	GetByUrlQuery = `SELECT id, url FROM links WHERE url = $1`
	InsertQuery   = `INSERT INTO links (id, url) values ($1, $2)`

	dsnTemplate = "postgres://%s:%s@%s:%v/%s"
)

type StoragePostgres struct {
	conn postgresInterface
}

type postgresInterface interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func initConnection(conn postgresInterface) *StoragePostgres {
	return &StoragePostgres{conn: conn}
}

func Init(ctx context.Context, host, user, db, password string, port uint16) (*StoragePostgres, error) {
	//подключение к базе через переменные окружения
	conn, err := pgx.Connect(ctx, fmt.Sprintf(dsnTemplate, user, password, host, port, db))
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to postgres")
	}

	return initConnection(conn), nil
}

func (s *StoragePostgres) PutBalance(ctx context.Context, id storage.Id, balance storage.Balance) (storage.UserInfo, error) {
	return storage.UserInfo{}, nil
}

func (s *StoragePostgres) GetBalance(ctx context.Context, id storage.Id) (storage.UserInfo, error) {
	return storage.UserInfo{}, nil
}
