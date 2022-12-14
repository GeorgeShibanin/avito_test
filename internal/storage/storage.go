package storage

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	StorageError        = errors.New("storage")
	ErrCollision        = fmt.Errorf("%w.collision", StorageError)
	ErrBelowZero        = errors.New("user does not have enough balance")
	ReserveAlreadyExist = errors.New("Reserve with such order id already exist")
	ErrNotFound         = fmt.Errorf("%w.not_found", StorageError)
	ErrWrongBalance     = errors.New("wrong balance")
	ErrAlreadyAccepted  = errors.New("error already accepted")
)

type Id string
type Balance int64
type IdService string
type IdOrder string
type Amout int64
type Date time.Time

type UserInfo struct {
	Id      string `json:"id"`
	Balance int64  `json:"balance"`
}

type Order struct {
	IdUser    string `json:"id_user"`
	IdService string `json:"id_service"`
	IdOrder   string `json:"id_order"`
	Amount    int64  `json:"amount"`
	Accepted  bool   `json:"accepted"`
}

type Deals struct {
	IdServise string `json:"idServise"`
	TotalSumm int    `json:"totalSumm"`
}

type Storage interface {
	PutBalance(ctx context.Context, id Id, balance Balance) (UserInfo, error)
	GetBalance(ctx context.Context, id Id) (UserInfo, error)

	PutReserve(ctx context.Context, id Id, servise IdService, order IdOrder, amout Amout) (Order, int64, error)
	PatchReserve(ctx context.Context, id Id, servise IdService, order IdOrder, amout Amout) (Order, error)
	DeleteReserve(ctx context.Context, id Id, servise IdService, order IdOrder, amout Amout) (string, error)

	GetReport(ctx context.Context, date1 Date, date2 Date) ([]Deals, error)
}
