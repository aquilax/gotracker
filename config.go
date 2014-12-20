package main

type Config struct {
	// general tracker options
	OpenTracker      bool `json:"open_tracker"`      /* track anything announced to it */
	AnnounceInterval int  `json:"announce_interval"` /* how often client will send requests */
	MinInterval      int  `json:"min_interval"`      /* how often client can force requests */
	DefaultPeers     int  `json:"default_peers"`     /* default # of peers to announce */
	MaxPeers         int  `json:"max_peers"`         /* max # of peers to announce */

	// advanced tracker options
	ExternalIp     bool `json:"external_ip"`      /* allow client to specify ip address */
	ForceCompact   bool `json:"force_compact"`    /* force compact announces only */
	FullScrape     bool `json:"full_scrape"`      /* allow scrapes without info_hash */
	RandomLimit    int  `json:"random_limit"`     /* if peers > #, use alternate SQL RANDOM() */
	CleanIdlePeers int  `json:"clean_idle_peers"` /* tweaks % of time tracker attempts idle peer removal, */
	/* if you have a busy tracker, you may adjust this */
	/* example: 10 = 10%, 20 = 5%, 50 = 2%, 100 = 1% */

	// database options
	DbPath string `json:"db_path"` /* file path to the SQLite3 database*/
}

func NewConfig() *Config {
	return &Config{
		OpenTracker:      true,
		AnnounceInterval: 1800,
		MinInterval:      900,
		DefaultPeers:     50,
		MaxPeers:         100,

		ExternalIp:     true,
		ForceCompact:   false,
		FullScrape:     false,
		RandomLimit:    500,
		CleanIdlePeers: 10,

		DbPath: "./db/tracker.db",
	}
}

func (c *Config) Load() error {
	// TODO: Load config
	return nil
}
