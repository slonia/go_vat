package main

type ExchangeRate struct {
	Bank  string  `json:"bank"`
	Buy   float32 `json:"buy"`
	Sell  float32 `json:"sell"`
	SetAt string  `json:"set_at"`
}
