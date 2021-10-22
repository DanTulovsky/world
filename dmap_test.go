package world

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetFunctions(t *testing.T) {
	dm := NewDmap()

	e := &Peep{
		id:      "test1",
		isalive: true,
		gender:  "red",
	}
	bad := &Peep{
		id:      "bad",
		isalive: true,
		gender:  "red",
	}
	l := Location{1, 2, 3}
	notOccupied := Location{1, 2, 4}

	// Set maps explicitly
	dm.mapExister[e] = l
	dm.mapLocation[l] = e

	Convey("Exister with uid 'test1' exists at location (1, 2, 3) ", t, func() {
		loc, err := dm.GetByExister(e)
		So(err, ShouldBeNil)
		So(loc.SameAs(l), ShouldBeTrue)
	})

	Convey("Location (1, 2, 3) should contain Exister with uid 'test1'", t, func() {
		So(dm.GetByLocation(l), ShouldEqual, e)
	})

	Convey("Location (1, 2, 4) should not exist", t, func() {
		So(dm.GetByLocation(notOccupied), ShouldBeNil)
	})

	Convey("Exister bad should not exist", t, func() {
		loc, err := dm.GetByExister(bad)
		So(err, ShouldNotBeNil)
		So(loc.SameAs(Location{}), ShouldBeTrue)
	})
}

func TestSetFunctions(t *testing.T) {
	dm := NewDmap()

	e := &Peep{
		id:      "test1",
		isalive: true,
		gender:  "red",
	}
	l := Location{1, 2, 3}

	dm.Set(e, l)

	Convey("Exister test1 is in the map at location (1, 2, 3)", t, func() {
		loc, err := dm.GetByExister(e)
		So(err, ShouldBeNil)
		So(loc.SameAs(l), ShouldBeTrue)
		So(dm.GetByLocation(l), ShouldEqual, e)
	})
}

func TestDeleteFunctions(t *testing.T) {
	dm := NewDmap()

	e := &Peep{
		id:      "test1",
		isalive: true,
		gender:  "red",
	}
	l := Location{1, 2, 3}

	dm.Set(e, l)

	dm.DelByExister(e)
	Convey("Exister test1 and location (1,2,3) are nil", t, func() {
		loc, err := dm.GetByExister(e)
		So(err, ShouldNotBeNil)
		So(loc.SameAs(Location{}), ShouldBeTrue)
		So(dm.GetByLocation(l), ShouldBeNil)
	})

	dm.Set(e, l)
	dm.DelByLocation(l)
	Convey("Exister test1 and location (1,2,3) are nil", t, func() {
		loc, err := dm.GetByExister(e)
		So(err, ShouldNotBeNil)
		So(loc.SameAs(Location{}), ShouldBeTrue)
		So(dm.GetByLocation(l), ShouldBeNil)
	})
}

func TestAllNonEmptyLocations(t *testing.T) {
	dm := NewDmap()

	Convey("All locations empty", t, func() {
		So(dm.AllNonEmptyLocations(), ShouldBeEmpty)
	})

	e := &Peep{
		id:      "test1",
		isalive: true,
		gender:  "red",
	}
	l := Location{1, 2, 3}

	dm.Set(e, l)
	Convey("Location (1,2,3) in list.", t, func() {
		So(dm.AllNonEmptyLocations(), ShouldNotBeEmpty)
		fmt.Println(dm.AllNonEmptyLocations())
		So(dm.AllNonEmptyLocations(), ShouldResemble, []Location{l})
	})
}

func BenchmarkDmap(b *testing.B) {
	dm := NewDmap()

	bench := func() {
		uuid := uuid.New()
		e := &Peep{
			id:      uuid.String(),
			isalive: true,
			gender:  "red",
		}
		l := Location{rand.Int31n(100), rand.Int31n(100), rand.Int31n(100)}

		dm.Set(e, l)
		Convey("Exister exists at location", b, func() {
			loc, err := dm.GetByExister(e)
			So(err, ShouldBeNil)
			So(loc.SameAs(l), ShouldBeTrue)
		})
		dm.DelByExister(e)
		Convey("Exister and location are nil", b, func() {
			_, err := dm.GetByExister(e)
			So(err, ShouldNotBeNil)
			So(dm.GetByLocation(l), ShouldBeNil)
		})
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			go bench()
		}
	}
}
