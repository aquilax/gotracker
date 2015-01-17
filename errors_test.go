package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestErrors(t *testing.T) {
	Convey("Given tracker error", t, func() {
		err := TrackerError{"test error"}
		Convey("Generates correct error message", func() {
			So(err.Error(), ShouldEqual, "d14:failure reason10:test errore")
		})
	})
}
