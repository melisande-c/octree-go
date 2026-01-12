package data

type BinData3D interface {
	Any() bool
	All() bool
	Get(i, j, k int) uint16
	GetSlice(i, iExtent, j, jExtent, k, kExtent int) BinData3D
	GetShape() [3]int
}

type BinData3DShaped struct {
	Data    [][][]uint16
	X, Y, Z int // Dimensions
}

func (d *BinData3DShaped) GetShape() [3]int {
	return [3]int{d.X, d.Y, d.Z}
}

func (d *BinData3DShaped) Any() bool {
	for i := range d.X {
		for j := range d.Y {
			for k := range d.Z {
				v := d.Data[i][j][k]
				if v != 0 {
					return true
				}
			}
		}
	}
	return false
}

func (d *BinData3DShaped) All() bool {
	for i := range d.X {
		for j := range d.Y {
			for k := range d.Z {
				v := d.Data[i][j][k]
				if v == 0 {
					return false
				}
			}
		}
	}
	return true
}

func (d *BinData3DShaped) Get(i, j, k int) uint16 {
	return d.Data[i][j][k]
}

func (d *BinData3DShaped) GetSlice(i, iExtent, j, jExtent, k, kExtent int) BinData3D {
	view := make([][][]uint16, iExtent)
	for l, lidx := i, 0; l < i+iExtent; l, lidx = l+1, lidx+1 {
		view[lidx] = make([][]uint16, jExtent)
		for m, midx := j, 0; m < j+jExtent; m, midx = m+1, midx+1 {
			view[lidx][midx] = d.Data[l][m][k : k+kExtent]
		}
	}
	return &BinData3DShaped{
		Data: view,
		X:    iExtent,
		Y:    jExtent,
		Z:    kExtent,
	}
}

type BinData3DFlat struct {
	Data    []uint16
	X, Y, Z int // Dimensions
}

func (d *BinData3DFlat) GetShape() [3]int {
	return [3]int{d.X, d.Y, d.Z}
}

func (d *BinData3DFlat) Any() bool {
	for _, v := range d.Data {
		if v != 0 {
			return true
		}
	}
	return false
}

func (d *BinData3DFlat) All() bool {
	for _, v := range d.Data {
		if v == 0 {
			return false
		}
	}
	return true
}

// Get value at (i, j, k)
func (d *BinData3DFlat) Get(i, j, k int) uint16 {
	return d.Data[i*d.Y*d.Z+j*d.Z+k]
}

func (d *BinData3DFlat) GetSlice(i, iExtent, j, jExtent, k, kExtent int) BinData3D {
	// TODO: make sure that extents are not greater than d.X d.Y d.Z
	if i+iExtent >= d.X {
		iExtent = d.X - i
	}
	if j+jExtent >= d.Y {
		jExtent = d.Y - j
	}
	if k+kExtent >= d.Z {
		kExtent = d.Z - k
	}
	view := make([][][]uint16, iExtent)
	for l, lidx := i, 0; l < i+iExtent; l, lidx = l+1, lidx+1 {
		view[lidx] = make([][]uint16, jExtent)
		for m, midx := j, 0; m < j+jExtent; m, midx = m+1, midx+1 {
			sd := l*d.Y*d.Z + m*d.Z
			view[lidx][midx] = d.Data[sd+k : sd+k+kExtent]
		}
	}
	return &BinData3DShaped{
		Data: view,
		X:    iExtent,
		Y:    jExtent,
		Z:    kExtent,
	}
}
