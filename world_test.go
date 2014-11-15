package world

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListContains(t *testing.T) {
	Convey("ListContains works properly", t, func() {
		locations := []Location{Location{0, 1, 0}, Location{3, 1, 0}}
		loc := Location{0, 1, 0}
		missingLoc := Location{3, 3, 0}
		So(ListContains(locations, loc), ShouldBeTrue)
		So(ListContains(locations, missingLoc), ShouldBeFalse)
	})
}

func TestHomebase(t *testing.T) {
	w := genWorld()

	Convey("SetHomebase set and get match.", t, func() {
		loc := Location{3, 4, 0}
		w.SetHomebase("red", loc)

		peep1, _ := w.NewPeep("red", Location{})
		So(peep1.Homebase().SameAs(loc), ShouldBeTrue)
	})

}

func TestAliveDeadPeeps(t *testing.T) {
	w := genWorld() // World can have 361 peeps based on size
	w.settings.MaxPeeps = 4000
	w.settings.MaxAge = 4000

	Convey("No peeps in the world.", t, func() {
		So(w.AlivePeepCount(), ShouldEqual, 0)
		So(w.DeadPeepCount(), ShouldEqual, 0)
	})

	Convey("Accurate count of alive peeps.", t, func() {
		for x := 1; x < 100; x++ {
			loc, err := w.FindAnyEmptyLocation()
			So(err, ShouldBeNil)

			if _, err := w.NewPeep("", loc); err != nil {
				fmt.Println(err)
			}
			So(w.AlivePeepCount(), ShouldEqual, x)
		}
	})

	Convey("Accurate count of alive/dead peeps.", t, func() {
		deadPeeps := 0
		for _, peep := range w.peeps {
			So(w.DeadPeepCount(), ShouldEqual, deadPeeps)
			So(w.AlivePeepCount(), ShouldEqual, len(w.peeps)-deadPeeps)
			peep.Die(0)
			deadPeeps++

		}
	})

}
