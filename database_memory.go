package main

import (
	"sort"
)

type HashMap map[string]*PeerList

type Stats struct {
	seeders  int
	leechers int
}

type DatabaseMemory struct {
	hm    HashMap
	stats Stats
}

func (hm HashMap) find(infoHash []byte) *PeerList {
	peerList, found := hm[string(infoHash)]
	if !found {
		return &PeerList{}
	}
	return peerList
}

func (pl PeerList) Len() int           { return len(pl) }
func (pl PeerList) Less(i, j int) bool { return string(pl[i].ID) < string(pl[j].ID) }
func (pl PeerList) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

func (pl PeerList) findPeer(peerId []byte) int {
	if !sort.IsSorted(pl) {
		sort.Sort(pl)
	}
	return sort.Search(len(pl), func(i int) bool {
		return string(pl[i].ID) == string(peerId)
	})
}

func NewDatabaseMemory() *DatabaseMemory {
	return &DatabaseMemory{
		hm: make(HashMap),
	}
}

func (dbm *DatabaseMemory) Init() {
	dbm.stats.seeders = 0
	dbm.stats.leechers = 0
}

func (dbm *DatabaseMemory) GetPeersCountForHash(infoHash []byte) (int, error) {
	peerList := dbm.hm.find(infoHash)
	return len(*peerList), nil
}

func (dbm *DatabaseMemory) GetPeerListForHash(infoHash []byte, total, limit int) (*PeerList, error) {
	// TODO: honor limits, shuffle results
	peerList := dbm.hm.find(infoHash)
	return peerList, nil
}

func (dbm *DatabaseMemory) GetPeerByHashAndId(infoHash, peerId []byte) (*Peer, error) {
	peerList := dbm.hm.find(infoHash)
	if len(*peerList) > 0 {
		n := peerList.findPeer(peerId)
		if n < len(*peerList) {
			return (*peerList)[n], nil
		}
		return nil, nil
	}
	return nil, nil
}

func (dbm *DatabaseMemory) DeletePeer(peer *Peer) error {
	peerList := dbm.hm.find(peer.InfoHash)
	if len(*peerList) > 0 {
		n := peerList.findPeer(peer.ID)
		if n < len(*peerList) {
			pl := append((*peerList)[:n], (*peerList)[n+1:]...)
			dbm.hm[string(peer.InfoHash)] = &pl
		}
	}
	return nil
}

func (dbm *DatabaseMemory) NewPeer(peer *Peer) error {
	ih := string(peer.InfoHash)
	_, found := dbm.hm[ih]
	if !found {
		dbm.hm[ih] = &PeerList{}
	}
	pl := append(*(dbm.hm[ih]), peer)
	dbm.hm[ih] = &pl
	if peer.State == stateDownloading {
		dbm.stats.leechers++
	} else {
		dbm.stats.seeders++
	}
	return nil
}

func (dbm *DatabaseMemory) UpdatePeer(peer *Peer) error {
	return nil
}

func (dbm *DatabaseMemory) UpdateLastAccess(peer *Peer) error {
	return nil
}

func (dbm *DatabaseMemory) Clean() error {
	return nil
}

func (dbm *DatabaseMemory) GetScrapeInfo(infoHash []byte) (*ScrapeList, error) {
	var sl ScrapeList
	return &sl, nil
}

func (dbm *DatabaseMemory) GetStats() (int, int, int, error) {
	return dbm.stats.seeders, dbm.stats.leechers, len(dbm.hm), nil
}
