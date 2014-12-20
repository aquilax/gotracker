package main

import (
	"fmt"
	"log"
	"bytes"
	"net/http"
)

const version = "gotorrent"

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
	peers.getPeersBuffer(cl.IsCompact, cl.NoPeerId).WriteTo(w)
	w.Write([]byte("e"))

	if err := cl.processEvent(t.db); err != nil {
		return err
	}
	t.db.clean()
	return nil
}

func (t *Tracker) scrapeHandler(w http.ResponseWriter, r *http.Request) error {
	stats := r.URL.Query().Get("stats")
	if stats != "" {
		return t.handleStats(w, stats)
	}
	infoHash := []byte(r.URL.Query().Get("info_hash"))
	if len(infoHash) != 20 {
		if !t.c.FullScrape {
			return TrackerError{"full scrape disabled"}
		}
		infoHash = nil
	}
	w.Write([]byte("d5:filesd"))
	scrapeList, err := t.db.GetScrapeInfo(infoHash)
	if err != nil {
		return err
	}
	for _, scrapeItem := range *scrapeList {
		w.Write([]byte(fmt.Sprintf("20:%sd8:completei%de10:downloadedi%de10:incompletei%dee", scrapeItem.InfoHash, scrapeItem.Complete, scrapeItem.Downloaded, scrapeItem.Incomplete)))
	}
	w.Write([]byte("ee"))
	return nil
}

func (t *Tracker) handleStats(w http.ResponseWriter, statsType string) error {
	var b bytes.Buffer
	seeders, leechers, torrents, err := t.db.GetStats()
	if err != nil {
		return err
	}
	switch statsType {
	case "xml":
		b.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="ISO-8859-1"?><tracker version="%s"><peers>%d</peers><seeders>%d</seeders><leechers>%d</leechers><torrents>%d</torrents></tracker>`, version, seeders+leechers, seeders, leechers, torrents))
		w.Header().Set("Content-Type", "text/xml")
	case "json":
		b.WriteString(fmt.Sprintf(`{"tracker":{"version":"%s", "peers": %d, "seeders": %d, "leechers": %d, "torrents": %d}}`, version, seeders+leechers, seeders, leechers, torrents))
		w.Header().Set("Content-Type", "text/javascript")
	default: 
		b.WriteString(fmt.Sprintf(`<!doctype html><html><head><meta charset='utf-8'><title>%s</title><body><pre>%d peers (%d seeders + %d leechers) in %d torrents</pre></body></html>`, version, seeders+leechers, seeders, leechers, torrents))
	}
	b.WriteTo(w)
	return nil
}
