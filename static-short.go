package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var lut map[string]string
var lutMutex = &sync.Mutex{}
var lastRefreshTs int64

func loadLut() {
	f, err := os.Open("test.csv")
	if err != nil {
		return
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return
	}
	lutMutex.Lock()
	lut = make(map[string]string, len(lines))
	for _, line := range lines {
		lut[line[0]] = line[1]
	}
	lutMutex.Unlock()
}

func main() {
	loadLut()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got an %s request from %s: %s (%s)",
			r.Proto, r.RemoteAddr, r.URL.Path, r.Host)
		if r.URL.Path == "/_refresh" {
			if time.Now().Unix() > lastRefreshTs+600 {
				log.Printf("Refreshing...")
				loadLut()
				lastRefreshTs = time.Now().Unix()
				fmt.Fprintf(w, "Refreshed.")
			} else {
				log.Printf("Not refreshing...")
				w.WriteHeader(423)
				fmt.Fprintf(w, "I will only refresh every 10 minutes max.")
			}
			return
		}
		short := r.URL.Path
		if short[0] == '/' {
			short = short[1:]
		}
		url, ok := lut[short]
		if !ok {
			w.WriteHeader(404)
			fmt.Fprintf(w, "Couldn't find what you were looking for")
			return
		}
		http.Redirect(w, r, url, 307)
	})

	http.ListenAndServe(":8087", nil)
}
