package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

func setupDBMock() (sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("failed to create sqlmock: " + err.Error())
	}

	return mock, func() {
		db.Close()
	}
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
	db, _, _ := sqlmock.New()
	router.HandleFunc("/counter/{bannerID}", CounterHandler(db))
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected StatusNoContent, have %d", rr.Code)
	}
}

func TestCounterHandler_BadBannerID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/counter/notint", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	db, _, _ := sqlmock.New()
	router.HandleFunc("/counter/{bannerID}", CounterHandler(db))
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, have %d", rr.Code)
	}
}

func TestCounterHandler_DBErrorBegin(t *testing.T) {
	mock, cleanup := setupDBMock()
	defer cleanup()
	mock.ExpectBegin().WillReturnError(fmt.Errorf("DB error"))
	req := httptest.NewRequest(http.MethodPost, "/counter/1", nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	db, _, _ := sqlmock.New()
	router.HandleFunc("/counter/{bannerID}", CounterHandler(db))
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected StatusInternalServerError, have %d", rr.Code)
	}
}
