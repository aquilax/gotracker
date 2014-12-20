package main

type Database struct {
}

func (db *Database) Init() {
	// sql := []string{
	// 	"PRAGMA synchronous = OFF;",
	// 	"PRAGMA journal_mode = MEMORY;",
	// 	"PRAGMA temp_store = MEMORY;",
	// }
}

func (db *Database) getPeersCountForHash(infoHash []byte) (int, error) {
	// TODO

	// // fetch peer total
	// $peer = self::$db->prepare(
	// 	// select a count of the number of peers that match the given info_hash
	// 	'SELECT COUNT(*) FROM (SELECT 1 FROM peers WHERE info_hash=:info_hash);'
	// );

	// // assign binary data
	// $peer->bindValue(':info_hash', $_GET['info_hash'], SQLITE3_BLOB);

	// // execute peer row count & cleanup
	// $success = $peer->execute() OR tracker_error('failed to select peer count');
	// $total = $success->fetchArray(SQLITE3_NUM);
	// $success->finalize();
	// $peer->close();

	return 1, nil
}

func (db *Database) getPeersForHash(infoHash []byte, total int, c *Config) (*Peers, error) {
	var peers Peers
	// // prepare query
	// $peer = self::$db->prepare(
	// 	// select
	// 	'SELECT ' .
	// 		// 6-byte compacted peer info
	// 		($_GET['compact'] ? 'compact ' :
	// 			// 20-byte peer_id
	// 			(!$_GET['no_peer_id'] ? 'peer_id, ' : '') .
	// 			// dotted decimal string ip, integer port
	// 			'ip, port '
	// 		) .
	// 	// from peers table matching info_hash
	// 	'FROM peers WHERE info_hash=:info_hash' .
	// 	// less peers than requested, so return them all
	// 	($total[0] <= $_GET['numwant'] ? ';' :
	// 		// if the total peers count is low, use SQL RANDOM
	// 		($total[0] <= $_SERVER['tracker']['random_limit'] ?
	// 			" ORDER BY RANDOM() LIMIT {$_GET['numwant']};" :
	// 			// use a more efficient but less accurate RANDOM
	// 			" LIMIT {$_GET['numwant']} OFFSET " .
	// 			mt_rand(0, ($total[0]-$_GET['numwant'])) . ';'
	// 		)
	// 	)
	// );

	return &peers, nil
}

func (db *Database) getPeerByHashAndId(infoHash, peerId []byte) (*Peer, error) {

	// // build peer query
	// $peer = self::$db->prepare(
	// 	// select a peer from the peers table that matches the given info_hash and peer_id
	// 	'SELECT ip, port, state FROM peers WHERE info_hash=:info_hash AND peer_id=:peer_id;'
	// );

	// // assign binary data
	// $peer->bindValue(':info_hash', $_GET['info_hash'], SQLITE3_BLOB);
	// $peer->bindValue(':peer_id', $_GET['peer_id'], SQLITE3_BLOB);

	// // execute peer select & cleanup
	// $success = $peer->execute() OR tracker_error('failed to select peer data');
	// $pState = $success->fetchArray(SQLITE3_NUM);
	// $success->finalize();
	// $peer->close();

	return &Peer{}, nil
}

func (db *Database) DeletePeer(peer *Peer) {

}

func (db *Database) NewPeer(client *Client) {

}

func (db *Database) UpdatePeer(peer *Peer) {

}

func (db *Database) UpdateLastAccess(peer *Peer) {

}
