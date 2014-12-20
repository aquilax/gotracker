package main

import (
	"net/http"
	"errors"
	"strconv"
)

type Client struct {
	Seeding int
	Compact bool
	NoPeerId bool
	Ip string
	NumWant int
}

func NewClient(c *Config, r *http.Request) (*Client, error) {
	cl := &Client{}

	// 20-bytes - info_hash
	// sha-1 hash of torrent metainfo
	infoHash := r.URL.Query().Get("info_hash")
	if len(infoHash) != 20 {
		return nil, errors.New("Bad info_hash: " + infoHash);
	}

	// 20-bytes - peer_id
	// client generated unique peer identifier
	peerId := r.URL.Query().Get("peer_id")
	if len(peerId) != 20 {
		return nil, errors.New("Bad peer_id: " + peerId);
	}

	// integer - port
	// port the client is accepting connections from
	port := r.URL.Query().Get("port")
	if _, err := strconv.Atoi(port); err != nil {
		return nil, &TrackerError{"client listening port is invalid"}
	}

	// integer - left
	// number of bytes left for the peer to download
	left := r.URL.Query().Get("left")
	lefti, err := strconv.Atoi(left)
	if err != nil {
		return nil, &TrackerError{"client data left field is invalid"}
	}
	cl.Seeding = 1
	if lefti > 0 {
		cl.Seeding = 0
	}

	// integer boolean - compact - optional
	// send a compact peer response
	// http://bittorrent.org/beps/bep_0023.html
	compact := r.URL.Query().Get("comapct")
	cl.Compact = false
	if compact == "1" || c.ForceCompact {
		cl.Compact = true
	}

	// integer boolean - no_peer_id - optional
	// omit peer_id in dictionary announce response
	noPeerId := r.URL.Query().Get("no_peer_id")
	cl.NoPeerId = false
	if noPeerId != "" {
		cl.NoPeerId = true
	}

	// string - ip - optional
	// ip address the peer requested to use
	ip := r.URL.Query().Get("ip")
	if ip != "" && c.ExternalIp {
		//TODO validate IP 
			// $_GET['ip'] = trim($_GET['ip'],'::ffff:');
			// if (!ip2long($_GET['ip'])) tracker_error('invalid ip, dotted decimal only');
	}

// // string - ip - optional
// // ip address the peer requested to use
// if (isset($_GET['ip']) && $_SERVER['tracker']['external_ip'])
// {
// 	// dotted decimal only
// 	$_GET['ip'] = trim($_GET['ip'],'::ffff:');
// 	if (!ip2long($_GET['ip'])) tracker_error('invalid ip, dotted decimal only');
// }
// // set ip to connected client
// elseif (isset($_SERVER['REMOTE_ADDR'])) $_GET['ip'] = trim($_SERVER['REMOTE_ADDR'],'::ffff:');
// // cannot locate suitable ip, must abort
// else tracker_error('could not locate clients ip');

	numwant := r.URL.Query().Get("numwant")
	cl.NumWant = c.DefaultPeers
	numwantI, err := strconv.Atoi(numwant)
	if err == nil {
		cl.NumWant = c.MaxPeers
		if numwantI < c.MaxPeers {
			cl.NumWant = numwantI
		}
	}
	return cl, nil
}