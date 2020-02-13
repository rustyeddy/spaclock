package main

import (
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
		var tlv TLV
		var err error

		if n, tlv.Buffer, err = conn.ReadMessage(); n < 3 || err != nil {
			log.Errorf("Error reading TLV from websocket len %d, err %v", n, err)
			continue
		}

		log.Debugf("TLV type %v, len %v and value %v\n", tlv.Type(), tlv.Len(), tlv.Value())
	}
}

// wsWriter spins forever waiting on messages (TLVs) containing messages
// that need to be sent to the web socket client
func wsWriter(conn *websocket.Conn, readQ chan TLV) {

	// Loop forever wating on the msgQ, when we recieve one (a string)
	// we'll wrap it in the single field JSON string and send it to
	// our client
	for {
		var unknown int
		select {

		case tlv := <-webQ:

			log.Debugf("WS SEND: %d, len: %d, value: %s\n", tlv.Type(), tlv.Len(), tlv.Value())
			if err := conn.WriteMessage(websocket.BinaryMessage, tlv.Buffer); err != nil {
				log.Errorf("TLV write message failed %v", err)
				continue
			}

		default:
			unknown++
			continue
		}

	}
}
