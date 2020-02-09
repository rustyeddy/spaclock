package main

import (
	"fmt"
	"github.com/tarm/serial"
	"sync"
	"strings"
	log "github.com/sirupsen/logrus"
)


func serial_loop(port string, wg *sync.WaitGroup) {
	defer wg.Done()

	if port == "" {
		log.Info("No serial port configured, skipping...")
		return
	}

	c := &serial.Config{Name: port, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Errorln("open serial: ", err)
		return
	}

	for {
		buf := make([]byte, 256)
		_, err := s.Read(buf)
		if err != nil {
			log.Errorf("Read Error %v\n", err)
			continue
		}
		processBuffer(buf)
	}
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
