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

// Box represents a box...
type Box struct {
	min Vec3
	max Vec3
}

// Partition partitions a box on an axis and says if the min and max on that axis
// are less/greater than the point
func (b *Box) Partition(axis Axis, point float64) (less, greater bool) {
	switch axis {
	case AxisX:
		less = b.min.X <= point
		greater = b.max.X >= point
	case AxisY:
		less = b.min.Y <= point
		greater = b.max.Y >= point
	case AxisZ:
		less = b.min.Z <= point
		greater = b.max.Z >= point
	}
	return
}

// Intersect calculates the intersection of a Ray with the box
func (b *Box) Intersect(r Ray) (float64, float64) {
	inverseDir := r.Direction.Inverse()

	tx1 := (b.min.X - r.Origin.X) * inverseDir.X
	tx2 := (b.max.X - r.Origin.X) * inverseDir.X

	tmin := math.Min(tx1, tx2)
	tmax := math.Max(tx1, tx2)

	ty1 := (b.min.Y - r.Origin.Y) * inverseDir.Y
	ty2 := (b.max.Y - r.Origin.Y) * inverseDir.Y

	tmin = math.Max(tmin, math.Min(ty1, ty2))
	tmax = math.Min(tmax, math.Max(ty1, ty2))

	tz1 := (b.min.Z - r.Origin.Z) * inverseDir.Z
	tz2 := (b.max.Z - r.Origin.Z) * inverseDir.Z

	tmin = math.Max(tmin, math.Min(tz1, tz2))
	tmax = math.Min(tmax, math.Max(tz1, tz2))

	return math.Max(0.0, tmin), tmax
}

// Len calculates the distance between box min and max
func (b *Box) Len() Vec3 {
	return b.max.Sub(b.min)
}

// Expand expands b to contain other
func (b *Box) Expand(other *Box) {
	if other.min.X < b.min.X {
		b.min.X = other.min.X
	}
	if other.min.Y < b.min.Y {
		b.min.Y = other.min.Y
	}
	if other.min.Z < b.min.Z {
		b.min.Z = other.min.Z
	}
	if b.max.X < other.max.X {
		b.max.X = other.max.X
	}
	if b.max.Y < other.max.Y {
		b.max.Y = other.max.Y
	}
	if b.max.Z < other.max.Z {
		b.max.Z = other.max.Z
	}
}

// LongestAxis returns the largest axis difference between box min and max
func (b *Box) LongestAxis() int {
	xLength := math.Abs(b.max.X - b.min.X)
	yLength := math.Abs(b.max.Y - b.min.Y)
	zLength := math.Abs(b.max.Z - b.min.Z)
	if xLength > yLength && xLength > zLength {
		return 0
	} else if yLength > xLength && yLength > zLength {
		return 1
	} else {
		return 2
	}
}

// Overlaps returns whether a box overlaps other
func (b *Box) Overlaps(other *Box) bool {
	x := b.max.X >= other.min.X && b.min.X <= other.max.X
	y := b.max.Y >= other.min.Y && b.min.Y <= other.max.Y
	z := b.max.Z >= other.min.Z && b.min.Z <= other.max.Z

	return x && y && z
}

func computeBoundingBox(vertSlice []Vec3) (*Box, error) {
	if len(vertSlice) < 2 {
		return nil, errors.New("vertSlice is too small to compute bounding box")
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
