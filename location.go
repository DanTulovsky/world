package world

type Location interface {
	MoveX(int32)
	MoveY(int32)
	MoveZ(int32)
}
