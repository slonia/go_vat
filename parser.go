package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"gopkg.in/iconv.v1"
)

const base = "http://kantory.pl/kursy/usd/"

func updateRates() {
	resp, err := http.Get(base)
	if err != nil {
		log.Fatal(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	cd, err := iconv.Open("utf-8", "iso-8859-2")
	if err != nil {
		log.Fatal(err)
	}
	// Search for the title
	rows := scrape.FindAll(root, scrape.ByClass("ju12"))
	for _, row := range rows {

		cells := scrape.FindAll(row, scrape.ByTag(atom.Td))
		if len(cells) > 4 {
			buy := cd.ConvString(scrape.Text(cells[0]))
			sell := cd.ConvString(scrape.Text(cells[1]))
			bank_a, ok := scrape.Find(cells[4], scrape.ByTag(atom.A))
			var bank string
			if ok {
				bank = cd.ConvString(scrape.Text(bank_a))
			}
			set_at := cd.ConvString(scrape.Text(cells[5]))
			rows, err := db.Query("SELECT id FROM exchange_rates WHERE bank = $1 LIMIT 1", bank)
			if err != nil {
				log.Fatal(err)
			}
			var id string
			var found bool = false
			cur_time := time.Now()
			for rows.Next() {
				if err := rows.Scan(&id); err != nil {
					log.Fatal(err)
				}
				db.Exec(`UPDATE exchange_rates SET buy = $1, sell = $2, set_at = $3, updated_at = $4 WHERE id = $5;`, buy, sell, set_at, cur_time, id)
				fmt.Println("Bank exist")
				found = true
			}
			if !found {
				fmt.Println("Creating bank")

				_, err := db.Exec(`INSERT INTO exchange_rates (bank, buy, sell, set_at, created_at, updated_at)  VALUES ($1, $2, $3, $4, $5, $5);`, bank, buy, sell, set_at, cur_time)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
