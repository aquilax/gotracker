package main

type Database interface {
	Init()
	GetPeersCountForHash(infoHash []byte) (int, error)
	GetPeersForHash(infoHash []byte, total, limit int) (*Peers, error)
	GetPeerByHashAndId(infoHash, peerId []byte) (*Peer, error)
	DeletePeer(peer *Peer) error
	NewPeer(client *Client) error
	UpdatePeer(peer *Peer) error
	UpdateLastAccess(peer *Peer) error
	Clean() error
	GetScrapeInfo(infoHash []byte) (*ScrapeList, error)
	GetStats() (int, int, int, error)
}
