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

type Sphere struct {
	center       Vec3
	radius       float64
	color        Vec3
	transparency float64
	reflection   float64
}

func (s *Sphere) GetColor() Vec3 {
	return s.color
}

func (s *Sphere) IsTransparent() bool {
	return s.transparency > 0
}

func (s *Sphere) GetTransparency() float64 {
	return s.transparency
}

func (s *Sphere) GetReflection() float64 {
	return s.reflection
}

func (s *Sphere) GetCenter() Vec3 {
	return s.center
}

func (s *Sphere) Intersect(r Ray) (bool, Hit) {
	distance := r.origin.Sub(s.center)
	b := dotProduct(distance, r.dir)
	c := dotProduct(distance, distance) - s.radius*s.radius

	if c > 0 && b > 0 {
		return false, Hit{infinity, zeroVec, zeroVec}
	}

	discr := b*b - c

	if discr < 0 {
		return false, Hit{infinity, zeroVec, zeroVec}
	}

	t := -b - math.Sqrt(discr)

	if t < 0 {
		t = 0
	}

	interSection := r.origin.Add(r.dir.Mul(t))
	n := interSection.Sub(s.center).normalize()
	return true, Hit{t, interSection, n}

	/*distance := r.origin.Sub(s.center)
	radiusSquared := s.radius * s.radius
	b := 2 * dotProduct(r.dir, distance)
	c := dotProduct(distance, distance) - radiusSquared
	discriminant := b*b - 4*c
	// No intersection
	if discriminant < 0.0 {
		return false, Hit{infinity, zeroVec, zeroVec}
	}
	d := math.Sqrt(discriminant)

	t0 := (-1 * b) - d
	t1 := (-1 * b) + d
	t := math.Min(t0, t1)
	t /= 2
	if t < 0 {
		return false, Hit{0, zeroVec, zeroVec}
	}
	intersectionPoint := r.origin.Add(r.dir.Mul(t))
	normal := intersectionPoint.Sub(s.center).normalize()
	solution1 := Hit{t, intersectionPoint, normal}
	return true, solution1
	*/

}
