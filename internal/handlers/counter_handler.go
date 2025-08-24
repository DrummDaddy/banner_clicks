package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bannerID, err := strconv.Atoi(vars["bannerID"])
	if err != nil {
		http.Error(w, "bad bannerID", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC().Truncate(time.Minute)

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
	INSERT INTO banner_clicks (ts, banner_id, cnt)
	VALUES($1, $2, 1)
	ON CONFLICT (ts, banner_id) DO UPDATE
	SET cnt = banner_clicks.cnt+1;
	`, now, bannerID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	tx.Commit()
	w.WriteHeader(http.StatusNoContent)

}
