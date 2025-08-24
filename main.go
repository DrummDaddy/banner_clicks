package main

import (
	"banner_clicks/internal/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	connStr := "user=postgres password=159357 dbname=banner_clicks sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)

	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/counter/{bannerID}", handlers.CounterHandler).Methods("GET")
	router.HandleFunc("/stats/{bannerID}", handlers.StatsHandler).Methods("POST")

	fmt.Println("Starting server on :3000")
	log.Fatal(http.ListenAndServe(":3000", router))

}
