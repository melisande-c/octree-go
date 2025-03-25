package data

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
	return d.Data[i*d.Y*d.Z+j*d.Z+k]
}

func (d *BinData3D) GetSlice(i, iExtent, j, jExtent, k, kExtent int) BinData3D {
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
	slice := make([]bool, iExtent*jExtent*kExtent)
	for l := i; l < i+iExtent; l++ {
		for m := j; m < j+jExtent; m++ {
			sd := l*d.Y*d.Z + m*d.Z
			s := (l-i)*jExtent*kExtent + (m-j)*kExtent
			copy(slice[s:s+kExtent], d.Data[sd+k:sd+k+kExtent])
		}
	}
	return BinData3D{Data: slice, X: iExtent, Y: jExtent, Z: kExtent}
}
