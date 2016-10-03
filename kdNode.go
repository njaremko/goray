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

type KdNode struct {
	bbox      Box
	left      *KdNode
	right     *KdNode
	triangles []*Triangle
}

func BuildKdNode(triangles []*Triangle, depth int) *KdNode {
	// Initialize a new node
	node := &KdNode{triangles: triangles}

	// Handle base cases
	if len(triangles) == 0 {
		return node
	}
	if len(triangles) == 1 {
		node.bbox = *triangles[0].BoundingBox()
		node.left = &KdNode{}
		node.right = &KdNode{}
		return node
	}

	// Ensure our bounding box contains all triangles
	node.bbox = *triangles[0].BoundingBox()
	for _, triangle := range node.triangles[1:] {
		node.bbox.Expand(triangle.BoundingBox())
	}

	// Calculate the middle point of all triangles
	midPoint := Vec3{0, 0, 0}
	inverseLen := 1.0 / float64(len(triangles))
	for _, triangle := range triangles {
		midPoint = midPoint.Add(triangle.MidPoint().Mul(inverseLen))
	}

	// Initialize slices to be filled in next section
	leftTriangles := []*Triangle{}
	rightTriangles := []*Triangle{}

	// This is where we break along optimal axis
	axis := node.bbox.LongestAxis()
	for _, triangle := range triangles {
		switch axis {
		case 0:
			if midPoint.x >= triangle.MidPoint().x {
				rightTriangles = append(rightTriangles, triangle)
			} else {
				leftTriangles = append(leftTriangles, triangle)
			}
		case 1:
			if midPoint.y >= triangle.MidPoint().y {
				rightTriangles = append(rightTriangles, triangle)
			} else {
				leftTriangles = append(leftTriangles, triangle)
			}
		case 2:
			if midPoint.z >= triangle.MidPoint().z {
				rightTriangles = append(rightTriangles, triangle)
			} else {
				leftTriangles = append(leftTriangles, triangle)
			}
		}
	}

	// If one of our slices is empty, just set it to the other one
	if len(leftTriangles) == 0 && len(rightTriangles) > 0 {
		leftTriangles = rightTriangles
	}
	if len(rightTriangles) == 0 && len(leftTriangles) > 0 {
		rightTriangles = leftTriangles
	}

	// Count how many triangles left and right tree share
	matches := 0
	for _, left := range leftTriangles {
		for _, right := range rightTriangles {
			if left.Equals(right) {
				matches++
			}
		}
	}

	// If they share less than half of the triangles in parent, go deeper
	if float64(matches)/float64(len(leftTriangles)) < 0.5 && float64(matches)/float64(len(rightTriangles)) < 0.5 {
		node.left = BuildKdNode(leftTriangles, depth+1)
		node.right = BuildKdNode(rightTriangles, depth+1)
	} else {
		// Both child trees share >=50% of parents triangles, so we can stop.
		node.left = &KdNode{}
		node.right = &KdNode{}
	}

	return node
}

func (node *KdNode) Intersect(r Ray) (bool, Hit) {
	// If the ray doesn't Hit the KdTree, then it's a miss.
	if !node.bbox.Intersect(r) {
		return false, noHit
	}
	if len(node.left.triangles) > 0 || len(node.right.triangles) > 0 {
		// This is an internal node
		// Check the left tree for a hit
		isHit, hit := node.left.Intersect(r)
		if isHit {
			return true, hit
		}
		// Check the right tree for a hit
		isHit, hit = node.right.Intersect(r)
		if isHit {
			return true, hit
		}
		return false, noHit
	} else {
		// This is a leaf
		for _, triangle := range node.triangles {
			isHit, hit := triangle.IntersectHit(r)
			if isHit {
				return true, hit
			}
		}
		return false, noHit
	}
}
