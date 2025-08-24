package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

func setupDBMock() (sqlmock.Sqlmock, func()) {
	mockDB, mock, _ := sqlmock.New()
	db = mockDB
	return mock, func() { db.Close() }
}

func TestCounterHandler_OK(t *testing.T) {
	mock, cleanup := setupDBMock()

	defer cleanup()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO banner_clicks").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest(http.MethodPost, "/counter/1", nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/counter/{bannerID}", CounterHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected StatusNoContent, have %d", rr.Code)
	}
}

func TestCounterHandler_BadBannerID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/counter/notint", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/counter/{bannerID}", CounterHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, have %d", rr.Code)
	}
}

func TestCounterHandler_DBErrorBegin(t *testing.T) {
	oldDB := db
	defer func() { db = oldDB }()
	db = nil

	req := httptest.NewRequest(http.MethodPost, "/counter/1", nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/counter/{bannerID}", CounterHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected StatusInternalServerError, have %d", rr.Code)
	}
}
