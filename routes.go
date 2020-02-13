package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// ============================================================================
// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// Start the router and register callbacks
func routes(wg *sync.WaitGroup) {
	defer wg.Done()

	// Setup the router
	router := mux.NewRouter()

	// New websocket connection requests
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

func handleMessage(w http.ResponseWriter, r *http.Request) {
	ok := true
	switch r.Method {

	case "GET":
		log.Warning("Todo GET")

	case "PUT", "POST":
		vars := mux.Vars(r)
		message := vars["message"]

		webQ <- NewTLV(tlvTypeMessage, len(message), message)
	default:
		log.Warning("handleMessage r.Method is not handled", r.Method)
		ok = false
	}
	json.NewEncoder(w).Encode(map[string]bool{"ok": ok})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {

	// TODO: run a check on http, rest, mqtt and websocket
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
