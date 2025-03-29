package main

import (
	"fmt"

	"github.com/melisande-c/octree-go/algo"
	"github.com/melisande-c/octree-go/data"
)

func circle(size int, r int) data.BinData3D {
	img := make([]bool, size*size*size)
	c := size / 2
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				x := i - c
				y := j - c
				z := k - c
				img[i*size*size+j*size+k] = (x*x + y*y + z*z) < r*r
			}
		}
	}
	return data.BinData3D{Data: img, X: size, Y: size, Z: size}
}

func main() {
	bin_image := circle(32, 8)
	bin_slice := bin_image.GetSlice(0, 32, 0, 32, 16-4, 8)

	tree := data.NewTree(1, bin_slice)
	offset := [3]int{0, 0, 0}
	scaling := [3]float64{2, 2, 1}
	min_d, loc := algo.FindMinLoc(tree, [3]int{25, 3, 4}, offset, scaling)
	fmt.Println(min_d, loc)

	// n := 6
	// node := tree.Root.Children[n]
	// is_not_leaf := !node.IsLeaf
	// level := 1
	// for is_not_leaf {
	// 	found := false
	// 	for _, n := range node.Children {
	// 		if !n.IsLeaf {
	// 			node = n
	// 			found = true
	// 			break
	// 		}
	// 	}
	// 	if !found {
	// 		node = node.Children[7]
	// 	}
	// 	is_not_leaf = !node.IsLeaf
	// 	fmt.Println(
	// 		level,
	// 		node.XBounds,
	// 		node.YBounds,
	// 		node.ZBounds,
	// 		node.IsLeaf,
	// 		node.ContainsData,
	// 	)
	// 	level++
	// }
}
