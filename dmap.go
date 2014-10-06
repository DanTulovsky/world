// doublemap provides a 1 <-> map implementation of Exister <-> Location
package world

import (
	"fmt"
	"sync"
)

type dmap struct {
	mapExister  map[Exister]Location
	mapLocation map[Location]Exister
	mapLock     sync.RWMutex
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
	d.mapLock.RLock()
	defer d.mapLock.RUnlock()

	_, ok := d.mapExister[e]
	if !ok {
		return Location{}, fmt.Errorf("No such exister on the map!")
	}
	return d.mapExister[e], nil
}

// GetByLocation returns an Exister by location
func (d *dmap) GetByLocation(l Location) Exister {
	d.mapLock.RLock()
	defer d.mapLock.RUnlock()

	return d.mapLocation[l]
}

// Set sets an element
func (d *dmap) Set(e Exister, l Location) {
	d.mapLock.RLock()
	defer d.mapLock.RUnlock()

	d.mapExister[e] = l
	d.mapLocation[l] = e
}

// DelByExister removes an element by Exister
func (d *dmap) DelByExister(e Exister) {
	d.mapLock.Lock()
	defer d.mapLock.Unlock()

	l := d.mapExister[e]
	delete(d.mapExister, e)
	delete(d.mapLocation, l)
}

// DelByLocation removes an element by Location
func (d *dmap) DelByLocation(l Location) {
	d.mapLock.Lock()
	defer d.mapLock.Unlock()

	e := d.mapLocation[l]
	delete(d.mapExister, e)
	delete(d.mapLocation, l)
}

// AllNonEmptyLocations returns a list of all locations with an exister in them
func (d *dmap) AllNonEmptyLocations() []Location {
	keys := make([]Location, 0, len(d.mapLocation))
	for k := range d.mapLocation {
		if d.mapLocation[k] != nil {
			keys = append(keys, k)
		}
	}
	return keys
}
