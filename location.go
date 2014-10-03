package world

import (
	"fmt"
	"math"
	"math/rand"
	"unicode"

	termbox "github.com/nsf/termbox-go"
)

// Grid holds information about what object is in what cell in the world.
type Grid struct {
	size *Size

	// Each object is represented by a unique string
	objects *dmap
}

// Location specifies one coordinate in the world.
type Location struct {
	X int32
	Y int32
	Z int32
}

func (l *Location) String() string {
	return fmt.Sprintf("(%v, %v, %v)", l.X, l.Y, l.Z)
}

// SameAs returns true if two locations are the same
func (l *Location) SameAs(other *Location) bool {

	if l.X == other.X && l.Y == other.Y && l.Z == other.Z {
		return true
	}
	return false

}

// NewLocation returns a new location at origin
func NewLocation() *Location {
	return &Location{0, 0, 0}
}

// NewLocationXYZ returns a new location at x,y,z
func NewLocationXYZ(x, y, z int32) *Location {
	return &Location{x, y, z}
}

// Size specifies the world size
// termbox starts at (0, 0) upper-left.
type Size struct {
	MaxX int32
	MaxY int32
	MaxZ int32

	MinX int32
	MinY int32
	MinZ int32
}

// Exister defines an object that exists in the grid
type Exister interface {
	Location() (l *Location)
	SetLocation(l *Location) error
	ID() string
	Gender() PeepGender
	Age() PeepAge
	IsAlive() bool
}

// Mover defines an object that can move through the world.
// Extends Exister
type Mover interface {
	Exister
	MoveX(int32) error
	MoveY(int32) error
	MoveZ(int32) error
}

// ExisterFromID returns an exister based on id
func (w *World) ExisterFromID(id string) Exister {
	for _, p := range w.peeps {
		if p.ID() == id {
			return p
		}
	}
	return nil
}

// LocationNeighbors returns all neighboring locations to the given one
func (w *World) LocationNeighbors(l *Location) []*Location {
	neighbors := []*Location{}

	for _, x := range []int32{-1, 0, 1} {
		for _, y := range []int32{-1, 0, 1} {
			newLoc := NewLocationXYZ(l.X+x, l.Y+y, l.Z)
			if newLoc.SameAs(l) {
				continue // skip our own location
			}
			if err := w.CheckOutsideGrid(newLoc.X, newLoc.Y, newLoc.Z); err == nil {
				neighbors = append(neighbors, newLoc)
			}
		}
	}
	return neighbors
}

// FindEmptyLocation returns an empty location next to one of the provided locations or an error if not able to find one
func (w *World) FindEmptyLocation(locations ...*Location) (*Location, error) {
	for _, l := range locations {
		neighbors := w.LocationNeighbors(l)
		for _, n := range neighbors {
			if !w.IsOccupiedLocation(n) {
				return n, nil
			}
		}
	}
	return nil, fmt.Errorf("No available locations next to %v", locations)
}

// OfSpawnAge returns true of Exister is old enough to spawn
func (w *World) OfSpawnAge(e Exister) bool {
	return e.Age() >= w.settings.SpawnAge
}

// SameGenderSpawn makes a new peep next to one of the provided peeps, if they are of the same gender
func (w *World) SameGenderSpawn(left, right Exister) {
	Log("Checking123")
	if left.Gender() == right.Gender() {
		Log("Checking")
		if w.OfSpawnAge(left) && w.OfSpawnAge(right) {
			Log("Spawning!")
			newLocation, err := w.FindEmptyLocation(left.Location(), right.Location())
			if err == nil {
				if rand.Float64() < w.settings.SpawnProbability {
					w.NewPeep(left.Gender(), newLocation)
				}

			}
		}
	}
}

// DiffGenderSpawn makes a new peep next to one of the provided peeps, if they are of a different gender
func (w *World) DiffGenderSpawn(left, right Exister) {
	if left.Gender() != right.Gender() {
		newLocation, err := w.FindEmptyLocation(left.Location(), right.Location())
		if err == nil {
			if rand.Float64() < w.settings.SpawnProbability {
				w.NewPeep("", newLocation)
			}
		}
	}
}

// Meet is called when two Movers bump into each other
func (w *World) Meet(left, right Exister) {
	// If they are of the same gender, they spawn a new one (yes yes, I know it's backwards)
	w.SameGenderSpawn(left, right)

	// If they are of a different gender, they spawn a random child.
	//w.DiffGenderSpawn(left, right)

}

// IsOccupiedLocation returns True if the given Location is occupied by something alive
func (w *World) IsOccupiedLocation(l *Location) bool {
	e := w.grid.objects.GetByLocation(l)

	if e == nil {
		return false
	}
	if e.IsAlive() {
		return true
	}
	return false
}

// UpdateGrid updates a location on the world grid with the current occupant
// If cell is already occupied, call Meet function and return an error.
// A dead peep is not an occupant
func (w *World) UpdateGrid(m Mover, src *Location, dst *Location) error {
	if src == dst {

		// Set explicitly again to catch new peeps being created
		w.grid.objects.Set(m, src)
		return nil
	}
	// Check if someone else is already squatting here
	if w.IsOccupiedLocation(dst) {
		squatter := w.grid.objects.GetByLocation(dst)
		if squatter != nil && squatter.ID() != m.ID() {
			// We have a meeting, perhaps something happens here
			w.Meet(m, squatter)
			return fmt.Errorf("Location (%v) taken by: %v", src, squatter.ID())
		}
	}

	// Check that new location is inside the grid
	if err := w.CheckOutsideGrid(dst.X, dst.Y, dst.Z); err != nil {
		return err
	}

	w.grid.objects.DelByLocation(src)
	w.grid.objects.Set(m, dst)

	return nil
}

// CheckOutsideGrid return error if the move would place object outside grid.
func (w *World) CheckMovementOutsideGrid(src *Location, x, y, z int32) error {
	newX := src.X + x
	newY := src.Y + y
	newZ := src.Z + z
	return w.CheckOutsideGrid(newX, newY, newZ)
}

// CheckOutsideGrid returns error if coordinates are outside the grid
// X and Y also remove 1 line for border
func (w *World) CheckOutsideGrid(x, y, z int32) error {
	if x > w.settings.Size.MaxX-1 || x < w.settings.Size.MinX+1 {
		return fmt.Errorf("Movement outside world not allowed!")
	}
	if y > w.settings.Size.MaxY-1 || y < w.settings.Size.MinY+1 {
		return fmt.Errorf("Movement outside world not allowed!")
	}
	if z > w.settings.Size.MaxZ || z < w.settings.Size.MinZ {
		return fmt.Errorf("Movement outside world not allowed!")
	}

	return nil
}

// Move moves a mover in direction and magnitude specified.
// All three must succeed.
func (w *World) Move(m Mover, x, y, z int32) error {
	// Save current location
	src := m.Location()
	if err := w.CheckMovementOutsideGrid(src, x, y, z); err != nil {
		return err
	}

	if err := m.MoveX(x); err != nil {
		m.SetLocation(src)
		return err
	}
	if err := m.MoveY(y); err != nil {
		m.SetLocation(src)
		return err
	}
	if err := m.MoveZ(z); err != nil {
		m.SetLocation(src)
		return err
	}

	// On success, update the Grid
	dst := m.Location()
	if err := w.UpdateGrid(m, src, dst); err != nil {
		m.SetLocation(src)
		return err
	}
	return nil
}

// ExisterIcon returns the correct icon to use based on criteria defined
func (w *World) ExisterIcon(e Exister) rune {
	midAge := w.settings.MaxAge / 2

	// icon is the first character of gender
	icon := rune(e.Gender()[0])

	// UpperCase for those who reach middle age
	if e.Age() < midAge {
		return unicode.ToLower(icon)
	}
	return unicode.ToUpper(icon)
}

// ExisterFg returns the correct foreground color for an Exister
func (w *World) ExisterFg(e Exister) termbox.Attribute {
	switch e.Gender() {
	case "blue":
		return termbox.ColorBlue
	case "red":
		return termbox.ColorRed
	case "green":
		return termbox.ColorGreen

	case "yellow":
		return termbox.ColorYellow
	}
	return termbox.ColorDefault
}

// ExisterBg returns the correct background color for an Exister
func (w *World) ExisterBg(e Exister) termbox.Attribute {
	//midAge := w.settings.MaxAge / 2

	// Mid-age gets a different background
	//if e.Age() > midAge {
	//	return termbox.ColorMagenta
	//}
	return termbox.ColorDefault
}

// Visuals describe visual attributes for displaying an Exister
type Visuals struct {
	Char rune              // character displayed
	Fg   termbox.Attribute // foreground color
	Bg   termbox.Attribute // background color
}

// ExisterVisuals returns all the visuals for a given Exister
func (w *World) ExisterVisuals(e Exister) *Visuals {
	v := &Visuals{
		Char: '*',
		Fg:   termbox.ColorDefault,
		Bg:   termbox.ColorDefault,
	}

	v.Char = w.ExisterIcon(e)
	v.Fg = w.ExisterFg(e)
	v.Bg = w.ExisterBg(e)

	return v
}

func (w *World) Draw() {
	//w, h := termbox.Size()
	//Log(w, h)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w.DrawGrid()

	for _, loc := range w.grid.objects.AllNonEmptyLocations() {
		e := w.grid.objects.GetByLocation(loc)
		if !e.IsAlive() {
			continue
		}

		// Convert our coordinates to termbox
		termX := int(loc.X) + int(math.Abs(float64(w.settings.Size.MinX)))
		termY := int(loc.Y) + int(math.Abs(float64(w.settings.Size.MinY)))
		visuals := w.ExisterVisuals(e)

		termbox.SetCell(termX, termY, visuals.Char, visuals.Fg, visuals.Bg)
	}
	termbox.Flush()
}

// DrawGrid draws borders around the world
func (w *World) DrawGrid() {
	width, height := termbox.Size()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Origin
	termbox.SetCell(0, 0, ' ', termbox.ColorYellow, termbox.ColorYellow)

	// top line
	for x := 0; x <= width; x++ {
		termbox.SetCell(x, 0, ' ', termbox.ColorDefault, termbox.Attribute(255))
	}

	// bottom line
	for x := 0; x <= width; x++ {
		termbox.SetCell(x, height-1, ' ', termbox.ColorDefault, termbox.Attribute(255))
	}

	// left border
	for y := 0; y <= height; y++ {
		termbox.SetCell(0, y, ' ', termbox.ColorDefault, termbox.Attribute(255))
	}

	// right border
	for y := 0; y <= height; y++ {
		termbox.SetCell(width-1, y, ' ', termbox.ColorDefault, termbox.Attribute(255))
	}

	termbox.Flush()
}
