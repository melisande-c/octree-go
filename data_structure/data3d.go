package data_structure

type BinData3D struct {
	Data    []bool
	X, Y, Z int // Dimensions
}

func (d *BinData3D) Any() bool {
	for _, v := range d.Data {
		if v {
			return true
		}
	}
	return false
}

func (d *BinData3D) All() bool {
	for _, v := range d.Data {
		if !v {
			return false
		}
	}
	return true
}

// Get value at (i, j, k)
func (d *BinData3D) Get(i, j, k int) bool {
	return d.Data[i*d.X*d.Z+j*d.Z+k]
}

func (d *BinData3D) GetSlice(i, i_extent, j, j_extent, k, k_extent int) BinData3D {
	// TODO: make sure that extents are not greater than d.X d.Y d.Z
	if i+i_extent >= d.X {
		i_extent = d.X - i
	}
	if j+j_extent >= d.Y {
		j_extent = d.Y - j
	}
	if k+k_extent >= d.Z {
		k_extent = d.Z - k
	}
	slice_data := make([]bool, i_extent*j_extent*k_extent)
	index := 0
	for l := i; l < i+i_extent; l++ {
		for m := j; m < j+j_extent; m++ {
			for n := k; n < k+k_extent; n++ {
				slice_data[index] = d.Get(l, m, n)
				index++
			}
		}
	}
	return BinData3D{Data: slice_data, X: i_extent, Y: j_extent, Z: k_extent}
}
