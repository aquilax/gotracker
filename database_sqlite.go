package main

type DatabaseSQLite struct{}

func (dbsl *DatabaseSQLite) Init() {
	// sql := []string{
	// 	"PRAGMA synchronous = OFF;",
	// 	"PRAGMA journal_mode = MEMORY;",
	// 	"PRAGMA temp_store = MEMORY;",
	// }
}

func (dbsl *DatabaseSQLite) getPeersCountForHash(infoHash []byte) (int, error) {
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

func (dbsl *DatabaseSQLite) getPeersForHash(infoHash []byte, total, limit int) (*PeerList, error) {
	var peerList PeerList
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

	return &peerList, nil
}

func (dbsl *DatabaseSQLite) getPeerByHashAndId(infoHash, peerId []byte) (*Peer, error) {

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

func (dbsl *DatabaseSQLite) DeletePeer(peer *Peer) {

}

func (dbsl *DatabaseSQLite) NewPeer(peer *Peer) {

}

func (dbsl *DatabaseSQLite) UpdatePeer(peer *Peer) {

}

func (dbsl *DatabaseSQLite) UpdateLastAccess(peer *Peer) {

}

func (dbsl *DatabaseSQLite) Clean() {

	// // database cleanup
	// public static function clean()
	// {
	// 	// run cleanup once per announce interval
	// 	// check 'clean_idle_peers'% of the time to avoid excess queries
	// 	if (mt_rand(1, $_SERVER['tracker']['clean_idle_peers']) == 1)
	// 	{
	// 		// unix timestamp
	// 		$time = time();

	// 		// fetch last cleanup time
	// 		if (($last = self::$db->querySingle(
	// 			// select last cleanup from tasks
	// 			"SELECT value FROM tasks WHERE name='prune';"
	// 		) + 0) == 0)
	// 		{
	// 			self::$db->exec(
	// 				// begin query transaction
	// 				'BEGIN TRANSACTION; ' .
	// 				// set tasks value prune to current unix timestamp
	// 				"INSERT OR REPLACE INTO tasks VALUES('prune', {$time}); " .
	// 				// delete peers that have been idle too long
	// 				'DELETE FROM peers WHERE updated < ' .
	// 				// idle length is announce interval x 2
	// 				($time - ($_SERVER['tracker']['announce_interval'] * 2)) . '; ' .
	// 				// end transaction and commit
	// 				'COMMIT;'
	// 			) OR tracker_error('could not perform maintenance');
	// 		}
	// 		// prune idle peers
	// 		elseif (($last + $_SERVER['tracker']['announce_interval']) < $time)
	// 		{
	// 			self::$db->exec(
	// 				// begin query transaction
	// 				'BEGIN TRANSACTION; ' .
	// 				// set tasks value prune to current unix timestamp
	// 				"UPDATE tasks SET value={$time} WHERE name='prune'; " .
	// 				// delete peers that have been idle too long
	// 				'DELETE FROM peers WHERE updated < ' .
	// 				// idle length is announce interval x 2
	// 				($time - ($_SERVER['tracker']['announce_interval'] * 2)) . '; ' .
	// 				// end transaction and commit
	// 				'COMMIT;'
	// 			) OR tracker_error('could not perform maintenance');
	// 		}
	// 	}
	// }

}

func (dbsl *DatabaseSQLite) GetScrapeInfo(infoHash []byte) (*ScrapeList, error) {
	return nil, nil
}

func (dbsl *DatabaseSQLite) GetStats() (int, int, int, error) {
	return 0, 0, 0, nil
}
