package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDatabaseMemory(t *testing.T) {
	Convey("Given DatabaseMemory", t, func() {
		dbm := NewDatabaseMemory()
		dbm.Init()
		Convey("It implements the database interface", func() {
			_, ok := interface{}(dbm).(Database)
			So(ok, ShouldBeTrue)
		})
		Convey("Adding peers works", func() {
			dbm.NewPeer(&Peer{
				InfoHash: []byte("01234567890123456789"),
				ID:       []byte("90123456789012345678"),
				Compact:  []byte("89012345678901234567"),
				IP:       "127.0.0.1",
				Port:     123,
				State:    stateDownloading,
				Updated:  666,
			})
			// First peer adds one leecher
			seeders, leechers, torrents, err := dbm.GetStats()
			So(leechers, ShouldEqual, 1)
			So(seeders, ShouldEqual, 0)
			So(torrents, ShouldEqual, 1)
			So(err, ShouldBeNil)

			dbm.NewPeer(&Peer{
				InfoHash: []byte("01234567890123456789"),
				ID:       []byte("90123456789012345671"),
				Compact:  []byte("89012345678901234561"),
				IP:       "127.0.0.2",
				Port:     123,
				State:    stateSeeding,
				Updated:  666,
			})
			// Second peer adds one seeder for the same torrent
			seeders, leechers, torrents, err = dbm.GetStats()
			So(leechers, ShouldEqual, 1)
			So(seeders, ShouldEqual, 1)
			So(torrents, ShouldEqual, 1)
			So(err, ShouldBeNil)

			dbm.NewPeer(&Peer{
				InfoHash: []byte("01234567890123456786"),
				ID:       []byte("90123456789012345631"),
				Compact:  []byte("89012345678901234561"),
				IP:       "127.0.0.3",
				Port:     123,
				State:    stateSeeding,
				Updated:  666,
			})
			// Second peer adds one seeder for the same torrent
			seeders, leechers, torrents, err = dbm.GetStats()
			So(leechers, ShouldEqual, 1)
			So(seeders, ShouldEqual, 2)
			So(torrents, ShouldEqual, 2)
			So(err, ShouldBeNil)
			Convey("GetPeersCountForHash returns number of peers for hash", func() {
				count, err := dbm.GetPeersCountForHash([]byte("01234567890123456789"))
				So(count, ShouldEqual, 2)
				So(err, ShouldBeNil)
				count, err = dbm.GetPeersCountForHash([]byte("01234567890123456786"))
				So(count, ShouldEqual, 1)
				So(err, ShouldBeNil)
				count, err = dbm.GetPeersCountForHash([]byte("none"))
				So(count, ShouldEqual, 0)
				So(err, ShouldBeNil)
			})
			Convey("GetPeerListForHash returns list of peers for hash", func() {
				pl, err := dbm.GetPeerListForHash([]byte("01234567890123456789"), 10, 3)
				So(len(*pl), ShouldEqual, 2)
				So(err, ShouldBeNil)
				pl, err = dbm.GetPeerListForHash([]byte("none"), 10, 3)
				So(len(*pl), ShouldEqual, 0)
				So(err, ShouldBeNil)
			})
			Convey("GetPeerByHashAndId returns peer", func() {
				peer, err := dbm.GetPeerByHashAndId([]byte("01234567890123456789"), []byte("90123456789012345678"))
				So(peer, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})
}
