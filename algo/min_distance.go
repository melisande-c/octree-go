package algo

import (
	"math"
	"slices"
	"sort"

	"github.com/melisande-c/octree-go/data"
)

type minFinder struct {
	tree          data.OcTree
	queue         []*data.OcNode
	queueDists    []float64
	currentMin    float64
	currentMinLoc [3]int
	scaling       [3]float64
}

func FindMinLoc(
	tree data.OcTree, coords [3]int, scaling [3]float64,
) (float64, [3]int) {
	finder := minFinder{
		tree: tree, scaling: scaling,
	}
	finder.initMin(coords)
	finder.search(coords)
	return finder.currentMin, finder.currentMinLoc
}

func (f *minFinder) initMin(coords [3]int) {
	bounds := f.tree.Root.Bounds()
	var scaled [2][3]float64
	for i, b := range bounds {
		for j, x := range b {
			scaled[i][j] = float64(x) * f.scaling[j]
		}
	}
	dists := make([]float64, 0, 6)
	for _, s := range scaled {
		for j, x := range s {
			dists = append(dists, math.Abs(float64(coords[j])-x))
		}
	}
	f.currentMin = slices.Max(dists)
}

func (f *minFinder) search(coords [3]int) {
	f.traverse(&f.tree.Root, coords)
	for len(f.queue) > 0 {
		f.filterQueue()
		f.sortQueue()
		// copy queue then clear
		queueSlice := f.queue[:]
		f.queue = make([]*data.OcNode, 0)
		f.queueDists = make([]float64, 0)
		for _, node := range queueSlice {
			f.traverse(node, coords)
		}
	}
}

func (f *minFinder) traverse(node *data.OcNode, coords [3]int) {
	dist, loc := distToCube(coords, node.Bounds(), f.scaling)
	// quick return if dist is greater than current min
	if dist >= f.currentMin {
		return
	}

	isInBounds := isInBounds(coords, node.XBounds, node.YBounds, node.ZBounds)
	if node.IsLeaf && node.ContainsData && isInBounds {
		f.currentMin = 0
		f.currentMinLoc = coords
		return
	} else if node.IsLeaf && node.ContainsData && !isInBounds {
		if dist < f.currentMin {
			f.currentMin, f.currentMinLoc = dist, loc
		}
		return
	} else if node.IsLeaf && !node.ContainsData {
		panic("At leaf node with no data")
	}
	// continue if node is not a leaf

	cubeDists := make([]float64, 0, 8)
	nodes := make([]*data.OcNode, 0, 8)
	for _, n := range node.Children {
		dist, _ := distToCube(coords, n.Bounds(), f.scaling)
		if n.ContainsData {
			cubeDists = append(cubeDists, dist)
			nodes = append(nodes, n)
		}
	}
	idx := argMin(cubeDists[:])
	for i, n := range nodes {
		if i != idx {
			f.queue = append(f.queue, n)
			f.queueDists = append(f.queueDists, cubeDists[i])
		}
	}
	f.traverse(nodes[idx], coords)
}

func (f *minFinder) filterQueue() {
	filteredQueue := make([]*data.OcNode, 0, len(f.queue))
	filteredQueueDists := make([]float64, 0, len(f.queueDists))
	for i, d := range f.queueDists {
		if d <= f.currentMin {
			filteredQueueDists = append(filteredQueueDists, d)
			filteredQueue = append(filteredQueue, f.queue[i])
		}
	}
	f.queue = filteredQueue
	f.queueDists = filteredQueueDists
}

func (f *minFinder) sortQueue() {
	sort.Slice(f.queueDists, func(i, j int) bool {
		// Swap elements in arr2 to maintain correspondence
		if f.queueDists[i] < f.queueDists[j] {
			f.queueDists[i], f.queueDists[j] = f.queueDists[j], f.queueDists[i]
			f.queue[i], f.queue[j] = f.queue[j], f.queue[i] // Apply same swaps to arr2
		}
		return f.queueDists[i] < f.queueDists[j]
	})
}

func isInBounds(coords [3]int, xBounds [2]int, yBounds [2]int, zBounds [2]int) bool {
	bounds := [3][2]int{xBounds, yBounds, zBounds}
	in := true
	for i, c := range coords {
		in = in && ((bounds[i][0] <= c) && (c < bounds[i][1]))
	}
	return in
}

func argMin(slice []float64) int {
	currentMin := slice[0]
	var currentMinIdx int
	currentMinIdx = 0
	for i, v := range slice {
		if v < currentMin {
			currentMin = v
			currentMinIdx = i
		}
	}
	return currentMinIdx
}

func distToCube(coords [3]int, bounds [2][3]int, scaling [3]float64) (float64, [3]int) {
	var closestPoint [3]int
	for i, c := range coords {
		if (bounds[0][i] <= c) && (c < bounds[1][i]) {
			closestPoint[i] = c
		} else if c < bounds[0][i] {
			closestPoint[i] = bounds[0][i]
		} else { // bounds[i][1] <= c
			closestPoint[i] = bounds[1][i] - 1
		}
	}
	var sCoords [3]float64
	var sClosestPoint [3]float64
	for i := range 3 {
		sCoords[i] = float64(coords[i]) * scaling[i]
		sClosestPoint[i] = float64(closestPoint[i]) * scaling[i]
	}

	dist := math.Sqrt(
		math.Pow(sCoords[0]-sClosestPoint[0], 2) +
			math.Pow(sCoords[1]-sClosestPoint[1], 2) +
			math.Pow(sCoords[2]-sClosestPoint[2], 2),
	)
	return dist, closestPoint
}
