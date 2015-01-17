package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPeer(t *testing.T) {
	Convey("Given Peer", t, func() {
		Convey("NewPeer returns new peer", func() {
			So(NewPeer, ShouldNotBeNil)
		})
		Convey("Given PeerList", func() {
			peerList := &PeerList{
				&Peer{
					ID:      []byte("12345678901234567890"),
					IP:      "127.0.0.1",
					Compact: []byte("23456789012345678901"),
					Port:    123,
				},
				&Peer{
					ID:      []byte("01234567890123456789"),
					IP:      "127.0.0.2",
					Compact: []byte("34567890123456789012"),
					Port:    124,
				},
			}
			Convey("Returns correct compact string", func() {
				result := peerList.getPeersBuffer(true, false)
				expected := "40:2345678901234567890134567890123456789012"
				So(result.String(), ShouldEqual, expected)
			})
			Convey("Returns correct noPeerId string", func() {
				result := peerList.getPeersBuffer(false, true)
				expected := "ld2:ip9:127.0.0.14:porti123eed2:ip9:127.0.0.24:porti124eee"
				So(result.String(), ShouldEqual, expected)
			})
			Convey("Returns correct PeerId string", func() {
				result := peerList.getPeersBuffer(false, false)
				expected := "ld2:ip9:127.0.0.17:peer id20:123456789012345678904:porti123eed2:ip9:127.0.0.27:peer id20:012345678901234567894:porti124eee"
				So(result.String(), ShouldEqual, expected)
			})

		})
	})
}
