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
	v1           Vec3
	v2           Vec3
	v3           Vec3
	normal       Vec3
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

func (t *Triangle) ComputeBoundingBox() *Box {
	box, err := computeBoundingBox([]Vec3{t.v1, t.v2, t.v3})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return box
}

func (t *Triangle) GetMidPoint() Vec3 {
	// Add three points of triangle and average
	xAverage := (t.v1.x + t.v2.x + t.v3.x) / 3
	yAverage := (t.v1.y + t.v2.y + t.v3.y) / 3
	zAverage := (t.v1.z + t.v2.z + t.v3.z) / 3
	return Vec3{xAverage, yAverage, zAverage}
}

func (t *Triangle) Equals(t2 *Triangle) bool {
	if t.v1.x == t2.v1.x && t.v1.y == t2.v1.y && t.v1.z == t2.v1.z {
		if t.v2.x == t2.v2.x && t.v2.y == t2.v2.y && t.v2.z == t2.v2.z {
			if t.v3.x == t2.v3.x && t.v3.y == t2.v3.y && t.v3.z == t2.v3.z {
				return true
			}
		}
	}
	return false
}

// Moller-Trumbore algorithm
func (t *Triangle) IntersectHit(r Ray) (bool, Hit) {
	//Find vectors for two edges sharing V1
	e1 := t.v2.Sub(t.v1)
	e2 := t.v3.Sub(t.v1)
	//Begin calculating determinant - also used to calculate u parameter
	p := crossProduct(r.dir, e2)
	//if determinant is near zero, ray lies in plane of triangle or ray is parallel to plane of triangle
	det := dotProduct(e1, p)
	// CULLING
	if /*det > -EPSILON &&*/ det < EPSILON {
		return false, noHit
	}
	invDet := 1.0 / det
	//calculate distance from V1 to ray origin
	s := r.origin.Sub(t.v1)
	//Calculate u parameter and test bound
	u := dotProduct(s, p) * invDet
	//The intersection lies outside of the triangle
	if u < 0.0 || u > 1.0 {
		return false, noHit
	}
	//Prepare to test v parameter
	q := crossProduct(s, e1)
	//Calculate V parameter and test bound
	v := dotProduct(r.dir, q) * invDet
	//The intersection lies outside of the triangle
	if v < 0.0 || u+v > 1.0 {
		return false, noHit
	}
	x := dotProduct(e2, q) * invDet
	if x > EPSILON { //ray intersection
		hitPoint := r.origin.Add(r.dir.Mul(x))
		normal := hitPoint.Add(crossProduct(e1, e2).Mul(1e-8))
		return true, Hit{x, hitPoint, normal}
	}
	return false, noHit
}
