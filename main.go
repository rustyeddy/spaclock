package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

type wsMessage struct {
	Message string `json:"message"`
	Clock   string `json:"clock"`
	Date    string `json:"date"`
}

// upgrader is used by the HTTP socket to establish a websocket
// connection with the client
var (
	upgrader *websocket.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func handleUpgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Websocket Upgrade failed %v", err)
		return
	}
	defer conn.Close()

	go wsLoop(conn)
}

func wsLoop(conn *websocket.Conn) {
	// ... Use conn to send a recieve messages: this basically turns
	// into an echo server.  We need to spawn this process off.  But
	// for now we'll just decode it and away we go.
	for {

		var err error
		var msg wsMessage

		if err = conn.ReadJSON(&msg); err != nil {
			log.Errorf("ws ReadJSON failed %v", err)
			return
		}
		log.Infof("ws recieved message: %+v", msg)

		msg = wsMessage{
			Message: "Monday, Jan. 13, 2020",
			Date:    "Today",
		}

		if err = conn.WriteJSON(msg); err != nil {
			log.Errorf("Websocket Write failed %v", err)
			return
		}
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/ws", handleUpgrade)
	router.HandleFunc("/api/health", handleHealth)
	router.HandleFunc("/api/message", handleMessage)

	spa := spaHandler{staticPath: "pub", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8003",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func handleMessage(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

	case "PUT", "POST":

	default:
	}
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
