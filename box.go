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

	tx1 := (b.min[0] - r.origin[0]) * inverseDir[0]
	tx2 := (b.max[0] - r.origin[0]) * inverseDir[0]

	tmin := math.Min(tx1, tx2)
	tmax := math.Max(tx1, tx2)

	ty1 := (b.min[1] - r.origin[1]) * inverseDir[1]
	ty2 := (b.max[1] - r.origin[1]) * inverseDir[1]

	tmin = math.Max(tmin, math.Min(ty1, ty2))
	tmax = math.Min(tmax, math.Max(ty1, ty2))

	tz1 := (b.min[2] - r.origin[2]) * inverseDir[2]
	tz2 := (b.max[2] - r.origin[2]) * inverseDir[2]

	tmin = math.Max(tmin, math.Min(tz1, tz2))
	tmax = math.Min(tmax, math.Max(tz1, tz2))

	return tmax >= math.Max(0.0, tmin)
}
