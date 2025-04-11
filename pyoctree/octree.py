import ctypes
import platform
from typing import Optional

import numpy as np
from numpy.typing import NDArray
import numpy.ctypeslib as npct

# load golib 
system = platform.system()
if system == "Windows":
    libpath = "./libgo.dll"
elif (system == 'Darwin') or (system=="Linux"):
    libpath = "./libgo.so"
else:
    raise RuntimeError(f"Unsupported system '{system}'.")

golib = ctypes.cdll.LoadLibrary(libpath)


array_1d_int = npct.ndpointer(dtype=np.int32, ndim=1, flags="CONTIGUOUS")

# Define function signatures
golib.NewOcTree.argtypes = [array_1d_int, ctypes.c_int, ctypes.c_int, ctypes.c_int]
golib.NewOcTree.restype = ctypes.c_void_p

golib.DeleteOcTree.argtypes = [ctypes.c_void_p]
golib.DeleteOcTree.restype = None

golib.FindMinDist.restype = None
golib.FindMinDist.argtypes = [
    ctypes.c_void_p,
    ctypes.c_int,
    ctypes.c_int,
    ctypes.c_int,
    ctypes.c_int,
    ctypes.c_int,
    ctypes.c_int,
    ctypes.c_double,
    ctypes.c_double,
    ctypes.c_double,
    ctypes.POINTER(ctypes.c_double),
    ctypes.POINTER(ctypes.c_int),
    ctypes.POINTER(ctypes.c_int),
    ctypes.POINTER(ctypes.c_int),
]


class Tree:
    def __init__(self, bin_data: NDArray[np.int_]):
        shape_c = (ctypes.c_int * 3)(*bin_data.shape)
        self.ptr = golib.NewOcTree(
            bin_data.flatten(), shape_c[0], shape_c[1], shape_c[2]
        )  # uintptr in Go

    def __del__(self):
        golib.DeleteOcTree(self.ptr)

    def find_min_dist(
        self,
        coords: tuple[int, int, int],
        offset: Optional[tuple[int, int, int]] = None,
        scaling: Optional[tuple[float, float, float]] = None,
    ) -> tuple[float, tuple[int, int, int]]:
        if offset is None:
            offset = (0, 0, 0)
        if scaling is None:
            scaling = (1, 1, 1)

        out_dist = ctypes.c_double(-1)
        out_dist_ptr = ctypes.pointer(out_dist)

        out_loc_x = ctypes.c_int(-1)
        out_loc_y = ctypes.c_int(-1)
        out_loc_z = ctypes.c_int(-1)
        out_loc_x_ptr = ctypes.pointer(out_loc_x)
        out_loc_y_ptr = ctypes.pointer(out_loc_y)
        out_loc_z_ptr = ctypes.pointer(out_loc_z)

        coords_c = (ctypes.c_int * 3)(*coords)
        offset_c = (ctypes.c_int * 3)(*offset)
        scaling_c = (ctypes.c_double * 3)(*scaling)
        golib.FindMinDist(
            self.ptr,
            coords_c[0],
            coords_c[1],
            coords_c[2],
            offset_c[0],
            offset_c[1],
            offset_c[2],
            scaling_c[0],
            scaling_c[1],
            scaling_c[2],
            out_dist_ptr,
            out_loc_x_ptr,
            out_loc_y_ptr,
            out_loc_z_ptr,
        )

        out_dist_python = out_dist.value
        out_loc_python = (out_loc_x.value, out_loc_y.value, out_loc_z.value)
        return (out_dist_python, out_loc_python)
