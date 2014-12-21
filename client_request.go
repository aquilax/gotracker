package main

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
)

const (
	stateDownloading = iota
	stateSeeding
)

type ScrapeItem struct {
	InfoHash   []byte
	Complete   int
	Downloaded int
	Incomplete int
}

type ScrapeList []ScrapeItem

type ClientRequest struct {
	isCompact bool
	noPeerId  bool
	numWant   int
	event     string
	peer      *Peer
}

func NewClientRequest(c *Config, query url.Values, clientIp string) (*ClientRequest, error) {
	cr := &ClientRequest{
		peer: NewPeer(),
	}

	// 20-bytes - info_hash
	// sha-1 hash of torrent metainfo
	cr.peer.InfoHash = []byte(query.Get("info_hash"))
	if len(cr.peer.InfoHash) != 20 {
		return nil, errors.New("Bad info_hash: " + string(cr.peer.InfoHash))
	}

	// 20-bytes - peer_id
	// client generated unique peer identifier
	cr.peer.ID = []byte(query.Get("peer_id"))
	if len(cr.peer.ID) != 20 {
		return nil, errors.New("Bad peer_id: " + string(cr.peer.ID))
	}

	// integer - port
	// port the client is accepting connections from
	port := query.Get("port")
	var err error
	if cr.peer.Port, err = strconv.Atoi(port); err != nil {
		return nil, &TrackerError{"client listening port is invalid"}
	}

	// integer - left
	// number of bytes left for the peer to download
	left := query.Get("left")
	lefti, err := strconv.Atoi(left)
	if err != nil {
		return nil, &TrackerError{"client data left field is invalid"}
	}
	cr.peer.State = stateSeeding
	if lefti > 0 {
		cr.peer.State = stateDownloading
	}

	// integer boolean - compact - optional
	// send a compact peer response
	// http://bittorrent.org/beps/bep_0023.html
	compact := query.Get("compact")
	cr.isCompact = false
	if (compact != "" && compact != "0") || c.ForceCompact {
		cr.isCompact = true
	}

	// integer boolean - no_peer_id - optional
	// omit peer_id in dictionary announce response
	noPeerId := query.Get("no_peer_id")
	cr.noPeerId = false
	if noPeerId != "" && noPeerId != "0" {
		cr.noPeerId = true
	}

	// string - ip - optional
	// ip address the peer requested to use
	ip := query.Get("ip")
	if ip != "" && c.ExternalIp {
		pip := net.ParseIP(ip)
		if pip == nil {
			return nil, TrackerError{"invalid ip, dotted decimal only"}
		}
		cr.peer.IP = pip.String()
	} else if clientIp != "" {
		ip = strings.Split(clientIp, ":")[0]
		pip := net.ParseIP(ip)
		if pip == nil {
			return nil, TrackerError{"invalid ip, dotted decimal only"}
		}
		cr.peer.IP = pip.String()
	} else {
		return nil, TrackerError{"could not locate clients ip"}
	}

	numwant := query.Get("numwant")
	cr.numWant = c.DefaultPeers
	numwantI, err := strconv.Atoi(numwant)
	if err == nil {
		cr.numWant = c.MaxPeers
		if numwantI < c.MaxPeers {
			cr.numWant = numwantI
		}
	}

	cr.event = query.Get("event")

	return cr, nil
}

func (cr *ClientRequest) processEvent(db Database) error {
	peer, err := db.GetPeerByHashAndId(cr.peer.InfoHash, cr.peer.ID)
	if err != nil {
		return err
	}

	switch cr.event {
	case "stopped":
		// remove peer
		if peer != nil {
			db.DeletePeer(peer)
		}
	case "completed":
		cr.peer.State = stateSeeding
	case "started":
		// client continuing download
	default:
		// new user
		if peer == nil {
			db.NewPeer(cr.getPeer())
		} else if cr.differs(peer) {
			db.UpdatePeer(peer)
		} else {
			db.UpdateLastAccess(peer)
		}
	}
	return nil
}

func (cr *ClientRequest) getPeer() *Peer {
	return cr.peer
}

func (cr *ClientRequest) differs(peer *Peer) bool {
	return !(peer.IP == cr.peer.IP && peer.Port == cr.peer.Port && peer.State == cr.peer.State)
}
