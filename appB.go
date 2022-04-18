package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"time"
)

func App_B() {
	log.Println("B app started ")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(httpServerAddr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := NatsGetData()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = fmt.Fprint(w, string(data))
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func NatsGetData() ([]byte, error) {
	nc, err := nats.Connect(NatsUrl)
	if err != nil {
		log.Println(err)
	}
	defer nc.Close()

	msg, err := nc.Request("Currency_Market", []byte("ping"), 2*time.Second)
	if err != nil {
		log.Println(err)
	}
	return msg.Data, err

}
