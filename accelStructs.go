package main

import (
	"errors"
	"math"
)

type BVH struct {
	tree []BoundingVolume
}

func computeBoundingSphere(vertSlice []Vec3) *Sphere {
	var avgPos Vec3
	maxSqDist := 0.0
	inverseNumVerts := 1.0 / float64(len(vertSlice))
	for _, vert := range vertSlice {
		avgPos.Add(vert.Mul(inverseNumVerts))
	}

	for _, vert := range vertSlice {
		diff := avgPos.Sub(vert)
		sqDist := dotProduct(diff, diff)
		if sqDist > maxSqDist {
			maxSqDist = sqDist
		}
	}
	return &Sphere{center: avgPos, radius: math.Sqrt(maxSqDist)}
}

func computeBoundingBox(vertSlice []Vec3) (*Box, error) {
	if len(vertSlice) < 2 {
		return nil, errors.New("vertSlice is too small to compute bounding box.")
	}
	min, max := vertSlice[0], vertSlice[0]
	for _, vert := range vertSlice[1:] {
		if vert[0] < min[0] {
			min[0] = vert[0]
		}
		if vert[1] < min[1] {
			min[1] = vert[1]
		}
		if vert[2] < min[2] {
			min[2] = vert[2]
		}
		if vert[0] > max[0] {
			max[0] = vert[0]
		}
		if vert[1] > max[1] {
			max[1] = vert[1]
		}
		if vert[2] > max[2] {
			max[2] = vert[2]
		}
	}
	return &Box{min, max}, nil
}
