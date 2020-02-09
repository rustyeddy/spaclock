package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"
	log "github.com/sirupsen/logrus"
)

// ============================ types ===============================
type Configuration struct {
	Addr       string // Addr:port
	Pubdir     string // The directory to publish
	SerialPort string // Must be provided if desired
	MQTTServer string
	MQTTTopic	string
}

type wsMessage struct {
	Message string `json:"message"`
}

type tempMessage struct {
	Tempf string `json:"tempf"`
}

type dateMessage struct {
	time.Time `json:"date"`
}

type clockMessage struct {
	time.Time `json:"clock"`
}

// ============================ Globals ===============================

// upgrader is used by the HTTP socket to establish a websocket
// connection with the client
var (
	config   Configuration
	msgQ  chan string
	tempQ chan string
	timeQ chan time.Time
	dateQ chan time.Time
)

// ============================ Init ===============================
func init() {
	msgQ = make(chan string)
	tempQ = make(chan string)
	timeQ = make(chan time.Time)
	dateQ = make(chan time.Time)

	flag.StringVar(&config.Addr, "addr", "0.0.0.0:2222", "Address:port default is :8000")
	flag.StringVar(&config.Pubdir, "pubdir", "./pub", "The directory to publish")
	flag.StringVar(&config.SerialPort, "serial", "", "Default is no serial port")
}

// ============================ Main ===============================
func main() {
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(4)

	go routes(&wg)
	go timers(&wg)
	go serial_loop(config.SerialPort, &wg)
	go mqtt_loop(config.MQTTServer, config.MQTTTopic, &wg)

	wg.Wait()
	log.Info("SPA Clock all done!")
}


func processBuffer(buf []byte) {
	var str string

	for i := 0; i < len(buf); i++ {
		if buf[i] == 0 || buf[i] == '\r' || buf[i] == '\n' {
			j := i - 1
			str = string(buf[0:j])
			break
		}
	}

	strs := strings.Split(str, "+")
	v := strings.Split(strs[2], ":")
	fmt.Printf("write floater to tempq: %s\n", v[1])
	tempQ <- v[1]
}
