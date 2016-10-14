package main

/*
   Copyright (C) 2016 Nathan Jaremko

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

import (
	"math"
	"sort"
)

// Axis represents which axis we partition on
type Axis uint8

const (
	// AxisNone is used when we don't have an axis to partition on
	AxisNone Axis = iota
	// AxisX is used when we partition on X
	AxisX Axis = iota
	// AxisY is used when we partition on Y
	AxisY Axis = iota
	// AxisZ is used when we partition on Z
	AxisZ Axis = iota
)

// KdTree represents the root of a kd-Tree
type KdTree struct {
	Box  *Box
	Root *Node
}

func buildTree(triangles []*Triangle) *KdTree {
	// Ensure our bounding box contains all triangles
	box := triangles[0].boundingBox()
	for _, triangle := range triangles[1:] {
		box.Expand(triangle.boundingBox())
	}
	node := newNode(triangles)
	node.split(0)
	return &KdTree{box, node}
}

// Intersect performs an intersection test on the kd-tree
func (tree *KdTree) Intersect(r Ray) Hit {
	tmin, tmax := tree.Box.Intersect(r)
	if tmax < tmin || tmax <= 0 {
		return NoHit
	}
	return tree.Root.intersect(r, tmin, tmax)
}

// Node represents a node in a kd-tree
type Node struct {
	Axis      Axis
	Point     float64
	Triangles []*Triangle
	Left      *Node
	Right     *Node
}

func newNode(shapes []*Triangle) *Node {
	return &Node{AxisNone, 0, shapes, nil, nil}
}

func (node *Node) intersect(r Ray, tmin, tmax float64) Hit {
	var tsplit float64
	var leftFirst bool
	switch node.Axis {
	case AxisNone:
		return node.intersectTriangles(r)
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
		return first.intersect(r, tmin, tmax)
	} else if tsplit < tmin {
		return second.intersect(r, tmin, tmax)
	} else {
		h1 := first.intersect(r, tmin, tsplit)
		if h1.T <= tsplit {
			return h1
		}
		h2 := second.intersect(r, tsplit, math.Min(tmax, h1.T))
		if h1.T <= h2.T {
			return h1
		}
		return h2
	}
}

func (node *Node) intersectTriangles(r Ray) Hit {
	hit := NoHit
	for _, triangle := range node.Triangles {
		_, h := triangle.IntersectHit(r)
		if h.T < hit.T {
			hit = h
		}
	}
	return hit
}

func (node *Node) partitionScore(axis Axis, point float64) int {
	left, right := 0, 0
	for _, shape := range node.Triangles {
		box := shape.boundingBox()
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
	}
	return right
}

func (node *Node) partition(size int, axis Axis, point float64) (left, right []*Triangle) {
	left = make([]*Triangle, 0, size)
	right = make([]*Triangle, 0, size)
	for _, shape := range node.Triangles {
		box := shape.boundingBox()
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

func (node *Node) split(depth int) {
	if len(node.Triangles) < 8 {
		return
	}
	xs := make([]float64, 0, len(node.Triangles)*2)
	ys := make([]float64, 0, len(node.Triangles)*2)
	zs := make([]float64, 0, len(node.Triangles)*2)
	for _, shape := range node.Triangles {
		box := shape.boundingBox()
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
	mx, my, mz := median(xs), median(ys), median(zs)
	best := int(float64(len(node.Triangles)) * 0.85)
	bestAxis := AxisNone
	bestPoint := 0.0
	sx := node.partitionScore(AxisX, mx)
	if sx < best {
		best = sx
		bestAxis = AxisX
		bestPoint = mx
	}
	sy := node.partitionScore(AxisY, my)
	if sy < best {
		best = sy
		bestAxis = AxisY
		bestPoint = my
	}
	sz := node.partitionScore(AxisZ, mz)
	if sz < best {
		best = sz
		bestAxis = AxisZ
		bestPoint = mz
	}
	if bestAxis == AxisNone {
		return
	}
	l, r := node.partition(best, bestAxis, bestPoint)
	node.Axis = bestAxis
	node.Point = bestPoint
	node.Left = newNode(l)
	node.Right = newNode(r)
	node.Left.split(depth + 1)
	node.Right.split(depth + 1)
	node.Triangles = nil
}
