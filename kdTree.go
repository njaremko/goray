package main

import (
	"math"
	"sort"
)

type Axis uint8

const (
	AxisNone Axis = iota
	AxisX    Axis = iota
	AxisY    Axis = iota
	AxisZ    Axis = iota
)

type KdTree struct {
	Box  *Box
	Root *Node
}

func BuildTree(triangles []*Triangle) *KdTree {
	// Ensure our bounding box contains all triangles
	box := triangles[0].BoundingBox()
	for _, triangle := range triangles[1:] {
		box.Expand(triangle.BoundingBox())
	}
	node := NewNode(triangles)
	node.Split(0)
	return &KdTree{box, node}
}

func (tree *KdTree) Intersect(r Ray) Hit {
	tmin, tmax := tree.Box.Intersect(r)
	if tmax < tmin || tmax <= 0 {
		return NoHit
	}
	return tree.Root.Intersect(r, tmin, tmax)
}

type Node struct {
	Axis      Axis
	Point     float64
	Triangles []*Triangle
	Left      *Node
	Right     *Node
}

func NewNode(shapes []*Triangle) *Node {
	return &Node{AxisNone, 0, shapes, nil, nil}
}

func (node *Node) Intersect(r Ray, tmin, tmax float64) Hit {
	var tsplit float64
	var leftFirst bool
	switch node.Axis {
	case AxisNone:
		return node.IntersectTriangles(r)
	case AxisX:
		tsplit = (node.Point - r.Origin.X) / r.Direction.X
		leftFirst = (r.Origin.X < node.Point) || (r.Origin.X == node.Point && r.Direction.X <= 0)
	case AxisY:
		tsplit = (node.Point - r.Origin.Y) / r.Direction.Y
		leftFirst = (r.Origin.Y < node.Point) || (r.Origin.Y == node.Point && r.Direction.Y <= 0)
	case AxisZ:
		tsplit = (node.Point - r.Origin.Z) / r.Direction.Z
		leftFirst = (r.Origin.Z < node.Point) || (r.Origin.Z == node.Point && r.Direction.Z <= 0)
	}
	var first, second *Node
	if leftFirst {
		first = node.Left
		second = node.Right
	} else {
		first = node.Right
		second = node.Left
	}
	if tsplit > tmax || tsplit <= 0 {
		return first.Intersect(r, tmin, tmax)
	} else if tsplit < tmin {
		return second.Intersect(r, tmin, tmax)
	} else {
		h1 := first.Intersect(r, tmin, tsplit)
		if h1.T <= tsplit {
			return h1
		}
		h2 := second.Intersect(r, tsplit, math.Min(tmax, h1.T))
		if h1.T <= h2.T {
			return h1
		} else {
			return h2
		}
	}
}

func (node *Node) IntersectTriangles(r Ray) Hit {
	hit := NoHit
	for _, shape := range node.Triangles {
		_, h := shape.IntersectHit(r)
		if h.T < hit.T {
			hit = h
		}
	}
	return hit
}

func (node *Node) PartitionScore(axis Axis, point float64) int {
	left, right := 0, 0
	for _, shape := range node.Triangles {
		box := shape.BoundingBox()
		l, r := box.Partition(axis, point)
		if l {
			left++
		}
		if r {
			right++
		}
	}
	if left >= right {
		return left
	} else {
		return right
	}
}

func (node *Node) Partition(size int, axis Axis, point float64) (left, right []*Triangle) {
	left = make([]*Triangle, 0, size)
	right = make([]*Triangle, 0, size)
	for _, shape := range node.Triangles {
		box := shape.BoundingBox()
		l, r := box.Partition(axis, point)
		if l {
			left = append(left, shape)
		}
		if r {
			right = append(right, shape)
		}
	}
	return
}

func (node *Node) Split(depth int) {
	if len(node.Triangles) < 8 {
		return
	}
	xs := make([]float64, 0, len(node.Triangles)*2)
	ys := make([]float64, 0, len(node.Triangles)*2)
	zs := make([]float64, 0, len(node.Triangles)*2)
	for _, shape := range node.Triangles {
		box := shape.BoundingBox()
		xs = append(xs, box.min.X)
		xs = append(xs, box.max.X)
		ys = append(ys, box.min.Y)
		ys = append(ys, box.max.Y)
		zs = append(zs, box.min.Z)
		zs = append(zs, box.max.Z)
	}
	sort.Float64s(xs)
	sort.Float64s(ys)
	sort.Float64s(zs)
	mx, my, mz := Median(xs), Median(ys), Median(zs)
	best := int(float64(len(node.Triangles)) * 0.85)
	bestAxis := AxisNone
	bestPoint := 0.0
	sx := node.PartitionScore(AxisX, mx)
	if sx < best {
		best = sx
		bestAxis = AxisX
		bestPoint = mx
	}
	sy := node.PartitionScore(AxisY, my)
	if sy < best {
		best = sy
		bestAxis = AxisY
		bestPoint = my
	}
	sz := node.PartitionScore(AxisZ, mz)
	if sz < best {
		best = sz
		bestAxis = AxisZ
		bestPoint = mz
	}
	if bestAxis == AxisNone {
		return
	}
	l, r := node.Partition(best, bestAxis, bestPoint)
	node.Axis = bestAxis
	node.Point = bestPoint
	node.Left = NewNode(l)
	node.Right = NewNode(r)
	node.Left.Split(depth + 1)
	node.Right.Split(depth + 1)
	node.Triangles = nil
}
