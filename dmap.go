// doublemap provides a 1 <-> map implementation of Exister <-> Location
package world

import (
	"fmt"
	"sync"
)

type dmap struct {
	mapExister     map[Exister]Location
	mapExisterLock sync.RWMutex

	mapLocation     map[Location]Exister
	mapLocationLock sync.RWMutex
}

// NewDmap returns a new double map
func NewDmap() *dmap {
	return &dmap{
		mapExister:  make(map[Exister]Location),
		mapLocation: make(map[Location]Exister),
	}
}

// GetByExister returns a location by Exister
func (d *dmap) GetByExister(e Exister) (Location, error) {
	d.mapExisterLock.RLock()
	defer d.mapExisterLock.RUnlock()

	d.mapLocationLock.RLock()
	defer d.mapLocationLock.RUnlock()

	l := d.mapExister[e]

	// check that this exister is in fact at that location
	if d.mapLocation[l] != nil && d.mapLocation[l].ID() == e.ID() {
		return d.mapExister[e], nil
	}
	return l, fmt.Errorf("Exister does not exist!")
}

// GetByLocation returns an Exister by location
func (d *dmap) GetByLocation(l Location) Exister {
	d.mapExisterLock.RLock()
	defer d.mapExisterLock.RUnlock()

	d.mapLocationLock.RLock()
	defer d.mapLocationLock.RUnlock()

	return d.mapLocation[l]
}

// Set sets an element
func (d *dmap) Set(e Exister, l Location) {
	d.mapExisterLock.RLock()
	defer d.mapExisterLock.RUnlock()

	d.mapLocationLock.RLock()
	defer d.mapLocationLock.RUnlock()

	d.mapExister[e] = l
	d.mapLocation[l] = e
}

// DelByExister removes an element by Exister
func (d *dmap) DelByExister(e Exister) {
	d.mapExisterLock.RLock()
	defer d.mapExisterLock.RUnlock()

	d.mapLocationLock.RLock()
	defer d.mapLocationLock.RUnlock()

	l := d.mapExister[e]
	delete(d.mapExister, e)
	delete(d.mapLocation, l)
}

// DelByLocation removes an element by Location
func (d *dmap) DelByLocation(l Location) {
	d.mapExisterLock.RLock()
	defer d.mapExisterLock.RUnlock()

	d.mapLocationLock.RLock()
	defer d.mapLocationLock.RUnlock()

	e := d.mapLocation[l]
	delete(d.mapExister, e)
	delete(d.mapLocation, l)
}

// AllNonEmptyLocations returns a list of all locations with an exister in them
func (d *dmap) AllNonEmptyLocations() []Location {
	keys := make([]Location, 0, len(d.mapLocation))
	for k := range d.mapLocation {
		keys = append(keys, k)
	}
	return keys
}
