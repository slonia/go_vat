package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func main() {
	connStr := "user=postgres dbname=vat password=password sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	updateRates()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8001", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	rates := extractData()
	renderResponse(rates, w)
}
func extractData() []ExchangeRate {
	rows, err := db.Query("SELECT bank, buy, sell, set_at FROM exchange_rates")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var rates = []ExchangeRate{}
	for rows.Next() {
		var bank, set_at string
		var buy, sell float32
		if err := rows.Scan(&bank, &buy, &sell, &set_at); err != nil {
			log.Fatal(err)
		}
		rate := ExchangeRate{bank, buy, sell, set_at}
		rates = append(rates, rate)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return rates
}
func renderResponse(rates []ExchangeRate, w http.ResponseWriter) {
	js, err := json.Marshal(rates)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
