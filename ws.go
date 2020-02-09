package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	upgrader *websocket.Upgrader
)

func init() {
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}

// ===================== Websocket =====================================
func handleUpgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Websocket Upgrade failed %v", err)
		return
	}
	defer conn.Close()

	// conn is a parameter to ensure the pointer does not change on
	// us a new client connects. The following go routine exists
	// when it recieves an error attempting to read from the connection.
	go func(conn *websocket.Conn) {
		for {
			var msg wsMessage
			if err = conn.ReadJSON(&msg); err != nil {
				log.Errorf("Failed reading message: %+v", err)
				return
			}
			log.Infof("ws recieved message: %+v", msg)

			// Do something with the message ...
		}
	}(conn)

	// Loop forever wating on the msgQ, when we recieve one (a string)
	// we'll wrap it in the single field JSON string and send it to
	// our client
	for {
		select {
		case message := <-msgQ:
			log.Debugf("msgQ %q", message)
			msg := &wsMessage{Message: message}
			if err := conn.WriteJSON(&msg); err != nil {
				log.Errorf("Websocket Write failed %v", err)
				return
			}
		case temp := <-tempQ:
			log.Debugf("tempQ %s", temp)
			tmp := &tempMessage{Tempf: temp}
			if err := conn.WriteJSON(&tmp); err != nil {
				log.Errorf("Websocket Write failed %v", err)
				return
			}

		case date := <-dateQ:
			log.Debugf("dateQ %s", date)
			d := &dateMessage{Time: date}
			if err := conn.WriteJSON(&d); err != nil {
				log.Errorf("Websocket Write failed %v", err)
				return
			}

		case clock := <-timeQ:
			log.Debugf("timeQ %s", clock)
			c := &clockMessage{Time: clock}
			if err := conn.WriteJSON(&c); err != nil {
				log.Errorf("Websocket Write failed %v", err)
				return
			}

		}
	}
}

