package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
)

func TestClientRequest(t *testing.T) {
	Convey("Given ClientRequest", t, func() {
		config := NewConfig()
		values, _ := url.ParseQuery("info_hash=12345678901234567890&peer_id=12345678901234567890&port=1&left=2")
		Convey("Returns new ClientRequest for valid request", func() {
			cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
			So(cr, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(string(cr.getPeer().InfoHash), ShouldEqual, "12345678901234567890")
			So(string(cr.getPeer().ID), ShouldEqual, "12345678901234567890")
			So(cr.getPeer().Port, ShouldEqual, 1)
			So(cr.getPeer().State, ShouldEqual, stateDownloading)
		})
		Convey("Returns error on invalid info_hash", func() {
			values.Set("info_hash", "123")
			cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
			So(cr, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Bad info_hash: 123")
		})
		Convey("Returns error on invalid peer_id", func() {
			values.Set("peer_id", "123")
			cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
			So(cr, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Bad peer_id: 123")
		})
		Convey("Returns error on invalid port", func() {
			values.Set("port", "sz")
			cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
			So(cr, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "d14:failure reason32:client listening port is invalide")
		})
		Convey("Test left", func() {
			Convey("Returns correct state if something is left", func() {
				values.Set("left", "1")
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.getPeer().State, ShouldEqual, stateDownloading)
			})
			Convey("Returns correct state if nothing is left", func() {
				values.Set("left", "0")
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.getPeer().State, ShouldEqual, stateSeeding)
			})
			Convey("Returns error if left is not a number", func() {
				values.Set("left", "z")
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "d14:failure reason33:client data left field is invalide")
			})
		})
		Convey("Test Compact", func() {
			Convey("No compact if 0", func() {
				values.Set("compact", "0")
				config.ForceCompact = false
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.isCompact, ShouldBeFalse)
			})
			Convey("No compact if empty", func() {
				values.Set("compact", "")
				config.ForceCompact = false
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.isCompact, ShouldBeFalse)
			})
			Convey("Compact if any", func() {
				values.Set("compact", "1")
				config.ForceCompact = false
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.isCompact, ShouldBeTrue)
			})
			Convey("Compact if forced", func() {
				values.Set("compact", "0")
				config.ForceCompact = true
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.isCompact, ShouldBeTrue)
			})
		})
		Convey("Sets noPeerId", func() {
			values.Set("no_peer_id", "1")
			cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
			So(cr, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(cr.noPeerId, ShouldBeTrue)
		})
		Convey("Test IP", func() {
			Convey("If correct IPv4 is passed", func() {
				values.Set("ip", "127.0.0.2")
				config.ExternalIp = true
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.getPeer().IP, ShouldEqual, "127.0.0.2")
			})
			Convey("If correct IPv6 is passed", func() {
				values.Set("ip", "2001:0db8:0000:0000:0000:ff00:0042:8329")
				config.ExternalIp = true
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.getPeer().IP, ShouldEqual, "2001:db8::ff00:42:8329")
			})
			Convey("Error if if Bad IP passed", func() {
				values.Set("ip", "bozo")
				config.ExternalIp = true
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "d14:failure reason31:invalid ip, dotted decimal onlye")
			})
			Convey("Use request ip if correct IPv4 is passed and external is disabled", func() {
				values.Set("ip", "127.0.0.2")
				config.ExternalIp = false
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.getPeer().IP, ShouldEqual, "127.0.0.1")
			})
			Convey("Error if if Bad request IP is passed", func() {
				values.Set("ip", "bozo")
				config.ExternalIp = false
				cr, err := NewClientRequest(config, values, "bozo:8000")
				So(cr, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "d14:failure reason31:invalid ip, dotted decimal onlye")
			})
			Convey("Error if if no ip is passed and request ip is empty", func() {
				values.Del("ip")
				config.ExternalIp = true
				cr, err := NewClientRequest(config, values, "")
				So(cr, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "d14:failure reason27:could not locate clients ipe")
			})
		})
		Convey("Numwant", func() {
			Convey("Get numwant from query string if below the limit", func() {
				values.Set("numwant", "2")
				config.MaxPeers = 3
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.numWant, ShouldEqual, 2)
			})
			Convey("Get max numwant if above the limit", func() {
				values.Set("numwant", "4")
				config.MaxPeers = 3
				config.DefaultPeers = 1
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.numWant, ShouldEqual, 3)
			})
			Convey("Get default numwant not set", func() {
				values.Del("numwant")
				config.MaxPeers = 3
				config.DefaultPeers = 2
				cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
				So(cr, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(cr.numWant, ShouldEqual, 2)
			})
		})
		Convey("Get event", func() {
			values.Set("event", "start")
			config.MaxPeers = 3
			config.DefaultPeers = 1
			cr, err := NewClientRequest(config, values, "127.0.0.1:8000")
			So(cr, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(cr.event, ShouldEqual, "start")
		})
	})
}
