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

type Mesh struct {
	triangles []*Triangle
	bv        BoundingVolume
}

func (m Mesh) GetColor() Vec3 {
	return Vec3{0.1, 0.7, 0.9}
}

func (m Mesh) IntersectHit(r Ray) (bool, Hit) {
	var pHit Hit
	var minDistance = infinity
	if isHit := m.bv.Intersect(r); !isHit {
		return false, zeroHit
	}
	for _, triangle := range m.triangles {
		if isHit, hit := triangle.IntersectHit(r); isHit {
			if hit.distance < minDistance {
				pHit = hit
				minDistance = pHit.distance
			}
		}
	}
	if minDistance == infinity {
		return false, zeroHit
	}
	return true, pHit
}
