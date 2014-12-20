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
	peers.getPeersBuffer(cl.IsCompact, cl.NoPeerId).WriteTo(w)
	w.Write([]byte("e"))

	if err := cl.processEvent(t.db); err != nil {
		return err
	}
	t.db.clean()
	return nil
}

func (t *Tracker) scrapeHandler(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Query().Get("stats") != "" {
		return t.handleStats()
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

func (t *Tracker) handleStats() error {
	// 	// get stats
	// 	$query = self::$db->query(
	// 		// select seeders and leechers
	// 		'SELECT SUM(state=1), SUM(state=0), ' .
	// 		// unique torrents from peers
	// 		'COUNT(DISTINCT info_hash) FROM peers;'
	// 	) OR tracker_error('failed to retrieve tracker statistics');
	// 	$stats = $query->fetchArray(SQLITE3_NUM);

	// 	// output format
	// 	switch ($_GET['stats'])
	// 	{
	// 		// xml
	// 		case 'xml':
	// 			header('Content-Type: text/xml');
	// 			echo '<?xml version="1.0" encoding="ISO-8859-1"?>' .
	// 			     '<tracker version="$Id: tracker.sqlite3.php 148 2009-11-16 23:18:28Z trigunflame $">' .
	// 			     '<peers>' . number_format($stats[0] + $stats[1]) . '</peers>' .
	// 			     '<seeders>' . number_format($stats[0]) . '</seeders>' .
	// 			     '<leechers>' . number_format($stats[1]) . '</leechers>' .
	// 			     '<torrents>' . number_format($stats[2]) . '</torrents></tracker>';
	// 			break;

	// 		// json
	// 		case 'json':
	// 			header('Content-Type: text/javascript');
	// 			echo '{"tracker":{"version":"$Id: tracker.sqlite3.php 148 2009-11-16 23:18:28Z trigunflame $",' .
	// 			     '"peers": "' . number_format($stats[0] + $stats[1]) . '",' .
	// 			     '"seeders":"' . number_format($stats[0]) . '",' .
	// 			     '"leechers":"' . number_format($stats[1]) . '",' .
	// 			     '"torrents":"' . number_format($stats[2]) . '"}}';
	// 			break;

	// 		// standard
	// 		default:
	// 			echo '<!doctype html><html><head><meta http-equiv="content-type" content="text/html; charset=UTF-8">' .
	// 			     '<title>PeerTracker: $Id: tracker.sqlite3.php 148 2009-11-16 23:18:28Z trigunflame $</title>' .
	// 			     '<body><pre>' . number_format($stats[0] + $stats[1]) .
	// 			     ' peers (' . number_format($stats[0]) . ' seeders + ' . number_format($stats[1]) .
	// 			     ' leechers) in ' . number_format($stats[2]) . ' torrents</pre></body></html>';
	// 	}

	// 	// cleanup
	// 	$query->finalize();
	// }
	return nil
}
