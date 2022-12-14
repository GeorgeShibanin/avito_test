package postgres

import (
	"context"
	"fmt"
	"github.com/GeorgeShibanin/avito_test/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"log"
	"time"
)

const (
	dsnTemplate = "postgres://%s:%s@%s:%v/%s"

	GetBalanceById = `SELECT id, balance FROM users WHERE id = $1`
	PatchBalance   = `UPDATE users SET balance = $1 WHERE id = $2
						RETURNING id, balance`
	InsertNewUser = `INSERT INTO users (id, balance) VALUES ($1, $2)
						RETURNING id, balance`

	GetReserveByOrderId = `SELECT id_user FROM orders WHERE id_order = $1`
	InsertNewReserve    = `INSERT INTO orders (id_user, id_service, id_order, amount, accepted) VALUES ($1, $2, $3, $4, $5)
							RETURNING id_user, id_service, id_order, amount, accepted`
	UpdateReserveAcceptance = `UPDATE orders SET accepted = $1 WHERE id_user = $2 AND id_service = $3 AND id_order = $4 AND amount = $5
							RETURNING id_user, id_service, id_order, amount, accepted`
	DeleteReserve = `DELETE FROM orders WHERE id_user = $1 AND id_service = $2 AND id_order = $3 AND amount = $4
						RETURNING id_order, accepted`

	InsertReport = `INSERT INTO report (id_user, id_service, id_order, amount, accepted_at)
						VALUES($1, $2, $3, $4, $5)
						RETURNING id_user`
	GetReport = `SELECT id_service, SUM(amount) FROM report WHERE accepted_at >= $1 AND accepted_at <= $2 GROUP BY id_service;`
)

type CommandTag string

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
	userInfo := &storage.UserInfo{}
	err = tx.QueryRow(ctx, GetBalanceById, id).Scan(&userInfo.Id, &userInfo.Balance)
	if err != nil {
		if balance < 0 {
			return storage.UserInfo{}, storage.ErrWrongBalance
		}
		err = tx.QueryRow(ctx, InsertNewUser, id, balance).Scan(&userInfo.Id, &userInfo.Balance)
		if err != nil {
			return storage.UserInfo{}, err
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
	err := s.conn.QueryRow(ctx, GetBalanceById, id).Scan(&userInfo.Id, &userInfo.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.UserInfo{}, fmt.Errorf("User does not exist - %w", storage.StorageError)
		}
		return storage.UserInfo{}, fmt.Errorf("Error while Searching for User", err)
	}
	return *userInfo, nil
}

func (s *StoragePostgres) PutReserve(ctx context.Context, id storage.Id, service storage.IdService, order storage.IdOrder, amout storage.Amout) (storage.Order, int64, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return storage.Order{}, 0, errors.Wrap(err, "can't create tx")
	}
	defer func() {
		if err != nil {
			log.Println(err)
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	userId := storage.Id("")
	err = tx.QueryRow(ctx, GetReserveByOrderId, order).Scan(&userId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Println(err)
		return storage.Order{}, 0, errors.Wrap(err, "can't reserve by order id")
	}
	if userId != "" {
		return storage.Order{}, 0, storage.ReserveAlreadyExist
	}

	//check if user is exists
	userInfo, err := s.GetBalance(ctx, id)
	if err != nil {
		return storage.Order{}, 0, err
	}
	//Check user balance to know if it's enough to make reserve
	if userInfo.Balance < int64(amout) || int64(amout) < 0 {
		return storage.Order{}, 0, storage.ErrBelowZero
	}

	//Insert New Order
	newReserve := &storage.Order{}
	err = tx.QueryRow(ctx, InsertNewReserve, id, service, order, amout, false).
		Scan(&newReserve.IdUser, &newReserve.IdService, &newReserve.IdOrder, &newReserve.Amount, &newReserve.Accepted)
	if err != nil {
		return storage.Order{}, 0, err
	}

	newBalance := userInfo.Balance - int64(amout)
	err = tx.QueryRow(ctx, PatchBalance, newBalance, id).Scan(&userInfo.Id, &userInfo.Balance)
	if err != nil {
		return storage.Order{}, 0, errors.Wrap(err, "Cant Update Balance of User")
	}
	return *newReserve, newBalance, nil
}

func (s *StoragePostgres) PatchReserve(ctx context.Context, id storage.Id, service storage.IdService, order storage.IdOrder, amount storage.Amout) (storage.Order, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return storage.Order{}, errors.Wrap(err, "can't create tx")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	reserve := &storage.Order{}
	//Update Reserve Status
	err = tx.QueryRow(ctx, UpdateReserveAcceptance, true, id, service, order, amount).
		Scan(&reserve.IdUser, &reserve.IdService, &reserve.IdOrder, &reserve.Amount, &reserve.Accepted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.Order{}, fmt.Errorf("Reserve not exist with such params - %w", storage.StorageError)
		}
		return storage.Order{}, fmt.Errorf("error while accept reserve - %w", err)
	}

	//Add accepted reserve to report
	userId := storage.Id("")
	err = tx.QueryRow(ctx, InsertReport, id, service, order, amount, time.Now().UTC().Format(time.RFC3339)).Scan(&userId)
	if err != nil {
		return storage.Order{}, storage.ErrAlreadyAccepted
	}
	return *reserve, nil
}

func (s *StoragePostgres) GetReport(ctx context.Context, date1 storage.Date, date2 storage.Date) ([]storage.Deals, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return []storage.Deals{}, errors.Wrap(err, "can't create tx")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	deals := &storage.Deals{}
	allDeals := []storage.Deals{}
	rows, err := tx.Query(ctx, GetReport, date1, date2)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []storage.Deals{}, fmt.Errorf("Such reports does not exist - %w", storage.StorageError)
		}
		return []storage.Deals{}, fmt.Errorf("Error while Searching for dates", err)
	}
	for rows.Next() {
		err = rows.Scan(&deals.IdServise, &deals.TotalSumm)
		allDeals = append(allDeals, *deals)
	}
	return allDeals, nil
}

func (s *StoragePostgres) DeleteReserve(ctx context.Context, id storage.Id, service storage.IdService, order storage.IdOrder, amount storage.Amout) (string, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return " storage.Item{}", errors.Wrap(err, "can't create tx")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	idOrder := storage.IdOrder("")
	acceptance := false
	err = tx.QueryRow(ctx, DeleteReserve, id, service, order, amount).Scan(&idOrder, &acceptance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("item not exist - %w", storage.StorageError)
		}
		return "", fmt.Errorf("cant delete item - %w", err)
	}
	if acceptance == true {
		return "", storage.ErrAlreadyAccepted
	}

	s.PutBalance(ctx, id, storage.Balance(amount))
	return string(idOrder), nil
}
