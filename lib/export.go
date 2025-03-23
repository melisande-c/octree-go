package main

import (
	"C"
	"fmt"
	"unsafe"

	"github.com/melisande-c/octree-go/algo"
	"github.com/melisande-c/octree-go/data_structure"
)

func numpy2go(data *C.int, length C.int) []int {
	slice := unsafe.Slice(data, length)
	slice_cast := make([]int, length)
	for i, v := range slice {
		slice_cast[i] = int(v)
	}
	return slice_cast
}

func numpy2BinData3D(data *C.int, shape [3]int) data_structure.BinData3D {
	slice := numpy2go(data, C.int(shape[0]*shape[1]*shape[2]))
	slice_bool := make([]bool, len(slice))
	for i, v := range slice {
		slice_bool[i] = v != 0
	}
	return data_structure.BinData3D{
		Data: slice_bool, X: shape[0], Y: shape[1], Z: shape[2],
	}
}

//export FindMinDist
func FindMinDist(
	x_coord C.int,
	y_coord C.int,
	z_coord C.int,
	data *C.int,
	x_data_shape C.int,
	y_data_shape C.int,
	z_data_shape C.int,
	out_dist *C.double,
	x_out_loc *C.int,
	y_out_loc *C.int,
	z_out_loc *C.int,
) {
	coords := [3]C.int{x_coord, y_coord, z_coord}
	data_shape := [3]C.int{x_data_shape, y_data_shape, z_data_shape}
	out_loc := [3]*C.int{x_out_loc, y_out_loc, z_out_loc}
	var coords_cast [3]int
	for i, c := range coords {
		coords_cast[i] = int(c)
	}
	var data_shape_cast [3]int
	for i, s := range data_shape {
		data_shape_cast[i] = int(s)
	}
	bin_data := numpy2BinData3D(data, data_shape_cast)
	tree := data_structure.NewTree(1, bin_data)
	min_dist, min_loc := algo.FindMinLoc(tree, coords_cast)
	*out_dist = C.double(min_dist)
	for i, v := range min_loc {
		*out_loc[i] = C.int(v)
	}
}

//export HelloWorld
func HelloWorld() {
	fmt.Println("Hello World")
}

func main() {}
