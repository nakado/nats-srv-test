package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/nats-io/nats.go"
	"golang.org/x/text/encoding/charmap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var QuerryData []ValCurs

func App_A() {
	log.Println("A app started")
	if err := NatsServe(); err != nil {
		log.Printf("error: %v\n", err)
	}
	worker()
	nc, err := nats.Connect(NatsUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	nc.Subscribe("Currency_Market", func(m *nats.Msg) {
		b, _ := json.Marshal(QuerryData)
		nc.Publish(m.Reply, b)
	})
	select {}
}

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func GetParseData(now string) {
	if xmlBytes, err := getXML(URLcbr + now); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		result := new(ValCurs)
		d := xml.NewDecoder(bytes.NewReader(xmlBytes))
		d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
			switch charset {
			case "windows-1251":
				return charmap.Windows1251.NewDecoder().Reader(input), nil
			default:
				return nil, fmt.Errorf("unknown charset: %s", charset)
			}
		}
		err = d.Decode(&result)
		if err != nil {
			log.Fatalf("decode: %v", err)
		}

		data := ValCurs{Date: now}
		for _, info := range result.Valute {
			switch info.CharCode {
			case "USD":
				data.Valute = append(data.Valute,
					Valute{
						CharCode: info.CharCode,
						Value:    info.Value,
					})
			case "EUR":
				data.Valute = append(data.Valute,
					Valute{
						CharCode: info.CharCode,
						Value:    info.Value,
					})
			}
		}
		QuerryData = append(QuerryData, data)
	}
}
func worker() {
	go func() {
		GetParseData(time.Now().Local().Format("02/01/2006")) // first call
		ticker := time.NewTicker(time.Minute)
		done := make(chan bool)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				h, m, _ := time.Now().Clock()
				if m == 0 && (h == 0) {
					GetParseData(time.Now().Local().Format("02/01/2006"))

				}
			}
		}
	}()
}

func NatsServe() error {
	var stdout bytes.Buffer
	cmd := exec.Command("nats-server")
	cmd.Stdout = &stdout
	return cmd.Run()
}
