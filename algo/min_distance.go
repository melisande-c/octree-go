package algo

import (
	"math"
	"sort"

	"github.com/melisande-c/octree-go/data_structure"
)

type minFinder struct {
	tree            data_structure.Tree
	search_queue    []*data_structure.Node
	queue_dists     []float64
	current_min     float64
	current_min_loc [3]int
}

func FindMinLoc(tree data_structure.Tree, coords [3]int) (float64, [3]int) {
	current_min := max(tree.Root.XBounds[1], tree.Root.YBounds[1], tree.Root.ZBounds[1])
	finder := minFinder{tree: tree, current_min: float64(current_min)}
	finder.search(coords)
	return finder.current_min, finder.current_min_loc
}

func (f *minFinder) search(coords [3]int) {
	f.traverse(&f.tree.Root, coords)
	for len(f.search_queue) > 0 {
		f.filterQueue()
		f.sortQueue()
		search_queue_slice := f.search_queue[:]
		f.search_queue = make([]*data_structure.Node, 0)
		f.queue_dists = make([]float64, 0)
		for _, node := range search_queue_slice {
			f.traverse(node, coords)
		}
	}
}

func (f *minFinder) traverse(node *data_structure.Node, coords [3]int) {
	dist, _ := distToCube(coords, node.XBounds, node.YBounds, node.ZBounds)
	if dist >= f.current_min {
		return
	}

	is_in_bounds := isInBounds(coords, node.XBounds, node.YBounds, node.ZBounds)
	if node.IsLeaf && node.ContainsData && is_in_bounds {
		f.current_min = 0
		f.current_min_loc = coords
		return
	} else if node.IsLeaf && node.ContainsData && !is_in_bounds {
		new_min, new_loc := distToCube(coords, node.XBounds, node.YBounds, node.ZBounds)
		if new_min < f.current_min {
			f.current_min, f.current_min_loc = new_min, new_loc
		}
		return
	} else if node.IsLeaf && !node.ContainsData {
		panic("At leaf node with no data")
	}

	cube_dist := make([]float64, 0, 8)
	nodes := make([]*data_structure.Node, 0, 8)
	for _, n := range node.Children {
		dist, _ := distToCube(coords, n.XBounds, n.YBounds, n.ZBounds)
		if n.ContainsData {
			// if n.ContainsData && (dist <= f.current_min) {
			cube_dist = append(cube_dist, dist)
			nodes = append(nodes, n)
		}
	}
	idx := argMin(cube_dist[:])
	for i, n := range nodes {
		if i != idx {
			f.search_queue = append(f.search_queue, n)
			f.queue_dists = append(f.queue_dists, cube_dist[i])
		}
	}
	f.traverse(nodes[idx], coords)
}
func (f *minFinder) filterQueue() {
	filtered_search_queue := make([]*data_structure.Node, 0, len(f.search_queue))
	filtered_queue_dists := make([]float64, 0, len(f.queue_dists))
	for i, d := range f.queue_dists {
		if d < f.current_min {
			filtered_queue_dists = append(filtered_queue_dists, d)
			filtered_search_queue = append(filtered_search_queue, f.search_queue[i])
		}
	}
	f.search_queue = filtered_search_queue
	f.queue_dists = filtered_queue_dists
}

func (f *minFinder) sortQueue() {
	sort.Slice(f.queue_dists, func(i, j int) bool {
		// Swap elements in arr2 to maintain correspondence
		if f.queue_dists[i] < f.queue_dists[j] {
			f.queue_dists[i], f.queue_dists[j] = f.queue_dists[j], f.queue_dists[i]
			f.search_queue[i], f.search_queue[j] = f.search_queue[j], f.search_queue[i] // Apply same swaps to arr2
		}
		return f.queue_dists[i] < f.queue_dists[j]
	})
}

func isInBounds(coords [3]int, x_bounds [2]int, y_bounds [2]int, z_bounds [2]int) bool {
	bounds := [3][2]int{x_bounds, y_bounds, z_bounds}
	in_bounds := true
	for i, c := range coords {
		in_bounds = in_bounds && ((bounds[i][0] <= c) && (c < bounds[i][1]))
	}
	return in_bounds
}

func argMin(slice []float64) int {
	current_min := slice[0]
	var current_min_idx int
	current_min_idx = 0
	for i, v := range slice {
		if v < current_min {
			current_min = v
			current_min_idx = i
		}
	}
	return current_min_idx
}

func distToCube(
	coords [3]int, x_bounds [2]int, y_bounds [2]int, z_bounds [2]int,
) (float64, [3]int) {
	bounds := [3][2]int{x_bounds, y_bounds, z_bounds}
	var closest_point [3]int
	for i, c := range coords {
		if (bounds[i][0] <= c) && (c < bounds[i][1]) {
			closest_point[i] = c
		} else if c < bounds[i][0] {
			closest_point[i] = bounds[i][0]
		} else { // bounds[i][1] <= c
			closest_point[i] = bounds[i][1] - 1
		}
	}
	distance := math.Sqrt(
		math.Pow(float64(coords[0])-float64(closest_point[0]), 2) +
			math.Pow(float64(coords[1])-float64(closest_point[1]), 2) +
			math.Pow(float64(coords[2])-float64(closest_point[2]), 2),
	)
	return distance, closest_point
}
