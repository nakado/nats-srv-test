package main

type ValCurs struct {
	Date   string   `json:"Date"`
	Name   string   `json:"-"`
	Valute []Valute `json:"Valute"`
}

type Valute struct {
	NumCode  string `json:"-"`
	CharCode string `json:"CharCode"`
	Nominal  string `json:"-"`
	Name     string `json:"-"`
	Value    string `json:"Value"`
	ID       string `json:"-"`
}
