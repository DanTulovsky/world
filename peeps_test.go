package world

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewPeep(t *testing.T) {
	w := genWorld()

	Convey("NewPeep is born.", t, func() {
		peep1, err := w.NewPeep("red", Location{})
		So(err, ShouldBeNil)
		So(peep1.Gender(), ShouldEqual, "red")
	})
}

func TestMeetings(t *testing.T) {
	w := genWorld()

	peep1, _ := w.NewPeep("red", Location{1, 2, 0})
	peep2, _ := w.NewPeep("red", Location{1, 1, 0})

	Convey("peep1 meets peep2 at turn 1.", t, func() {
		peep1.Meet(peep2, 1)
		So(peep1.Met()[peep2], ShouldEqual, 1)
		// So(peep2.Met()[peep1], ShouldEqual, 1)

		So(peep1.MetPeep(peep2), ShouldBeTrue)
		// So(peep2.MetPeep(peep1), ShouldBeTrue)
	})
}

func TestIsAlive(t *testing.T) {
	w := genWorld()

	peep1, err := w.NewPeep("red", NewLocationXYZ(1, 1, 0))

	Convey("peep1 is alive", t, func() {
		So(err, ShouldBeNil)
		So(peep1.IsAlive(), ShouldBeTrue)
	})

	Convey("peep1 is 0 years old", t, func() {
		So(peep1.Age(), ShouldEqual, 0)
	})

	w.NextTurn()
	Convey("peep1 is 1 years old", t, func() {
		So(peep1.Age(), ShouldEqual, 1)
	})

	peep1.AddAge()
	Convey("peep1 is 2 years old", t, func() {
		So(peep1.Age(), ShouldEqual, 2)
	})

	Convey("peep1 is red", t, func() {
		So(peep1.Gender(), ShouldEqual, "red")
	})

	peep1.Die(w.turn)
	Convey("peep1 is dead", t, func() {
		So(peep1.IsAlive(), ShouldBeFalse)
	})

	Convey("peep1 is 0 years old", t, func() {
		So(peep1.IsAlive(), ShouldBeFalse)
	})

	peep2, err := w.NewPeep("red", NewLocationXYZ(3, 3, 0))
	Convey("peep2 is alive", t, func() {
		So(err, ShouldBeNil)
		So(peep2.IsAlive(), ShouldBeTrue)
	})

	for x := 0; x < int(w.settings.MaxAge+1); x++ {
		w.NextTurn()
	}
	Convey("peep2 is dead", t, func() {
		So(peep2.IsAlive(), ShouldBeFalse)
	})
}
