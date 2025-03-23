package data_structure

type Tree struct {
	max_resolution int
	Root           Node
}

func NewTree(max_resolution int, data BinData3D) Tree {
	tree := Tree{
		max_resolution: max_resolution,
		Root:           createNode(max_resolution, data, [3]int{0, 0, 0}),
	}
	return tree
}

func createNode(max_resolution int, data BinData3D, coords [3]int) Node {
	for _, d := range [3]int{data.X, data.Y, data.Z} {
		if d <= max_resolution {
			return Node{
				IsLeaf:       true,
				ContainsData: data.Any(),
				XBounds:      [2]int{coords[0], coords[0] + data.X},
				YBounds:      [2]int{coords[1], coords[1] + data.Y},
				ZBounds:      [2]int{coords[2], coords[2] + data.Z},
			}
		}
	}
	if data.All() || !data.Any() {
		return Node{
			IsLeaf:       true,
			ContainsData: data.Any(),
			XBounds:      [2]int{coords[0], coords[0] + data.X},
			YBounds:      [2]int{coords[1], coords[1] + data.Y},
			ZBounds:      [2]int{coords[2], coords[2] + data.Z},
		}
	}
	oc_data, oc_coords := splitOcs(data, coords)
	var child_nodes [8]*Node
	for i := 0; i < 8; i++ {
		child_node := createNode(max_resolution, oc_data[i], oc_coords[i])
		child_nodes[i] = &child_node
	}
	return Node{
		Children:     child_nodes,
		IsLeaf:       false,
		ContainsData: data.Any(),
		XBounds:      [2]int{coords[0], coords[0] + data.X},
		YBounds:      [2]int{coords[1], coords[1] + data.Y},
		ZBounds:      [2]int{coords[2], coords[2] + data.Z},
	}
}

func splitOcs(data BinData3D, coords [3]int) ([8]BinData3D, [8][3]int) {
	var split_data [8]BinData3D
	var split_coords [8][3]int

	extent_0 := [3]int{data.X / 2, data.Y / 2, data.Z / 2}
	extent_1 := [3]int{data.X - extent_0[0], data.Y - extent_0[1], data.Z - extent_0[2]}
	extent := [2][3]int{extent_0, extent_1}

	idx := 0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				split_extent := [3]int{extent[i][0], extent[j][1], extent[k][2]}
				complement_extent := [3]int{
					extent[1-i][0],
					extent[1-j][1],
					extent[1-k][2],
				}
				relative_coords := [3]int{
					i * complement_extent[0],
					j * complement_extent[1],
					k * complement_extent[2],
				}
				split_data[idx] = data.GetSlice(
					relative_coords[0], split_extent[0],
					relative_coords[1], split_extent[1],
					relative_coords[2], split_extent[2],
				)
				split_coords[idx] = [3]int{
					coords[0] + relative_coords[0],
					coords[1] + relative_coords[1],
					coords[2] + relative_coords[2],
				}
				idx++
			}
		}
	}
	return split_data, split_coords
}
