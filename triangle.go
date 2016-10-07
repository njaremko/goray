package main

import (
	"fmt"
	"os"
)

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

// EPSILON added to normal vector to prevent acne
const EPSILON = 0.00001

type Triangle struct {
	V1, V2, V3   Vec3
	N1, N2, N3   Vec3
	T1, T2, T3   Vec3
	color        Vec3
	transparency float64
	reflection   float64
	center       Vec3
}

func (t *Triangle) GetColor() Vec3 {
	return t.color
}

func (t *Triangle) IsTransparent() bool {
	return t.transparency > 0
}

func (t *Triangle) GetTransparency() float64 {
	return t.transparency
}

func (t *Triangle) GetReflection() float64 {
	return t.reflection
}

func (t *Triangle) GetCenter() Vec3 {
	return t.center
}

func (t *Triangle) BoundingBox() *Box {
	box, err := computeBoundingBox([]Vec3{t.V1, t.V2, t.V3})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return box
}

func (t *Triangle) MidPoint() Vec3 {
	// Add three points of triangle and average
	xAverage := (t.V1.X + t.V2.X + t.V3.X) / 3
	yAverage := (t.V1.Y + t.V2.Y + t.V3.Y) / 3
	zAverage := (t.V1.Z + t.V2.Z + t.V3.Z) / 3
	return Vec3{xAverage, yAverage, zAverage}
}

func (t *Triangle) Normal() Vec3 {
	e1 := t.V2.Sub(t.V1)
	e2 := t.V3.Sub(t.V1)
	return crossProduct(e1, e2).Normalize()
}

func (t *Triangle) FixNormals() {
	n := t.Normal()
	zero := Vec3{}
	if t.N1 == zero {
		t.N1 = n
	}
	if t.N2 == zero {
		t.N2 = n
	}
	if t.N3 == zero {
		t.N3 = n
	}
}
func (t *Triangle) Equals(t2 *Triangle) bool {
	if t.V1.Equals(t2.V1) && t.V2.Equals(t2.V2) && t.V3.Equals(t2.V3) {
		return true
	}
	return false
}

// Moller-Trumbore algorithm
func (t *Triangle) IntersectHit(r Ray) (bool, Hit) {
	//Find vectors for two edges sharing V1
	e1 := t.V2.Sub(t.V1)
	e2 := t.V3.Sub(t.V1)
	//Begin calculating determinant - also used to calculate u parameter
	p := crossProduct(r.Direction, e2)
	//if determinant is near zero, ray lies in plane of triangle or ray is parallel to plane of triangle
	det := dotProduct(e1, p)
	// CULLING
	if /*det > -EPSILON &&*/ det < EPSILON {
		return false, NoHit
	}
	invDet := 1.0 / det
	//calculate distance from V1 to ray origin
	s := r.Origin.Sub(t.V1)
	//Calculate u parameter and test bound
	u := dotProduct(s, p) * invDet
	//The intersection lies outside of the triangle
	if u < 0.0 || u > 1.0 {
		return false, NoHit
	}
	//Prepare to test v parameter
	q := crossProduct(s, e1)
	//Calculate V parameter and test bound
	v := dotProduct(r.Direction, q) * invDet
	//The intersection lies outside of the triangle
	if v < 0.0 || u+v > 1.0 {
		return false, NoHit
	}
	x := dotProduct(e2, q) * invDet
	if x > EPSILON { //ray intersection
		hitPoint := r.Origin.Add(r.Direction.Mul(x))
		normal := crossProduct(e1, e2)
		return true, Hit{x, hitPoint, normal}
	}
	return false, NoHit
}
