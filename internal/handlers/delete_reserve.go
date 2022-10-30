package handlers

import (
	"encoding/json"
	"github.com/GeorgeShibanin/avito_test/internal/storage"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func (h *HTTPHandler) HandleCancelReserve(rw http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(rw, "invalid query params", http.StatusBadRequest)
		return
	}

	id_servise := r.URL.Query().Get("id_servise")
	if id_servise == "" {
		http.Error(rw, "invalid query params", http.StatusBadRequest)
		return
	}

	id_order := r.URL.Query().Get("id_order")
	if id_order == "" {
		http.Error(rw, "invalid query params", http.StatusBadRequest)
		return
	}

	cost := r.URL.Query().Get("amount")
	if cost == "" {
		http.Error(rw, "invalid query params", http.StatusBadRequest)
		return
	}
	amount, err := strconv.Atoi(cost)
	if err != nil {
		http.Error(rw, "invalid query params", http.StatusBadRequest)
		return
	}
	reserve, err := h.storage.DeleteReserve(r.Context(), storage.Id(id), storage.IdServise(id_servise),
		storage.IdOrder(id_order), storage.Amout(amount))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rawResponse, err := json.Marshal(reserve + " this order deleted")
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
