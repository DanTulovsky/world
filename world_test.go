package world

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpd(t *testing.T) {
	w := genWorld()

	// Peep 1
	peep1, _ := w.NewPeep("", NewLocationXYZ(1, 0, 0))

	// Peep 2
	peep2, _ := w.NewPeep("", NewLocation())

	w.Show()

	Convey("Cannot move peep2 over peep1", t, func() {
		loc1, _ := w.ExisterLocation(peep1)
		loc2, _ := w.ExisterLocation(peep2)
		So(w.UpdateGrid(peep2, loc2, loc1), ShouldNotBeNil)
	})
	Convey("Can move peep2 to new location", t, func() {
		loc2, _ := w.ExisterLocation(peep2)
		So(w.UpdateGrid(peep2, loc2, NewLocationXYZ(3, 0, 0)), ShouldBeNil)
	})
	Convey("Can move peep2 to same location", t, func() {
		loc2, _ := w.ExisterLocation(peep2)
		So(w.UpdateGrid(peep2, loc2, loc2), ShouldBeNil)
	})
	Convey("Cannot move peep2 off the X grid", t, func() {
		loc2, _ := w.ExisterLocation(peep2)
		So(w.UpdateGrid(peep2, loc2, NewLocationXYZ(w.MaxX()+1, 0, 0)), ShouldNotBeNil)
		So(w.UpdateGrid(peep2, loc2, NewLocationXYZ(w.MinX()-1, 0, 0)), ShouldNotBeNil)
	})
	Convey("Cannot move peep2 off the Y grid", t, func() {
		loc2, _ := w.ExisterLocation(peep2)
		So(w.UpdateGrid(peep2, loc2, NewLocationXYZ(0, w.MaxY()+1, 0)), ShouldNotBeNil)
		So(w.UpdateGrid(peep2, loc2, NewLocationXYZ(0, w.MinY()-1, 0)), ShouldNotBeNil)
	})
	Convey("Cannot move peep2 off the Z grid", t, func() {
		loc2, _ := w.ExisterLocation(peep2)
		So(w.UpdateGrid(peep2, loc2, NewLocationXYZ(0, 0, w.MaxZ()+1)), ShouldNotBeNil)
		So(w.UpdateGrid(peep2, loc2, NewLocationXYZ(0, 0, w.MinZ()-1)), ShouldNotBeNil)
	})
}
