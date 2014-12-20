package main

import (
	"bytes"
	"fmt"
)

type Peer struct {
	InfoHash []byte `db:"info_hash"`
	ID       []byte `db:"peer_id"`
	Compact  []byte `db:"compact"`
	IP       string `db:"ip"`
	Port     int    `db:"port"`
	State    int    `db:"state"`
	Updated  int    `db:"updated"`
}

type Peers []*Peer

func (ps Peers) getPeersBuffer(compact, noPeerId bool) *bytes.Buffer {
	var result bytes.Buffer

	// compact announce
	if compact {
		// peers list
		var tb bytes.Buffer
		// build response
		for _, peer := range ps {
			tb.Write(peer.Compact)
		}
		// 6-byte compacted peer info
		result.WriteString(fmt.Sprintf("%d:", tb.Len))
		result.Write(tb.Bytes())
		return &result
	}

	// dictionary announce
	result.WriteString("l")
	for _, peer := range ps {
		if noPeerId {
			result.WriteString(fmt.Sprintf("d2:ip%d:%s7:peer id20:%s4:porti%dee", len(peer.IP), peer.IP, peer.ID, peer.Port))
		} else {
			result.WriteString(fmt.Sprintf("d2:ip%d:%s4:porti%dee", len(peer.IP), peer.IP, peer.Port))
		}
	}
	result.WriteString("e")
	return &result
}
