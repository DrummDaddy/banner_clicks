package main

import (
	"banner_clicks/internal/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	dbHost := getenv("DB_HOST", "db")
	dbPort := getenv("DB_PORT", "5432")
	dbUser := getenv("DB_USER", "postgres")
	dbPassword := getenv("DB_PASSWORD", "postgres")
	dbName := getenv("DB_NAME", "banner_clicks")
	sslMode := getenv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/counter/{bannerID}", handlers.CounterHandler(db)).Methods("GET")
	router.HandleFunc("/stats/{bannerID}", handlers.StatsHandler(db)).Methods("POST")

	fmt.Println("Starting server on :3000")
	log.Fatal(http.ListenAndServe(":3000", router))

}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
