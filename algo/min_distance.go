package algo

import (
	"math"
	"reflect"
	"sort"

	"github.com/melisande-c/octree-go/data"
)

type linTransform struct {
	Offset  [3]int
	Scaling [3]float64
}

type minFinder struct {
	tree          data.OcTree
	queue         []*data.OcNode
	queueDists    []float64
	currentMin    float64
	currentMinLoc [3]int
	linTransform  linTransform
}

func FindMinLoc(
	tree data.OcTree, coords [3]int, offset [3]int, scaling [3]float64,
) (float64, [3]int) {
	linT := linTransform{Offset: offset, Scaling: scaling}
	bounds := tree.Root.Bounds()
	var tBounds [2][3]float64
	for i, b := range bounds {
		tBounds[i] = applyTransform(b, linT)
	}
	current_min := max(tBounds[1][0], tBounds[1][1], tBounds[1][2])
	finder := minFinder{
		tree: tree, currentMin: float64(current_min), linTransform: linT,
	}
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
	dist, loc := distToCube(coords, node.Bounds(), f.linTransform)
	isInBounds := isInBounds(coords, node.XBounds, node.YBounds, node.ZBounds)
	// quick return if dist is greater than current min
	if dist > f.currentMin {
		return
	}

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
		dist, _ := distToCube(coords, n.Bounds(), f.linTransform)
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

func applyTransform(location [3]int, linTransform linTransform) [3]float64 {
	var transformed [3]float64
	for i := range 3 {
		transformed[i] = float64(location[i]+linTransform.Offset[i]) * (linTransform.Scaling[i])

	}
	return transformed
}

func applyInverseTransform(tLocation [3]float64, linTransform linTransform) [3]int {
	var inv_transformed [3]int
	for i := range 3 {
		// TODO: find a way to prevent floating point errors
		v := (tLocation[i] / linTransform.Scaling[i]) - float64(linTransform.Offset[i])
		inv_transformed[i] = int(math.Round(v))
	}
	return inv_transformed
}

func distToCube(
	coords [3]int, bounds [2][3]int, linTransform linTransform,
) (float64, [3]int) {
	isIdentity := (reflect.DeepEqual(linTransform.Offset, [3]int{0, 0, 0}) &&
		reflect.DeepEqual(linTransform.Scaling, [3]float64{1, 1, 1}))
	if isIdentity {
		return distToCubeNoTransform(coords, bounds)
	} else {
		return distToCubeTransform(coords, bounds, linTransform)
	}
}

func distToCubeNoTransform(coords [3]int, bounds [2][3]int) (float64, [3]int) {
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
	dist := math.Sqrt(
		math.Pow(float64(coords[0])-float64(closestPoint[0]), 2) +
			math.Pow(float64(coords[1])-float64(closestPoint[1]), 2) +
			math.Pow(float64(coords[2])-float64(closestPoint[2]), 2),
	)
	return dist, closestPoint
}

func distToCubeTransform(
	coords [3]int, bounds [2][3]int, linTransform linTransform,
) (float64, [3]int) {
	// fmt.Printf("Applying transform %+v\n", linTransform)
	var tBounds [2][3]float64 // transformed bounds
	for i, b := range bounds {
		tBounds[i] = applyTransform(b, linTransform)
	}
	tCoords := applyTransform(coords, linTransform) // transformed coords

	var tClosestPoint [3]float64
	var closestPoint [3]int
	for i, c := range tCoords {
		if (tBounds[0][i] <= c) && (c < tBounds[1][i]) {
			tClosestPoint[i] = c
			closestPoint[i] = coords[i]
		} else if c < tBounds[0][i] {
			tClosestPoint[i] = tBounds[0][i]
			closestPoint[i] = bounds[0][i]
		} else { // bounds[i][1] <= c
			tClosestPoint[i] = tBounds[1][i] - linTransform.Scaling[i]
			closestPoint[i] = bounds[1][i] - 1
		}
	}
	dist := math.Sqrt(
		math.Pow(tCoords[0]-tClosestPoint[0], 2) +
			math.Pow(tCoords[1]-tClosestPoint[1], 2) +
			math.Pow(tCoords[2]-tClosestPoint[2], 2),
	)
	return dist, closestPoint
}
