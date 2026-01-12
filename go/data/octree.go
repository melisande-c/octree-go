package data

import "sync"

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

func (n *OcNode) Bounds() [2][3]int {
	var bounds [2][3]int
	for i, dimBounds := range [3][2]int{n.XBounds, n.YBounds, n.ZBounds} {
		for j, b := range dimBounds {
			bounds[j][i] = b
		}
	}
	return bounds
}

func NewTree(maxRes int, data BinData3D, root_offset [3]int) OcTree {
	tree := OcTree{
		MaxRes: maxRes,
		Root:   createRoot(maxRes, data, root_offset),
	}
	return tree
}

// TODO find a way to reduce duplication across createRoot and createNode
// createRoot is the same a createNode but creates the child nodes concurrently
func createRoot(maxRes int, data BinData3D, coords [3]int) OcNode {
	shape := data.GetShape()
	allMaxRes := true
	for _, d := range shape {
		allMaxRes = allMaxRes && d <= maxRes
	}
	if allMaxRes {
		return OcNode{
			IsLeaf:       true,
			ContainsData: data.Any(),
			XBounds:      [2]int{coords[0], coords[0] + shape[0]},
			YBounds:      [2]int{coords[1], coords[1] + shape[1]},
			ZBounds:      [2]int{coords[2], coords[2] + shape[2]},
		}
	}
	if data.All() || !data.Any() {
		return OcNode{
			IsLeaf:       true,
			ContainsData: data.Any(),
			XBounds:      [2]int{coords[0], coords[0] + shape[0]},
			YBounds:      [2]int{coords[1], coords[1] + shape[1]},
			ZBounds:      [2]int{coords[2], coords[2] + shape[2]},
		}
	}

	ocData, ocCoords := splitOcs(data, coords, maxRes)

	return OcNode{
		Children:     childNodesAsync(maxRes, ocData, ocCoords),
		IsLeaf:       false,
		ContainsData: data.Any(),
		XBounds:      [2]int{coords[0], coords[0] + shape[0]},
		YBounds:      [2]int{coords[1], coords[1] + shape[1]},
		ZBounds:      [2]int{coords[2], coords[2] + shape[2]},
	}
}

func childNodesAsync(maxRes int, ocData [8]BinData3D, ocCoords [8][3]int) [8]*OcNode {
	var childNodes [8]*OcNode
	var wg sync.WaitGroup
	wg.Add(8)
	resultChan := make(chan OcNode, 8)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for i := range 8 {
		go func() {
			defer wg.Done()
			resultChan <- createNode(maxRes, ocData[i], ocCoords[i])
		}()
	}
	i := 0
	for n := range resultChan {
		childNodes[i] = &n
		i++
	}
	return childNodes
}

func createNode(
	maxRes int, data BinData3D, coords [3]int,
) OcNode {
	shape := data.GetShape()
	allMaxRes := true
	for _, d := range shape {
		allMaxRes = allMaxRes && d <= maxRes
	}
	if allMaxRes {
		return OcNode{
			IsLeaf:       true,
			ContainsData: data.Any(),
			XBounds:      [2]int{coords[0], coords[0] + shape[0]},
			YBounds:      [2]int{coords[1], coords[1] + shape[1]},
			ZBounds:      [2]int{coords[2], coords[2] + shape[2]},
		}
	}
	if data.All() || !data.Any() {
		return OcNode{
			IsLeaf:       true,
			ContainsData: data.Any(),
			XBounds:      [2]int{coords[0], coords[0] + shape[0]},
			YBounds:      [2]int{coords[1], coords[1] + shape[1]},
			ZBounds:      [2]int{coords[2], coords[2] + shape[2]},
		}
	}

	ocData, ocCoords := splitOcs(data, coords, maxRes)
	var childNodes [8]*OcNode

	for i := range 8 {
		cn := createNode(maxRes, ocData[i], ocCoords[i])
		childNodes[i] = &cn
	}

	return OcNode{
		Children:     childNodes,
		IsLeaf:       false,
		ContainsData: data.Any(),
		XBounds:      [2]int{coords[0], coords[0] + shape[0]},
		YBounds:      [2]int{coords[1], coords[1] + shape[1]},
		ZBounds:      [2]int{coords[2], coords[2] + shape[2]},
	}
}

func splitOcs(data BinData3D, coords [3]int, maxRes int) ([8]BinData3D, [8][3]int) {
	shape := data.GetShape()
	var ocData [8]BinData3D
	var ocCoords [8][3]int

	e0 := [3]int{
		max(maxRes, shape[0]/2),
		max(maxRes, shape[1]/2),
		max(maxRes, shape[2]/2),
	}
	e1 := [3]int{shape[0] - e0[0], shape[1] - e0[1], shape[2] - e0[2]}
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
