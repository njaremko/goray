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
)

type Box struct {
	min Vec3
	max Vec3
}

func (b *Box) Intersect(r Ray) bool {
	inverseDir := r.dir.Inverse()

	tx1 := (b.min.x - r.origin.x) * inverseDir.x
	tx2 := (b.max.x - r.origin.x) * inverseDir.x

	tmin := math.Min(tx1, tx2)
	tmax := math.Max(tx1, tx2)

	ty1 := (b.min.y - r.origin.y) * inverseDir.y
	ty2 := (b.max.y - r.origin.y) * inverseDir.y

	tmin = math.Max(tmin, math.Min(ty1, ty2))
	tmax = math.Min(tmax, math.Max(ty1, ty2))

	tz1 := (b.min.z - r.origin.z) * inverseDir.z
	tz2 := (b.max.z - r.origin.z) * inverseDir.z

	tmin = math.Max(tmin, math.Min(tz1, tz2))
	tmax = math.Min(tmax, math.Max(tz1, tz2))

	return tmax >= math.Max(0.0, tmin)
}

func (b *Box) Len() Vec3 {
	return b.max.Sub(b.min)
}

func (b *Box) Expand(other *Box) {
	if other.min.x < b.min.x {
		b.min.x = other.min.x
	}
	if other.min.y < b.min.y {
		b.min.y = other.min.y
	}
	if other.min.z < b.min.z {
		b.min.z = other.min.z
	}
	if b.max.x < other.max.x {
		b.max.x = other.max.x
	}
	if b.max.y < other.max.y {
		b.max.y = other.max.y
	}
	if b.max.z < other.max.z {
		b.max.z = other.max.z
	}
}

func (b *Box) LongestAxis() int {
	xLength := math.Abs(b.max.x - b.min.x)
	yLength := math.Abs(b.max.y - b.min.y)
	zLength := math.Abs(b.max.z - b.min.z)
	if xLength > yLength && xLength > zLength {
		return 0
	} else if yLength > xLength && yLength > zLength {
		return 1
	} else {
		return 2
	}
}

func (b *Box) Overlaps(other *Box) bool {
	x := b.max.x >= other.min.x && b.min.x <= other.max.x
	y := b.max.y >= other.min.y && b.min.y <= other.max.y
	z := b.max.z >= other.min.z && b.min.z <= other.max.z

	return x && y && z
}
