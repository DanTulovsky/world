package world

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNextMoveToGetFromTo(t *testing.T) {
	w := genWorld()
	src := Location{0, 0, 0}
	dst := Location{1, 2, 0}

	Convey("Best magnitude to move from {0, 0, 0} to {1, 2, 0} is {1, 1, 0}", t, func() {
		x, y, z := w.NextMoveToGetFromTo(src, dst)
		So(x, ShouldEqual, 1)
		So(y, ShouldEqual, 1)
		So(z, ShouldEqual, 0)
	})

	src = Location{0, 0, 0}
	dst = Location{0, 0, 0}

	Convey("Best magnitude to move from {0, 0, 0} to {0, 0, 0} is {0, 0, 0}", t, func() {
		x, y, z := w.NextMoveToGetFromTo(src, dst)
		So(x, ShouldEqual, 0)
		So(y, ShouldEqual, 0)
		So(z, ShouldEqual, 0)
	})

	src = Location{w.MinX(), w.MinY(), 0}
	dst = Location{w.MaxX(), w.MaxY(), 0}

	Convey("Best magnitude to move from {minX, minY, 0} to {maxX, MaxY, 0} is {1, 1, 0}", t, func() {
		x, y, z := w.NextMoveToGetFromTo(src, dst)
		So(x, ShouldEqual, 1)
		So(y, ShouldEqual, 1)
		So(z, ShouldEqual, 0)
	})

	src = Location{w.MaxX(), w.MaxY(), 0}
	dst = Location{w.MinX(), w.MinY(), 0}

	Convey("Best magnitude to move from {maxX, maxY, 0} to {minX, MinY, 0} is {-1, -1, 0}", t, func() {
		x, y, z := w.NextMoveToGetFromTo(src, dst)
		So(x, ShouldEqual, -1)
		So(y, ShouldEqual, -1)
		So(z, ShouldEqual, 0)
	})
}

func TestNextMoveToGetFromToAlternatives(t *testing.T) {
	w := genWorld()
	w.NewPeep("red", NewLocationXYZ(0, 0, 0))
	w.NewPeep("red", NewLocationXYZ(1, 1, 0)) // in the way

	src := Location{0, 0, 0}
	dst := Location{1, 2, 0}

	Convey("Best magnitude to move from {0, 0, 0} to {1, 2, 0} is {0, 1, 0}", t, func() {
		x, y, z := w.NextMoveToGetFromTo(src, dst)
		So(x, ShouldEqual, 0)
		So(y, ShouldEqual, 1)
		So(z, ShouldEqual, 0)
	})
}

func TestNextMoveToGetAwayFrom(t *testing.T) {
	w := genWorld()
	w.NewPeep("red", NewLocationXYZ(w.MinX(), w.MaxY(), 0))

	src := Location{0, 0, 0}
	loc := Location{w.MinX(), w.MinY(), 0}

	Convey("Best magnitude to move away from {MinX, MaxY, 0} while on {0, 0, 0} is {1, 1, 0}", t, func() {
		x, y, z := w.NextMoveToGetAwayFrom(src, loc)
		So(x, ShouldEqual, 1)
		So(y, ShouldEqual, 1)
		So(z, ShouldEqual, 0)
	})
}
