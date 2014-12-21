package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDatabaseMemory(t *testing.T) {
	Convey("Given DatabaseMemory", t, func() {
		dbm := NewDatabaseMemory()
		Convey("It implements the database interface", func() {
			_, ok := interface{}(dbm).(Database)
			So(ok, ShouldBeTrue)
		})
	})
}
