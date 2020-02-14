package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	upgrader *websocket.Upgrader
)

type wsMessage struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

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

	go wsReader(conn)
	wsWriter(conn, webQ)
}

// conn is a parameter to ensure the pointer does not change on
// us a new client connects. The following go routine exists
// when it recieves an error attempting to read from the connection.
func wsReader(conn *websocket.Conn) {
	for {
		var n int
		var msg Message = Message{}
		var buf []byte
		var err error

		if n, buf, err = conn.ReadMessage(); err != nil {
			log.Errorf("Error reading TLV from websocket len %d, err %v", n, err)
			continue
		}
		err = json.Unmarshal(buf, &msg)
		if err != nil {
			log.Errorf("Error unmarshalling json %v", err)
			return
		}
		log.Debugf("Read JSON %+v\n", msg)
	}
}

// wsWriter spins forever waiting on messages (TLVs) containing messages
// that need to be sent to the web socket client
func wsWriter(conn *websocket.Conn, readQ chan Message) {

	// Loop forever wating on the msgQ, when we recieve one (a string)
	// we'll wrap it in the single field JSON string and send it to
	// our client
	for {
		var unknown int
		select {

		case msg := <-webQ:
			var buf []byte

			log.Debugf("WS Send JSON %+v", msg)
			if buf = msg.Marshal(); buf == nil {
				log.Errorln("msg.Marshal failed to JSONify the message")
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, buf); err != nil {
				log.Errorf("TLV write message failed %v", err)
				continue
			}
			log.Debugln("sent message...")

		default:
			unknown++
			continue
		}

	}
}
