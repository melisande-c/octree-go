import os
import time
from typing import Optional
import ctypes

import numpy as np
from numpy.typing import NDArray
import numpy.ctypeslib as npct
import matplotlib.pyplot as plt
import tifffile

from pyoctree.octree import Tree

print(os.getcwd())

img = tifffile.imread(
    "/Users/melisande.croft/Documents/Data/5639253/Multilabel_U-Net_dataset_B.subtilis/training/instance_segmentation_GT/train_18.tif"
)
array = np.zeros((*img.shape, 256), dtype=np.int32)
array[img != 0, :] = 1

rng = np.random.default_rng(seed=42)
n_points = 256
start = np.zeros(3, dtype=int)
start[-1] = array.shape[2] // 2
end = np.array(array.shape)
end[-1] = (array.shape[2] // 2) + 1
points = rng.integers(start, end, (n_points, 3))

t0 = time.time()
tree = Tree(array)
print(f"Time taken to build Octree {time.time()-t0:.2f}s")
locs: list[tuple[int, int, int]] = []
t_points = time.time()
for point in points:
    query = tuple(point)
    # query = (0, 0, array.shape[2]//2)
    t0 = time.time()
    dist, loc = tree.find_min_dist(query, scaling=(0.5,0.5,0.5))
    # print(f"Time taken to query point {time.time()-t0:.2f}s")
    locs.append(loc)
    print(dist, loc)
print(f"Time taken to all points {time.time()-t_points:.2f}s")


# locs: list[tuple[int, int, int]] = []
# for point in points:
#     query = tuple(point)
#     # query = (0, 0, array.shape[2]//2)
#     dist, loc = find_min_dist(query, array)
#     locs.append(loc)
#     print(dist, loc)

fig, ax = plt.subplots()
# z = 0
z = array.shape[2] // 2
# z = 8
ax.imshow(array[:, :, z], "gray")

for query, loc in zip(points, locs):
    ax.plot(query[1], query[0], "x", c="cyan")
    ax.plot(loc[1], loc[0], "x", c="magenta")
    ax.plot([query[1], loc[1]], [query[0], loc[0]], "--", c="yellow", linewidth=1.5)

plt.show()
