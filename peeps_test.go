package world

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsAlive(t *testing.T) {
	w := genWorld()

	peep1, _ := w.NewPeep("red", NewLocation())

	Convey("peep1 is alive", t, func() {
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

	peep1.Die()
	Convey("peep1 is dead", t, func() {
		So(peep1.IsAlive(), ShouldBeFalse)
	})

	Convey("peep1 is 0 years old", t, func() {
		So(peep1.IsAlive(), ShouldBeFalse)
	})

	peep2, _ := w.NewPeep("red", NewLocation())
	Convey("peep2 is alive", t, func() {
		So(peep2.IsAlive(), ShouldBeTrue)
	})

	for x := 0; x < int(w.settings.MaxAge+1); x++ {
		w.NextTurn()
	}
	Convey("peep2 is dead", t, func() {
		So(peep2.IsAlive(), ShouldBeFalse)
	})
}
