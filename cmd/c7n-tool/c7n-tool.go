package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	http.HandleFunc("/c7n/acme-challenge", handleCheckDomain)
	http.HandleFunc("/", handleIndex)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Error(err)
	}
}

func handleCheckDomain(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s%s", r.Host, r.RequestURI)

	_, err := w.Write([]byte(url))
	if err != nil {
		log.Error(err)
	}
}
func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
