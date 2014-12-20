package main

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	stateDownloading = iota
	stateSeeding
)

type Client struct {
	IsCompact bool
	NoPeerId  bool
	NumWant   int
	Event     string
	Peer
}

func NewClient(c *Config, r *http.Request) (*Client, error) {
	cl := &Client{}

	// 20-bytes - info_hash
	// sha-1 hash of torrent metainfo
	cl.InfoHash = []byte(r.URL.Query().Get("info_hash"))
	if len(cl.InfoHash) != 20 {
		return nil, errors.New("Bad info_hash: " + string(cl.InfoHash))
	}

	// 20-bytes - peer_id
	// client generated unique peer identifier
	cl.ID = []byte(r.URL.Query().Get("peer_id"))
	if len(cl.ID) != 20 {
		return nil, errors.New("Bad peer_id: " + string(cl.ID))
	}

	// integer - port
	// port the client is accepting connections from
	port := r.URL.Query().Get("port")
	var err error
	if cl.Port, err = strconv.Atoi(port); err != nil {
		return nil, &TrackerError{"client listening port is invalid"}
	}

	// integer - left
	// number of bytes left for the peer to download
	left := r.URL.Query().Get("left")
	lefti, err := strconv.Atoi(left)
	if err != nil {
		return nil, &TrackerError{"client data left field is invalid"}
	}
	cl.State = stateDownloading
	if lefti > 0 {
		cl.State = stateSeeding
	}

	// integer boolean - compact - optional
	// send a compact peer response
	// http://bittorrent.org/beps/bep_0023.html
	compact := r.URL.Query().Get("comapct")
	cl.IsCompact = false
	if compact == "1" || c.ForceCompact {
		cl.IsCompact = true
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

	cl.Event = r.URL.Query().Get("event")

	return cl, nil
}

func (cl *Client) processEvent(db *Database) error {
	peer, err := db.getPeerByHashAndId(cl.InfoHash, cl.ID)
	if err != nil {
		return err
	}

	switch cl.Event {
	case "stopped":
		// remove peer
		if peer != nil {
			db.DeletePeer(peer)
		}
	case "completed":
		cl.State = stateSeeding
	case "started":
		// client continuing download
	default:
		// new user
		if peer == nil {
			db.NewPeer(cl)
		} else if peer.IP == cl.IP && peer.Port == cl.Port && peer.State == cl.State {
			db.UpdatePeer(peer)
		} else {
			db.UpdateLastAccess(peer)
		}
	}
	return nil
}
