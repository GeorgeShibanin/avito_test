package handlers

import (
	"encoding/json"
	"github.com/GeorgeShibanin/avito_test/internal/storage"
	"github.com/pkg/errors"
	"net/http"
)

func (h *HTTPHandler) HandlePutBalance(rw http.ResponseWriter, r *http.Request) {
	var userInfo PostUserBalance
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	newbalance, err := h.storage.PutBalance(r.Context(), storage.Id(userInfo.Id), storage.Balance(userInfo.Balance))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	response := ResponseUserInfo(newbalance)
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
