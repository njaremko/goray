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

func (s *Sphere) Intersect(r Ray) bool {
	distance := r.Origin.Sub(s.center)
	b := dotProduct(distance, r.Direction)
	c := dotProduct(distance, distance) - s.radius*s.radius

	if c > 0 && b > 0 {
		return false
	}

	discr := b*b - c

	if discr < 0 {
		return false
	}

	t := -b - math.Sqrt(discr)

	if t < 0 {
		t = 0
	}

	return true
}

func (s *Sphere) IntersectHit(r Ray) Hit {
	distance := r.Origin.Sub(s.center)
	b := dotProduct(distance, r.Direction)
	c := dotProduct(distance, distance) - s.radius*s.radius

	if c > 0 && b > 0 {
		return NoHit
	}

	discr := b*b - c

	if discr < 0 {
		return NoHit
	}

	t := -b - math.Sqrt(discr)

	if t < 0 {
		t = 0
	}

	intersection := r.Origin.Add(r.Direction.Mul(t))
	n := intersection.Sub(s.center).Normalize()
	return Hit{t, intersection, n}
}
