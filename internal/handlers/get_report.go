package handlers

import (
	"encoding/csv"
	"encoding/json"
	"github.com/GeorgeShibanin/avito_test/internal/storage"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (h *HTTPHandler) HandleGetReport(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Disposition", `attachment; filename="report.csv"`)
	from := r.URL.Query().Get("from")
	if from == "" {
		http.Error(rw, "invalid query params", http.StatusBadRequest)
		return
	}
	to := r.URL.Query().Get("to")
	if to == "" {
		http.Error(rw, "invalid query params", http.StatusBadRequest)
		return
	}
	layout := "2006-01-02"
	date1, err := time.Parse(layout, from)
	if err != nil {
		http.Error(rw, "invalid query params", http.StatusBadRequest)
		return
	}
	date2, err := time.Parse(layout, to)
	if err != nil {
		http.Error(rw, "invalid query params", http.StatusBadRequest)
		return
	}

	report, err := h.storage.GetReport(r.Context(), storage.Date(date1), storage.Date(date2))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	//Save to csv file
	csvFile, err := os.Create("report.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()
	for _, record := range report {
		row := []string{record.IdServise, strconv.Itoa(record.TotalSumm)}
		if err := csvwriter.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	csvwriter.Flush()
	rawResponse, err := json.Marshal(report)
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
