package storage

import (
	"context"
	"errors"
	"fmt"
)

var (
	StorageError    = errors.New("storage")
	ErrCollision    = fmt.Errorf("%w.collision", StorageError)
	ErrAlreadyExist = errors.New("key already exist")
	ErrNotFound     = fmt.Errorf("%w.not_found", StorageError)
)

type Id string
type Balance int64

type UserInfo struct {
	Id      string `json:"id"`
	Balance string `json:"balance"`
}

type Storage interface {
	PutBalance(ctx context.Context, id Id, balance Balance) (UserInfo, error)
	GetBalance(ctx context.Context, id Id) (UserInfo, error)
}
