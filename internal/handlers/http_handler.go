package handlers

import (
	"github.com/GeorgeShibanin/avito_test/internal/storage"
)

type HTTPHandler struct {
	storage storage.Storage
}

func NewHTTPHandler(storage storage.Storage) *HTTPHandler {
	return &HTTPHandler{
		storage: storage,
	}
}

type PostUserBalance struct {
	Id      string `json:"id"`
	Balance int64  `json:"balance"`
}

type ResponseUserInfo struct {
	Id      string `json:"id"`
	Balance int64  `json:"balance"`
}

type PostOrder struct {
	Id        string `json:"id"`
	IdServise string `json:"id_servise"`
	IdOrder   string `json:"id_order"`
	Amount    int64  `json:"amount"`
}

type AcceptOrder struct {
	Id        string `json:"id"`
	IdServise string `json:"id_servise"`
	IdOrder   string `json:"id_order"`
	Amount    int64  `json:"amount"`
}

type CancelOrder struct {
	Id        string `json:"id"`
	IdServise string `json:"id_servise"`
	IdOrder   string `json:"id_order"`
	Amount    int64  `json:"amount"`
}

type ResponseAcceptedOrder struct {
	Id        string `json:"id"`
	IdServise string `json:"id_servise"`
	IdOrder   string `json:"id_order"`
	Amount    int64  `json:"reserved_balance"`
}

type ResponseOrder struct {
	Id        string `json:"id"`
	Balance   int64  `json:"balance"`
	IdServise string `json:"id_servise"`
	IdOrder   string `json:"id_order"`
	Amount    int64  `json:"reserved_balance"`
}
