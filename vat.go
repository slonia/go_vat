package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error
var args map[string]string

func main() {
	extractArgs()
	setupConnection()
	setupServer()
	updateRates()
}

func extractArgs() {
	args = make(map[string]string)
	plainArgs := os.Args[1:]
	for _, el := range plainArgs {
		strs := strings.Split(el, "=")
		param, value := strs[0], strs[1]
		args[param] = value
	}
}

func argOrDefault(name string, def string) string {
	val, ok := args[name]
	if ok {
		return val
	} else {
		return def
	}
}

func setupConnection() {
	user := argOrDefault("user", "postgres")
	database := argOrDefault("database", "")
	password := argOrDefault("password", "")
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", user, database, password)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func setupServer() {
	http.HandleFunc("/", handler)
	port := argOrDefault("port", "8001")
	serverString := fmt.Sprintf("localhost:%s", port)
	log.Fatal(http.ListenAndServe(serverString, nil))
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
