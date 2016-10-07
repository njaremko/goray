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
		if vert.X < min.X {
			min.X = vert.X
		}
		if vert.Y < min.Y {
			min.Y = vert.Y
		}
		if vert.Z < min.Z {
			min.Z = vert.Z
		}
		if vert.X > max.X {
			max.X = vert.X
		}
		if vert.Y > max.Y {
			max.Y = vert.Y
		}
		if vert.Z > max.Z {
			max.Z = vert.Z
		}
	}
	return &Box{min, max}, nil
}
