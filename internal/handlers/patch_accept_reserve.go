package handlers

import (
	"encoding/json"
	"github.com/GeorgeShibanin/avito_test/internal/storage"
	"github.com/pkg/errors"
	"net/http"
)

func (h *HTTPHandler) HandlePatchAcceptReserve(rw http.ResponseWriter, r *http.Request) {
	var reserveInfo AcceptOrder
	err := json.NewDecoder(r.Body).Decode(&reserveInfo)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	reserve, err := h.storage.PatchReserve(r.Context(), storage.Id(reserveInfo.Id), storage.IdService(reserveInfo.IdServise),
		storage.IdOrder(reserveInfo.IdOrder), storage.Amout(reserveInfo.Amount))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	response := ResponseAcceptedOrder{
		Id:        reserve.IdUser,
		IdServise: reserve.IdService,
		IdOrder:   reserve.IdOrder,
		Amount:    reserve.Amount,
		Accepted:  true,
	}
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
