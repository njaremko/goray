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

const EPSILON = 0.0001

type Triangle struct {
	v1           Vec3
	v2           Vec3
	v3           Vec3
	normal       Vec3
	color        Vec3
	transparency float64
	reflection   float64
	center       Vec3
}

func (s *Triangle) GetColor() Vec3 {
	return s.color
}

func (s *Triangle) IsTransparent() bool {
	return s.transparency > 0
}

func (s *Triangle) GetTransparency() float64 {
	return s.transparency
}

func (s *Triangle) GetReflection() float64 {
	return s.reflection
}

func (s *Triangle) GetCenter() Vec3 {
	return s.center
}

// Moller-Trumbore algorithm
func (s *Triangle) Intersect(r Ray) (bool, Hit) {
	//Find vectors for two edges sharing V1
	e1 := s.v2.Sub(s.v1)
	e2 := s.v3.Sub(s.v1)
	//Begin calculating determinant - also used to calculate u parameter
	p := crossProduct(r.dir, e2)
	//if determinant is near zero, ray lies in plane of triangle or ray is parallel to plane of triangle
	det := dotProduct(e1, p)
	//NOT CULLING
	if det > -EPSILON && det < EPSILON {
		return false, Hit{0, zeroVec, zeroVec}
	}
	inv_det := 1 / det
	//calculate distance from V1 to ray origin
	t := r.origin.Sub(s.v1)
	//Calculate u parameter and test bound
	u := dotProduct(t, p) * inv_det
	//The intersection lies outside of the triangle
	if u < 0 || u > 1 {
		return false, Hit{0, zeroVec, zeroVec}
	}
	//Prepare to test v parameter
	q := crossProduct(t, e1)
	//Calculate V parameter and test bound
	v := dotProduct(r.dir, q) * inv_det
	//The intersection lies outside of the triangle
	if v < 0 || u+v > 1 {
		return false, Hit{0, zeroVec, zeroVec}
	}
	x := dotProduct(e2, q) * inv_det
	if x > EPSILON { //ray intersection
		return true, Hit{x, r.origin.Add(r.dir.Mul(x)), crossProduct(e1, e2)}
	}
	return false, Hit{0, zeroVec, zeroVec}
}
