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
