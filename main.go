package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func spin(period int) {
	b, err := json.Marshal(Message{fmt.Sprintf("Spinning for %d secs", period)})
	if err == nil {
		shipToLogstash(append(b, []byte("\n")...))
	}

	now := time.Now()
	end := now.Add(time.Duration(period) * time.Second)

	doSpin := func() {
		for time.Now().Before(end) {
		}
	}
	for i := 0; i < 8; i++ {
		go doSpin()
	}
}

type Message struct {
	Text string
}

func shipToLogstash(b []byte) error {
	logstashAddr := "localhost:5000"
	conn, err := net.Dial("tcp", logstashAddr)

	if err != nil {
		log.Fatalf("Dial failed: %s", err)
	}
	defer conn.Close()

	_, err = conn.Write(b)
	if err != nil {
		log.Println("Write to logstash failed:", err.Error())
	}
	return err
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!")
	m := r.URL.Query()
	if period, ok := m["period"]; ok {
		i, err := strconv.Atoi(period[0])
		if err != nil {
			return
		}
		go spin(i)
	} else {
		go spin(0)
	}
}

func main() {
	runtime.GOMAXPROCS(8)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
