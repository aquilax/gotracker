package main

import (
	"log"
	"net/http"
)

type Tracker struct{}

type appHandler func(http.ResponseWriter, *http.Request) error

func NewTracker() *Tracker{
	return &Tracker{}
}

func (t *Tracker) Run() {
	http.Handle("/announce", appHandler(t.announceHandler))
	http.Handle("/scrape", appHandler(t.scrapeHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Println(err.Error());
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *Tracker) announceHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (t *Tracker) scrapeHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}
