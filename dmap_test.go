package world

import (
	"math/rand"
	"reflect"
	"testing"
	"code.google.com/p/go-uuid/uuid"

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
	empty := Location{0, 0, 0}

	// Set maps explicitly
	dm.mapExister[e] = l
	dm.mapLocation[l] = e

	Convey("Exister with uid 'test1' exists at location (1, 2, 3) ", t, func() {
		loc, err := dm.GetByExister(e)
		So(err, ShouldBeNil)
		So(reflect.DeepEqual(loc, l), ShouldBeTrue)
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
		So(reflect.DeepEqual(loc, empty), ShouldBeTrue)
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
		So(reflect.DeepEqual(loc, l), ShouldBeTrue)
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
	empty := Location{0, 0, 0}

	dm.Set(e, l)

	dm.DelByExister(e)
	Convey("Exister test1 and location (1,2,3) are nil", t, func() {
		loc, err := dm.GetByExister(e)
		So(err, ShouldNotBeNil)

		So(reflect.DeepEqual(loc, empty), ShouldBeTrue)
		So(dm.GetByLocation(l), ShouldBeNil)
	})

	dm.Set(e, l)
	dm.DelByLocation(l)
	Convey("Exister test1 and location (1,2,3) are nil", t, func() {
		loc, err := dm.GetByExister(e)
		So(err, ShouldNotBeNil)

		So(reflect.DeepEqual(loc, empty), ShouldBeTrue)
		So(dm.GetByLocation(l), ShouldBeNil)
	})
}

func BenchmarkDmap(b *testing.B) {
	dm := NewDmap()

	bench := func() {
		uuid := uuid.New()
		e := &Peep{
			id:      uuid,
			isalive: true,
			gender:  "red",
		}
		l := Location{rand.Int31n(100), rand.Int31n(100), rand.Int31n(100)}
		empty := Location{0, 0, 0}

		dm.Set(e, l)
		Convey("Exister exists at location", b, func() {
			loc, err := dm.GetByExister(e)
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(loc, l), ShouldBeTrue)
		})
		dm.DelByExister(e)
		Convey("Exister and location are nil", b, func() {
			loc, err := dm.GetByExister(e)
			So(err, ShouldNotBeNil)

			So(reflect.DeepEqual(loc, empty), ShouldBeTrue)
			So(dm.GetByLocation(l), ShouldBeNil)
		})
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			go bench()
		}
	}
}