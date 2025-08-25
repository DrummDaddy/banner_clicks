package handlers

import (
	"banner_clicks/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

func TestStatsHandler_OK(t *testing.T) {
	db, mock, cleanup := setupDBMock()
	defer cleanup()

	rows := sqlmock.NewRows([]string{"ts", "cnt"}).
		AddRow(time.Date(2025, 1, 2, 15, 0, 0, 0, time.UTC), 5).
		AddRow(time.Date(2025, 1, 2, 16, 0, 0, 0, time.UTC), 7)

	mock.ExpectQuery("SELECT ts, cnt FROM banner_clicks").WillReturnRows(rows)

	payload := `{"from": "2025-01-02T00:00:00Z", "to":"2025-01-03T00:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/stats/1", strings.NewReader(payload))
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/stats/{bannerID}", StatsHandler(db))
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Errorf("expected successful response, have %d", rr.Code)
	}

	var resp models.StatsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal("couldn't decode JSON:", err)
	}
	if len(resp.Stats) != 2 {
		t.Errorf("expected 2 statistic records, have %d", len(resp.Stats))
	}
}

func TestStatsHandler_BadBannerID(t *testing.T) {
	payload := `{"from" : "2025-01-02T00:00:00Z", "to":"2025-01-03T00:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/stats/abc", strings.NewReader(payload))
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	db, _, _ := sqlmock.New()
	router.HandleFunc("/stats/{bannerID}", StatsHandler(db))
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, have %d", rr.Code)

	}
}

func TestStatsHandler_BadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/stats/1", strings.NewReader("notjson"))
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	db, _, _ := sqlmock.New()
	router.HandleFunc("/stats/{bannerID}", StatsHandler(db))
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, have %d", rr.Code)
	}
}

func TestStatsHandler_BadTimeFormat(t *testing.T) {
	payload := `{"from": "notdate", "to":"2025-01-03T00:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/stats/1", strings.NewReader(payload))
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	db, _, _ := sqlmock.New()
	router.HandleFunc("/stats/{bannerID}", StatsHandler(db))
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, have %d", rr.Code)
	}
}
