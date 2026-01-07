import ctypes
import platform
from typing import Optional
from pathlib import Path

import numpy as np
from numpy.typing import NDArray
import numpy.ctypeslib as npct

# load golib
system = platform.system()
if (system != "Darwin") and (system != "Linux"):
    raise RuntimeError(
        f"Unsupported system '{system}'. Currently only Linux and Darwin is supported."
    )
machine = platform.machine()
if machine == "x86_64":
    machine = "amd64"
if machine not in ["arm64", "amd64"]:
    raise RuntimeError(
        f"Unsupported CPU '{machine}'. Currently only arm64 and amd64 is supported."
    )


libdir = Path(__file__).parent.resolve() / "_lib"
libfile = f"octree-{system.lower()}-{machine}.so"
libpath = libdir / libfile
print("File:", libpath)

go_octree = ctypes.cdll.LoadLibrary(libpath)


array_1d_int = npct.ndpointer(dtype=np.int32, ndim=1, flags="CONTIGUOUS")

# Define function signatures
go_octree.NewOcTree.argtypes = [
    array_1d_int,
    ctypes.c_int,
    ctypes.c_int,
    ctypes.c_int,
    ctypes.c_int,
    ctypes.c_int,
    ctypes.c_int,
]
go_octree.NewOcTree.restype = ctypes.c_void_p

go_octree.DeleteOcTree.argtypes = [ctypes.c_void_p]
go_octree.DeleteOcTree.restype = None

go_octree.FindMinDist.restype = None
go_octree.FindMinDist.argtypes = [
    ctypes.c_void_p,
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
    def __init__(
        self, bin_data: NDArray[np.int_], root_offset: Optional[tuple[int, int, int]] = None
    ):
        if root_offset is None:
            root_offset = (0, 0, 0)
        shape_c = (ctypes.c_int * 3)(*bin_data.shape)
        self.ptr = go_octree.NewOcTree(
            bin_data.flatten(),
            shape_c[0],
            shape_c[1],
            shape_c[2],
            root_offset[0],
            root_offset[1],
            root_offset[2],
        )  # uintptr in Go

    def __del__(self):
        go_octree.DeleteOcTree(self.ptr)

    def find_min_dist(
        self,
        coords: tuple[int, int, int],
        scaling: Optional[tuple[float, float, float]] = None,
    ) -> tuple[float, tuple[int, int, int]]:
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
        scaling_c = (ctypes.c_double * 3)(*scaling)
        go_octree.FindMinDist(
            self.ptr,
            coords_c[0],
            coords_c[1],
            coords_c[2],
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
