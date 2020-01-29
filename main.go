package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/tarm/serial"

	log "github.com/sirupsen/logrus"
)

// ============================ types ===============================
type Configuration struct {
	Addr       string // Addr:port
	Pubdir     string // The directory to publish
	SerialPort string // Must be provided if desired
}

type wsMessage struct {
	Message string `json:"message"`
}

// ============================ Globals ===============================

// upgrader is used by the HTTP socket to establish a websocket
// connection with the client
var (
	config   Configuration
	upgrader *websocket.Upgrader
	msgQ     chan string
)

// ============================ Init ===============================
func init() {
	msgQ = make(chan string)
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	flag.StringVar(&config.Addr, "addr", "0.0.0.0:8000", "Address:port default is :8000")
	flag.StringVar(&config.Pubdir, "pubdir", "./pub", "The directory to publish")
	flag.StringVar(&config.SerialPort, "serial", "", "Default is no serial port")
}

// ============================ Main ===============================
func main() {
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(2)
	go doRouter(&wg)
	if config.SerialPort != "" {
		go doSerial(config.SerialPort, &wg)
	}

	wg.Wait()
	log.Info("SPA Clock all done!")
}

// Start the router and register callbacks
func doRouter(wg *sync.WaitGroup) {
	defer wg.Done()

	// Setup the router
	router := mux.NewRouter()
	router.HandleFunc("/ws", handleUpgrade)
	router.HandleFunc("/api/health", handleHealth)
	router.HandleFunc("/api/message/{message}", handleMessage)

	spa := spaHandler{staticPath: config.Pubdir, indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    config.Addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}

// ============================================================================
// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ============ Handle HTTP REST and static file Requests =====================

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

func handleHealth(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		log.Warning("Todo GET")

	case "PUT", "POST":
		vars := mux.Vars(r)
		message := vars["message"]
		if message != "" {
			msgQ <- message
		} else {
			log.Info("\tmessage not found")
		}

	default:
		log.Warning("handleMessage DEFAULT")
	}
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
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
			msg := &wsMessage{Message: message}
			if err := conn.WriteJSON(&msg); err != nil {
				log.Errorf("Websocket Write failed %v", err)
				return
			}
		}
	}
}

// ===================== SerialPort =====================================
func doSerial(port string, wg *sync.WaitGroup) {
	defer wg.Done()

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
	var err error

	for i := 0; i < len(buf); i++ {
		if buf[i] == 0 || buf[i] == '\r' || buf[i] == '\n' {
			j := i - 1
			str = string(buf[0:j])
			break
		}
	}

	var f float64
	strs := strings.Split(str, "+")
	v := strings.Split(strs[2], ":")
	if f, err = strconv.ParseFloat(v[1], 32); err != nil {
		log.Errorf("failed to parse %v\n", err)
		return
	}
	fmt.Println("Floater: ", f)
}
