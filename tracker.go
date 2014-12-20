package main

import (
	"fmt"
	"log"
	"net/http"
)

type Tracker struct {
	c  *Config
	db *Database
}

type appHandler func(http.ResponseWriter, *http.Request) error

func NewTracker() *Tracker {
	return &Tracker{}
}

func (t *Tracker) Run() {
	t.c = NewConfig()
	if err := t.c.Load(); err != nil {
		log.Panic(err)
	}

	http.Handle("/announce", appHandler(t.announceHandler))
	http.Handle("/scrape", appHandler(t.scrapeHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *Tracker) announceHandler(w http.ResponseWriter, r *http.Request) error {
	cl, err := NewClient(t.c, r)
	if err != nil {
		return err
	}
	total, err := t.db.getPeersCountForHash(cl.InfoHash)
	if total < 1 {
		// No peers found
	}
	w.Write([]byte(fmt.Sprintf("d8:intervali%de12:min intervali%de5:peers", t.c.AnnounceInterval, t.c.MinInterval)))
	peers, err := t.db.getPeersForHash(cl.InfoHash, total, t.c)
	if err != nil {
		return err
	}
	peers.getPeersBuffer(cl.Compact, cl.NoPeerId).WriteTo(w)
	w.Write([]byte("e"))

	if err := cl.Event(); err != nil {
		return err
	}
	return nil
}

func (t *Tracker) scrapeHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}
