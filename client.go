package main

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	InfoHash string
	Seeding  int
	Compact  bool
	NoPeerId bool
	IP       string
	NumWant  int
}

func NewClient(c *Config, r *http.Request) (*Client, error) {
	cl := &Client{}

	// 20-bytes - info_hash
	// sha-1 hash of torrent metainfo
	cl.InfoHash = r.URL.Query().Get("info_hash")
	if len(cl.InfoHash) != 20 {
		return nil, errors.New("Bad info_hash: " + cl.InfoHash)
	}

	// 20-bytes - peer_id
	// client generated unique peer identifier
	peerId := r.URL.Query().Get("peer_id")
	if len(peerId) != 20 {
		return nil, errors.New("Bad peer_id: " + peerId)
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
		pip := net.ParseIP(ip)
		if pip == nil {
			return nil, TrackerError{"invalid ip, dotted decimal only"}
		}
		cl.IP = pip.String()
	} else if r.RemoteAddr != "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
		pip := net.ParseIP(ip)
		if pip == nil {
			return nil, TrackerError{"invalid ip, dotted decimal only"}
		}
		cl.IP = pip.String()
	} else {
		return nil, TrackerError{"could not locate clients ip"}
	}

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

func (cl *Client) Event() error {

	// public static function event()
	// {
	// 	// build peer query
	// 	$peer = self::$db->prepare(
	// 		// select a peer from the peers table that matches the given info_hash and peer_id
	// 		'SELECT ip, port, state FROM peers WHERE info_hash=:info_hash AND peer_id=:peer_id;'
	// 	);

	// 	// assign binary data
	// 	$peer->bindValue(':info_hash', $_GET['info_hash'], SQLITE3_BLOB);
	// 	$peer->bindValue(':peer_id', $_GET['peer_id'], SQLITE3_BLOB);

	// 	// execute peer select & cleanup
	// 	$success = $peer->execute() OR tracker_error('failed to select peer data');
	// 	$pState = $success->fetchArray(SQLITE3_NUM);
	// 	$success->finalize();
	// 	$peer->close();

	// 	// process tracker event
	// 	switch ((isset($_GET['event']) ? $_GET['event'] : false))
	// 	{
	// 		// client gracefully exited
	// 		case 'stopped':
	// 			// remove peer
	// 			if (isset($pState[2])) self::delete_peer();
	// 			break;
	// 		// client completed download
	// 		case 'completed':
	// 			// force seeding status
	// 			$_SERVER['tracker']['seeding'] = 1;
	// 		// client started download
	// 		case 'started':
	// 		// client continuing download
	// 		default:
	// 			// new peer
	// 			if (!isset($pState[2])) self::new_peer();
	// 			// peer status
	// 			elseif (
	// 				// check that ip addresses match
	// 				$pState[0] != $_GET['ip'] ||
	// 				// check that listening ports match
	// 				($pState[1]+0) != $_GET['port'] ||
	// 				// check whether seeding status match
	// 				($pState[2]+0) != $_SERVER['tracker']['seeding']
	// 			) self::update_peer();
	// 			// update time
	// 			else self::update_last_access();
	// 	}
	// }
	return nil
}
