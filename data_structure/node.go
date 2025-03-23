package data_structure

type Node struct {
	Children     [8]*Node
	IsLeaf       bool
	ContainsData bool
	XBounds      [2]int
	YBounds      [2]int
	ZBounds      [2]int
}
