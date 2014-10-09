package world

import (
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
