package main

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// ==================== handleTicker -===================================
func timers(wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case t := <-ticker.C:
			log.Debugln("ticker went off")
			tstr := t.Format("Mon Jan 2")
			msg := NewMessage("date", tstr)
			webQ <- *msg
		}
	}
}
