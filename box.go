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
	inverseRay := r.dir.Inverse()
	tx1 := (b.min.x - r.origin.x) * inverseRay.x
	tx2 := (b.max.x - r.origin.x) * inverseRay.x

	tmin := math.Min(tx1, tx2)
	tmax := math.Max(tx1, tx2)

	ty1 := (b.min.y - r.origin.y) * inverseRay.y
	ty2 := (b.max.y - r.origin.y) * inverseRay.y

	tmin = math.Max(tmin, math.Min(ty1, ty2))
	tmax = math.Min(tmax, math.Max(ty1, ty2))

	tz1 := (b.min.z - r.origin.z) * inverseRay.z
	tz2 := (b.max.z - r.origin.z) * inverseRay.z

	tmin = math.Max(tmin, math.Min(tz1, tz2))
	tmax = math.Min(tmax, math.Max(tz1, tz2))

	return tmax >= tmin
}
