package world

import (
	"fmt"
	"testing"

	termbox "github.com/nsf/termbox-go"
	. "github.com/smartystreets/goconvey/convey"
)

func genWorld() *World {
	// Setup
	s := &Settings{
		NewPeep:          1,
		MaxAge:           10,
		MaxPeeps:         20,
		RandomDeath:      0, // No randomness in tests
		NewPeepMax:       2,
		NewPeepModifier:  1000,
		Size:             &Size{10, 10, 0, -10, -10, 0},
		SpawnAge:         5,
		SpawnProbability: 1, // No randomness in tests
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

func TestMove(t *testing.T) {
	w := genWorld()

	peep1, _ := w.NewPeep("", NewLocationXYZ(3, 4, 0))

	w.Move(peep1, 1, 0, 0)

	Convey("peep1 moves to location 4, 4, 0", t, func() {
		loc1, _ := w.ExisterLocation(peep1)
		So(loc1.SameAs(Location{4, 4, 0}), ShouldBeTrue)
	})

}

func TestLocationExister(t *testing.T) {
	w := genWorld()
	peep1, _ := w.NewPeep("", NewLocationXYZ(3, 4, 0))

	Convey("peep1 is at location (3,4,0)", t, func() {
		So(w.LocationExister(Location{3, 4, 0}), ShouldEqual, peep1)
	})
}

func TestIsOccupiedLocation(t *testing.T) {
	w := genWorld()
	peep1, _ := w.NewPeep("", NewLocationXYZ(3, 4, 0))

	Convey("Location (1, 2, 0) is not occupied.", t, func() {
		So(w.IsOccupiedLocation(Location{1, 2, 0}), ShouldBeFalse)
	})

	Convey("Location (3, 4, 0) is occupied.", t, func() {
		So(w.IsOccupiedLocation(Location{3, 4, 0}), ShouldBeTrue)
	})

	peep1.Die(w.turn)
	Convey("Location (1, 2, 0) is not occupied.", t, func() {
		So(w.IsOccupiedLocation(Location{1, 2, 0}), ShouldBeFalse)
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

	loc1 := Location{0, 0, 0}

	Convey("Origin should have 8 neightbors", t, func() {
		So(len(w.LocationNeighbors(loc1)), ShouldEqual, 8)
	})

	loc2 := Location{w.MinX(), w.MinY(), 0}
	fmt.Println("Location: ", loc2)
	fmt.Println("Neighbors: ", w.LocationNeighbors(loc2))

	Convey("TopLeft should have 3 neightbors", t, func() {
		So(len(w.LocationNeighbors(loc2)), ShouldEqual, 3)
	})

}

func TestLocation(t *testing.T) {
	expected := Location{1, 2, 3}
	origin := Location{}

	Convey("Expecting location(1,2,3) to exist", t, func() {
		So(NewLocationXYZ(1, 2, 3).SameAs(expected), ShouldBeTrue)
	})
	Convey("Expecting location(0,0,0) to exist", t, func() {
		So(NewLocation().SameAs(origin), ShouldBeTrue)
	})
}

func TestSameAs(t *testing.T) {
	Convey("Locations are the same", t, func() {
		So(Location{1, 2, 3}.SameAs(Location{1, 2, 3}), ShouldBeTrue)
	})
	Convey("Locations are different", t, func() {
		So(Location{1, 2, 3}.SameAs(Location{1, 3, 3}), ShouldBeFalse)
	})
}

func TestMaxMin(t *testing.T) {
	w := genWorld()

	Convey("Correct Max values for the world.", t, func() {
		So(w.MaxX(), ShouldEqual, 9)
		So(w.MaxY(), ShouldEqual, 9)
		So(w.MaxZ(), ShouldEqual, 0)
		So(w.MinX(), ShouldEqual, -9)
		So(w.MinY(), ShouldEqual, -9)
		So(w.MinZ(), ShouldEqual, 0)
	})
}

// shouldBeInLocations checks that a list of Locations contains a given location
func shouldBeInLocations(loc interface{}, locList ...interface{}) string {
	for _, l := range locList {
		if l.(Location).SameAs(loc.(Location)) {
			return ""
		}
	}
	return "Error!"
}

func TestFindEmptyLocation(t *testing.T) {
	w := genWorld()

	peep1, _ := w.NewPeep("red", NewLocationXYZ(-9, -9, 0))
	loc, err := w.ExisterLocation(peep1)
	Convey("Location for peep1 found.", t, func() {
		So(err, ShouldBeNil)
	})

	// neighbors := w.LocationNeighbors(loc)

	emptyLoc, err := w.FindEmptyLocation(loc)
	Convey("Empty location found.", t, func() {
		So(err, ShouldBeNil)
		// So(emptyLoc, shouldBeInLocations, neighbors...)  // Why doesn't this work?
		So(emptyLoc, shouldBeInLocations, Location{-9, -8, 0}, Location{-8, -9, 0}, Location{-8, -8, 0})
	})

	// Fill up the neighbors
	for _, l := range w.LocationNeighbors(loc) {
		w.NewPeep("red", NewLocationXYZ(l.X, l.Y, l.Z))
	}

	emptyLoc, err = w.FindEmptyLocation(loc)
	Convey("Empty location not found.", t, func() {
		So(err, ShouldNotBeNil)
	})

}

func TestColors(t *testing.T) {
	w := genWorld()
	peep1, _ := w.NewPeep("red", NewLocation())

	Convey("Peeps is ColorRed", t, func() {
		So(w.ExisterFg(peep1), ShouldEqual, termbox.ColorRed)
	})

	Convey("Peeps is ColorRed", t, func() {
		So(w.ExisterBg(peep1), ShouldEqual, termbox.ColorDefault)
	})
}

func TestSpawnLocations(t *testing.T) {
	w := genWorld()

	locations := w.SpawnLocations()

	Convey("Correct spawn locations present.", t, func() {
		So(ListContains(locations, Location{9, 9, 0}), ShouldBeTrue)
		So(ListContains(locations, Location{-9, -9, 0}), ShouldBeTrue)
		So(ListContains(locations, Location{9, -9, 0}), ShouldBeTrue)
		So(ListContains(locations, Location{-9, 9, 0}), ShouldBeTrue)
	})
}

func TestSameGenderSpawn(t *testing.T) {

	w := genWorld()

	left, _ := w.NewPeep("red", Location{1, 1, 0})
	right, _ := w.NewPeep("red", Location{1, 0, 0})
	wrong, _ := w.NewPeep("blue", Location{1, -1, 0})

	Convey("Different genders don't spawn.", t, func() {
		So(w.SameGenderSpawn(left, wrong), ShouldNotBeNil)
	})

	Convey("Not of spawn age, no spawn.", t, func() {
		So(w.SameGenderSpawn(left, right), ShouldNotBeNil)
	})

	Convey("Same gender of spawn age spawn", t, func() {
		left.age = w.settings.SpawnAge + 1
		right.age = w.settings.SpawnAge + 1
		So(w.SameGenderSpawn(left, right), ShouldBeNil)
	})

	Convey("No empty location, no spawn", t, func() {
		// Populate all available locations around
		for _, loc := range []Location{Location{1, 1, 0}, Location{1, 0, 0}} {
			for _, l := range w.LocationNeighbors(loc) {
				w.NewPeep("red", l)
			}
		}
		So(w.SameGenderSpawn(left, right), ShouldNotBeNil)

	})
}

func TestMeet(t *testing.T) {
	w := genWorld()

	peep1, _ := w.NewPeep("red", Location{1, 2, 0})
	peep2, _ := w.NewPeep("red", Location{1, 1, 0})
	peep3, _ := w.NewPeep("red", Location{1, -1, 0})

	Convey("peep1 meets peep2 at turn 1.", t, func() {
		w.NextTurn()
		w.Meet(peep1, peep2)
		So(peep1.Met()[peep2], ShouldEqual, 1)
		So(peep2.Met()[peep1], ShouldEqual, 1)

		So(peep1.MetPeep(peep2), ShouldBeTrue)
		So(peep2.MetPeep(peep1), ShouldBeTrue)

		So(w.AlivePeeps(), ShouldEqual, 4)

		w.Meet(peep1, peep3)
		So(w.AlivePeeps(), ShouldEqual, 3)

		w.Meet(peep1, peep2)
		So(w.AlivePeeps(), ShouldEqual, 3)
	})
}
