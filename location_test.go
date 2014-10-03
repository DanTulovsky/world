package world

import (
	"testing"

	termbox "github.com/nsf/termbox-go"
	. "github.com/smartystreets/goconvey/convey"
)

func genWorld() *World {
	// Setup
	s := &Settings{
		NewPeep:         1,
		MaxAge:          999,
		MaxPeeps:        2,
		RandomDeath:     0.0001,
		NewPeepMax:      2,
		NewPeepModifier: 1000,
		Size:            &Size{10, 10, 0, -10, -10, 0},
	}

	// Listen for input events on keyboard, required to test
	event_queue := make(chan termbox.Event)

	return NewWorld("Alpha1", *s, event_queue)
}

func TestUpdateGrid(t *testing.T) {
	w := genWorld()

	// Peep 1
	peep1, _ := w.NewPeep("", NewLocationXYZ(1, 0, 0))

	// Peep 2
	peep2, _ := w.NewPeep("", NewLocation())

	w.Show()

	Convey("Cannot move peep2 over peep1", t, func() {
		So(w.UpdateGrid(peep2, peep2.Location(), peep1.Location()), ShouldNotBeNil)
	})
	Convey("Can move peep2 to new location", t, func() {
		So(w.UpdateGrid(peep2, peep2.Location(), NewLocationXYZ(3, 0, 0)), ShouldBeNil)
	})
	Convey("Can move peep2 to same location", t, func() {
		So(w.UpdateGrid(peep2, peep2.Location(), peep2.Location()), ShouldBeNil)
	})
	Convey("Cannot move peep2 off the X grid", t, func() {
		So(w.UpdateGrid(peep2, peep2.Location(), NewLocationXYZ(w.settings.Size.MaxX+1, 0, 0)), ShouldNotBeNil)
		So(w.UpdateGrid(peep2, peep2.Location(), NewLocationXYZ(w.settings.Size.MinX-1, 0, 0)), ShouldNotBeNil)
	})
	Convey("Cannot move peep2 off the Y grid", t, func() {
		So(w.UpdateGrid(peep2, peep2.Location(), NewLocationXYZ(0, w.settings.Size.MaxY+1, 0)), ShouldNotBeNil)
		So(w.UpdateGrid(peep2, peep2.Location(), NewLocationXYZ(0, w.settings.Size.MinY-1, 0)), ShouldNotBeNil)
	})
	Convey("Cannot move peep2 off the Z grid", t, func() {
		So(w.UpdateGrid(peep2, peep2.Location(), NewLocationXYZ(0, 0, w.settings.Size.MaxZ+1)), ShouldNotBeNil)
		So(w.UpdateGrid(peep2, peep2.Location(), NewLocationXYZ(0, 0, w.settings.Size.MinZ-1)), ShouldNotBeNil)
	})
	Convey("Cannot move peep2 off the X grid into wall", t, func() {
		So(w.UpdateGrid(peep2, peep2.Location(), &Location{w.settings.Size.MaxX, 0, 0}), ShouldNotBeNil)
		So(w.UpdateGrid(peep2, peep2.Location(), NewLocationXYZ(w.settings.Size.MinX, 0, 0)), ShouldNotBeNil)
	})

}

func TestExisterIcon(t *testing.T) {
	w := genWorld()

	// Peep 1
	peep1, _ := w.NewPeep("red", NewLocation())
	peep1.age = w.settings.MaxAge/2 + 4

	Convey("Peep should show up as above mid-age", t, func() {
		So(w.ExisterIcon(peep1), ShouldEqual, 'R')
	})

	peep1.age = w.settings.MaxAge/2 - 4
	Convey("Peep should show up as below mid-age", t, func() {
		So(w.ExisterIcon(peep1), ShouldEqual, 'r')
	})
}

func TestLocationNeighbors(t *testing.T) {
	w := genWorld()

	loc1 := &Location{0, 0, 0}

	Convey("Origin should have 8 neightbors", t, func() {
		So(len(w.LocationNeighbors(loc1)), ShouldEqual, 8)
	})

	//loc2 := &Location{w.settings.Size.MinX, w.settings.Size.MinY, 0}
	//Convey("TopLeft should have 3 neightbors", t, func() {
	//	So(len(w.LocationNeighbors(loc2)), ShouldEqual, 3)
	//})

}

func TestLocation(t *testing.T) {
	expected := &Location{1, 2, 3}
	origin := &Location{}

	Convey("Expecting location(1,2,3) to exist", t, func() {
		So(NewLocationXYZ(1, 2, 3).SameAs(expected), ShouldBeTrue)
	})
	Convey("Expecting location(0,0,0) to exist", t, func() {
		So(NewLocation().SameAs(origin), ShouldBeTrue)
	})
}
