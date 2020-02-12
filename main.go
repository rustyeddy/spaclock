package main

import (
	"flag"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// ============================ types ===============================
type Configuration struct {
	Addr       string // Addr:port
	Debug      bool   // turn on debugging
	Pubdir     string // The directory to publish
	SerialPort string // Must be provided if desired

	// MQTT Broker and topic(s)
	Broker string
	Topic  string
}

// ============================ Globals ===============================

// upgrader is used by the HTTP socket to establish a websocket
// connection with the client
var (
	config Configuration
	msgQ   chan string
	tempQ  chan string
	timeQ  chan time.Time
	dateQ  chan time.Time
)

// ============================ Init ===============================
func init() {
	msgQ = make(chan string)
	tempQ = make(chan string)
	timeQ = make(chan time.Time)
	dateQ = make(chan time.Time)

	flag.StringVar(&config.Addr, "addr", "0.0.0.0:2222", "Address:port default is :8000")
	flag.BoolVar(&config.Debug, "debug", false, "Turn on debugging")
	flag.StringVar(&config.Pubdir, "pubdir", "./pub", "The directory to publish")
	flag.StringVar(&config.SerialPort, "serial", "", "Default is no serial port")
	flag.StringVar(&config.Broker, "broker", "tcp://10.24.10.10:1883/topic/tempf", "Broker addr:port")
	flag.StringVar(&config.Topic, "topic", "/topic/tempf", "topic to subscribe to")
}

// ============================ Main ===============================
func main() {
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(4)

	go routes(&wg)                         // handle http
	go timers(&wg)                         // timed events
	go serial_loop(config.SerialPort, &wg) // serial port
	go mqtt_loop(config.Broker, &wg)       // MQTT

	wg.Wait() // wait for everything to return or die
	log.Info("SPA Clock all done!")
}
