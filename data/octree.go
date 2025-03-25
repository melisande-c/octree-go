package data

type OcTree struct {
	MaxRes int
	Root   OcNode
}

type OcNode struct {
	Children     [8]*OcNode
	IsLeaf       bool
	ContainsData bool
	XBounds      [2]int
	YBounds      [2]int
	ZBounds      [2]int
}

func NewTree(maxRes int, data BinData3D) OcTree {
	tree := OcTree{
		MaxRes: maxRes,
		Root:   createNode(maxRes, data, [3]int{0, 0, 0}),
	}
	return tree
}

func createNode(maxRes int, data BinData3D, coords [3]int) OcNode {
	allMaxRes := true
	for _, d := range [3]int{data.X, data.Y, data.Z} {
		allMaxRes = allMaxRes && d <= maxRes
	}
	if allMaxRes {
		return OcNode{
			IsLeaf:       true,
			ContainsData: data.Any(),
			XBounds:      [2]int{coords[0], coords[0] + data.X},
			YBounds:      [2]int{coords[1], coords[1] + data.Y},
			ZBounds:      [2]int{coords[2], coords[2] + data.Z},
		}
	}
	if data.All() || !data.Any() {
		return OcNode{
			IsLeaf:       true,
			ContainsData: data.Any(),
			XBounds:      [2]int{coords[0], coords[0] + data.X},
			YBounds:      [2]int{coords[1], coords[1] + data.Y},
			ZBounds:      [2]int{coords[2], coords[2] + data.Z},
		}
	}
	ocData, OcCoords := splitOcs(data, coords, maxRes)
	var childNodes [8]*OcNode
	for i := 0; i < 8; i++ {
		cn := createNode(maxRes, ocData[i], OcCoords[i])
		childNodes[i] = &cn
	}
	return OcNode{
		Children:     childNodes,
		IsLeaf:       false,
		ContainsData: data.Any(),
		XBounds:      [2]int{coords[0], coords[0] + data.X},
		YBounds:      [2]int{coords[1], coords[1] + data.Y},
		ZBounds:      [2]int{coords[2], coords[2] + data.Z},
	}
}

func splitOcs(data BinData3D, coords [3]int, maxRes int) ([8]BinData3D, [8][3]int) {
	var ocData [8]BinData3D
	var ocCoords [8][3]int

	e0 := [3]int{
		max(maxRes, data.X/2),
		max(maxRes, data.Y/2),
		max(maxRes, data.Z/2),
	}
	e1 := [3]int{data.X - e0[0], data.Y - e0[1], data.Z - e0[2]}
	extent := [2][3]int{e0, e1}

	idx := 0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				ocExtent := [3]int{extent[i][0], extent[j][1], extent[k][2]}
				complementExtent := [3]int{
					extent[1-i][0],
					extent[1-j][1],
					extent[1-k][2],
				}
				relCoords := [3]int{
					i * complementExtent[0],
					j * complementExtent[1],
					k * complementExtent[2],
				}
				ocData[idx] = data.GetSlice(
					relCoords[0], ocExtent[0],
					relCoords[1], ocExtent[1],
					relCoords[2], ocExtent[2],
				)
				ocCoords[idx] = [3]int{
					coords[0] + relCoords[0],
					coords[1] + relCoords[1],
					coords[2] + relCoords[2],
				}
				idx++
			}
		}
	}
	return ocData, ocCoords
}
