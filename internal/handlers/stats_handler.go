package handlers

import (
	"banner_clicks/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bannerID, err := strconv.Atoi(vars["bannerID"])
	if err != nil {
		http.Error(w, "bad bannerID", http.StatusBadRequest)
		return
	}

	var req models.StatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	from, err := time.Parse(time.RFC3339, req.From)

	if err != nil {
		http.Error(w, "bad from format", http.StatusBadRequest)
		return
	}

	to, err := time.Parse(time.RFC3339, req.To)
	if err != nil {
		http.Error(w, "bad to format", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
	SELECT ts, cnt FROM banner_clicks
	WHERE banner_id = $1 
	AND ts >= $2 AND ts < $3 
	ORDER BY ts asc;
	`, bannerID, from, to)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	res := models.StatsResponse{}
	for rows.Next() {
		var ts time.Time
		var cnt int
		err := rows.Scan(&ts, &cnt)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		res.Stats = append(res.Stats, models.Stat{
			Ts: ts.Format(time.RFC3339),
			V:  cnt,
		})
	}

	w.Header().Set("Counter-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}
