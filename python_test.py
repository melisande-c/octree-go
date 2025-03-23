import os

import numpy as np
import matplotlib.pyplot as plt
import tifffile

import ctypes

from numpy.typing import NDArray
import numpy.ctypeslib as npct

import ctypes

print(os.getcwd())

# golib = ctypes.CDLL("libgo.so")
# golib = npct.load_library("libgo_amd64.so", ".")
golib = ctypes.cdll.LoadLibrary("./libgo_amd64.so")

array_1d_int = npct.ndpointer(dtype=np.int32, ndim=1, flags="CONTIGUOUS")
golib.FindMinDist.restype = None
golib.FindMinDist.argtypes = [
    ctypes.c_int,
    ctypes.c_int,
    ctypes.c_int,
    array_1d_int,
    ctypes.c_int,
    ctypes.c_int,
    ctypes.c_int,
    ctypes.POINTER(ctypes.c_double),
    ctypes.POINTER(ctypes.c_int),
    ctypes.POINTER(ctypes.c_int),
    ctypes.POINTER(ctypes.c_int),
]


def find_min_dist(
    coords: tuple[int, int, int], bin_data: NDArray[np.int_]
) -> tuple[float, tuple[int, int, int]]:
    out_dist = ctypes.c_double(-1)
    out_dist_ptr = ctypes.pointer(out_dist)

    out_loc_x = ctypes.c_int(-1)
    out_loc_y = ctypes.c_int(-1)
    out_loc_z = ctypes.c_int(-1)
    out_loc_x_ptr = ctypes.pointer(out_loc_x)
    out_loc_y_ptr = ctypes.pointer(out_loc_y)
    out_loc_z_ptr = ctypes.pointer(out_loc_z)

    coords_c = (ctypes.c_int * 3)(*coords)
    shape_c = (ctypes.c_int * 3)(*bin_data.shape)

    golib.FindMinDist(
        coords_c[0],
        coords_c[1],
        coords_c[2],
        bin_data.flatten(),
        shape_c[0],
        shape_c[1],
        shape_c[2],
        out_dist_ptr,
        out_loc_x_ptr,
        out_loc_y_ptr,
        out_loc_z_ptr,
    )

    out_dist_python = out_dist.value
    out_loc_python = (out_loc_x.value, out_loc_y.value, out_loc_z.value)
    return (out_dist_python, out_loc_python)


# n = 128
# array = np.zeros((n, n, n), dtype=np.int32)
# ii, jj, kk = np.mgrid[:n, :n, :n]
# r = n / 4
# ii = ii - n/2
# jj = jj - n/2
# kk = kk - n/2
# array[(ii**2 + jj**2 + kk**2) < r**2] = 1

# array[6:10, 6:10, 16] = 1

img = tifffile.imread("/Users/melisande.croft/Documents/Data/5639253/Multilabel_U-Net_dataset_B.subtilis/training/instance_segmentation_GT/train_18.tif")
array = np.zeros((*img.shape, 1024), dtype=np.int32)
array[img != 0, 1024//2] = 1

query = (616, 313, 1024//2)
# query = (0, 0, 0 )
dist, loc = find_min_dist(query, array)
print(dist, loc)

fig, ax = plt.subplots()
# z = 0
z = array.shape[2] //2
# z = 8
ax.imshow(array[:,:,z])

ax.plot(query[1], query[0], "x", c="cyan")
ax.plot(loc[1], loc[0], "rx")

plt.show()
