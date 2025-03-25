package algo

import (
	"math"
	"sort"

	"github.com/melisande-c/octree-go/data"
)

type minFinder struct {
	tree          data.OcTree
	queue         []*data.OcNode
	queueDists    []float64
	currentMin    float64
	currentMinLoc [3]int
}

func FindMinLoc(tree data.OcTree, coords [3]int) (float64, [3]int) {
	current_min := max(tree.Root.XBounds[1], tree.Root.YBounds[1], tree.Root.ZBounds[1])
	finder := minFinder{tree: tree, currentMin: float64(current_min)}
	finder.search(coords)
	return finder.currentMin, finder.currentMinLoc
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
	dist, _ := distToCube(coords, node.XBounds, node.YBounds, node.ZBounds)
	if dist >= f.currentMin {
		return
	}

	isInBounds := isInBounds(coords, node.XBounds, node.YBounds, node.ZBounds)
	if node.IsLeaf && node.ContainsData && isInBounds {
		f.currentMin = 0
		f.currentMinLoc = coords
		return
	} else if node.IsLeaf && node.ContainsData && !isInBounds {
		newMin, newLoc := distToCube(coords, node.XBounds, node.YBounds, node.ZBounds)
		if newMin < f.currentMin {
			f.currentMin, f.currentMinLoc = newMin, newLoc
		}
		return
	} else if node.IsLeaf && !node.ContainsData {
		panic("At leaf node with no data")
	}

	cubeDists := make([]float64, 0, 8)
	nodes := make([]*data.OcNode, 0, 8)
	for _, n := range node.Children {
		dist, _ := distToCube(coords, n.XBounds, n.YBounds, n.ZBounds)
		if n.ContainsData {
			// if n.ContainsData && (dist <= f.current_min) {
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
		if d < f.currentMin {
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

func distToCube(
	coords [3]int, xBounds [2]int, yBounds [2]int, zBounds [2]int,
) (float64, [3]int) {
	bounds := [3][2]int{xBounds, yBounds, zBounds}
	var closestPoint [3]int
	for i, c := range coords {
		if (bounds[i][0] <= c) && (c < bounds[i][1]) {
			closestPoint[i] = c
		} else if c < bounds[i][0] {
			closestPoint[i] = bounds[i][0]
		} else { // bounds[i][1] <= c
			closestPoint[i] = bounds[i][1] - 1
		}
	}
	dist := math.Sqrt(
		math.Pow(float64(coords[0])-float64(closestPoint[0]), 2) +
			math.Pow(float64(coords[1])-float64(closestPoint[1]), 2) +
			math.Pow(float64(coords[2])-float64(closestPoint[2]), 2),
	)
	return dist, closestPoint
}
