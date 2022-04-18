package main

import "github.com/nats-io/nats.go"

const URLcbr = "https://www.cbr.ru/scripts/XML_daily_eng.asp?date_req="

const httpServerAddr = ":80"

const NatsUrl = nats.DefaultURL

func main() {
	go App_A()
	go App_B()
	select {}
}
