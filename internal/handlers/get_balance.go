package handlers

import (
	"encoding/json"
	"github.com/GeorgeShibanin/avito_test/internal/storage"
	"github.com/pkg/errors"
	"net/http"
)

func (h *HTTPHandler) HandleGetBalance(rw http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	if userId == "" {
		http.Error(rw, "invalid query params", http.StatusBadRequest)
		return
	}

	userBalance, err := h.storage.GetBalance(r.Context(), storage.Id(userId))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	response := ResponseUserInfo(userBalance)
	rawResponse, err := json.Marshal(response)
	if err != nil {
		err = errors.Wrap(err, "can't marshall response")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = rw.Write(rawResponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
