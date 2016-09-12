package main

import (
	"errors"
	"math"
)

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
		if vert.x < min.x {
			min.x = vert.x
		}
		if vert.y < min.y {
			min.y = vert.y
		}
		if vert.z < min.z {
			min.z = vert.z
		}
		if vert.x > max.x {
			max.x = vert.x
		}
		if vert.y > max.y {
			max.y = vert.y
		}
		if vert.z > max.z {
			max.z = vert.z
		}
	}
	return &Box{min, max}, nil
}
