package main

// #include <stdlib.h>

import (
	"C"
	"fmt"
	"unsafe"

	"github.com/melisande-c/octree-go/algo"
	"github.com/melisande-c/octree-go/data"
)

var treeRefs = make(map[uintptr]*data.OcTree)

func numpy2go(data *C.int, length C.int) []int {
	slice := unsafe.Slice(data, length)
	slice_cast := make([]int, length)
	for i, v := range slice {
		slice_cast[i] = int(v)
	}
	return slice_cast
}

func numpy2BinData3D(array *C.int, shape [3]int) data.BinData3D {
	slice := numpy2go(array, C.int(shape[0]*shape[1]*shape[2]))
	slice_bool := make([]bool, len(slice))
	for i, v := range slice {
		slice_bool[i] = v != 0
	}
	return data.BinData3D{
		Data: slice_bool, X: shape[0], Y: shape[1], Z: shape[2],
	}
}

//export NewOcTree
func NewOcTree(
	array *C.int,
	x_data_shape C.int,
	y_data_shape C.int,
	z_data_shape C.int,
) uintptr {
	data_shape := [3]C.int{x_data_shape, y_data_shape, z_data_shape}
	var data_shape_cast [3]int
	for i, s := range data_shape {
		data_shape_cast[i] = int(s)
	}
	bin_data := numpy2BinData3D(array, data_shape_cast)

	tree := data.NewTree(1, bin_data)
	tree_ref := &tree
	ptr := uintptr(unsafe.Pointer(tree_ref))
	treeRefs[ptr] = tree_ref

	return ptr
}

//export DeleteOcTree
func DeleteOcTree(ptr uintptr) {
	delete(treeRefs, ptr)
}

//export FindMinDist
func FindMinDist(
	ptr unsafe.Pointer,
	x_coord C.int,
	y_coord C.int,
	z_coord C.int,
	x_offset C.int,
	y_offset C.int,
	z_offset C.int,
	x_scaling C.double,
	y_scaling C.double,
	z_scaling C.double,
	out_dist *C.double,
	x_out_loc *C.int,
	y_out_loc *C.int,
	z_out_loc *C.int,
) {
	tree := (*data.OcTree)(ptr)
	coords := [3]int{int(x_coord), int(y_coord), int(z_coord)}
	offset := [3]int{int(x_offset), int(y_offset), int(z_offset)}
	scaling := [3]float64{float64(x_scaling), float64(y_scaling), float64(z_scaling)}

	out_loc := [3]*C.int{x_out_loc, y_out_loc, z_out_loc}

	min_dist, min_loc := algo.FindMinLoc(*tree, coords, offset, scaling)
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
