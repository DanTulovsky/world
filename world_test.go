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
		NewPeep:          0, // no randomness in tests
		MaxAge:           10,
		MaxPeeps:         20,
		RandomDeath:      0, // No randomness in tests
		NewPeepMax:       0, // No randomness in tests
		NewPeepModifier:  0, // no randomness in tests
		Size:             &Size{10, 10, 0, -10, -10, 0},
		SpawnAge:         5,
		SpawnProbability: 1, // No randomness in tests
		PeepViewDistance: 2,
		MaxGenders:       4,
	}

	// Listen for input events on keyboard, required to test
	event_queue := make(chan termbox.Event)

	// turn off random moves
	allowMoves = false

	return NewWorld("Alpha1", *s, event_queue, false)
}

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
}

func TestOnePeep(t *testing.T) {
	w := genWorld()
	maxAge := PeepAge(1000)
	w.settings.MaxPeeps = 1
	w.settings.MaxAge = maxAge
	w.settings.RandomDeath = 0

	peep1, _ := w.NewPeep("red", Location{})
	id := peep1.ID()

	for turn := 0; turn < int(maxAge); turn++ {
		Convey("One peep in the world.", t, func() {
			So(w.AlivePeepCount(), ShouldEqual, 1)
		})
		Convey("Peep's gender is 'red'", t, func() {
			So(peep1.Gender(), ShouldEqual, "red")
		})
		Convey("Peep's age is same as world age.", t, func() {
			So(peep1.Age(), ShouldEqual, turn)
		})
		Convey("Peep's ID has not changed!", t, func() {
			So(peep1.ID(), ShouldEqual, id)
		})
		w.NextTurn()
	}
}
