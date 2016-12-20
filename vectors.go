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

// Vec3 is a representation of a 3 dimentional vector
type Vec3 struct {
	X, Y, Z float64
}

// Add implements vector addition
func (v Vec3) Add(vectors ...Vec3) Vec3 {
	for _, vector := range vectors {
		v.X += vector.X
		v.Y += vector.Y
		v.Z += vector.Z
	}
	return v
}

// Sub implements vector subtraction
func (v Vec3) Sub(vectors ...Vec3) Vec3 {
	for _, vector := range vectors {
		v.X -= vector.X
		v.Y -= vector.Y
		v.Z -= vector.Z
	}
	return v
}

// Mul implements scalar multiplication
func (v Vec3) Mul(x float64) Vec3 {
	v.X *= x
	v.Y *= x
	v.Z *= x
	return v
}

// MulVec implements vector multiplication
func (v Vec3) MulVec(v2 Vec3) Vec3 {
	v.X *= v2.X
	v.Y *= v2.Y
	v.Z *= v2.Z
	return v
}

// Reflect implements vector reflection
func (v Vec3) Reflect(n Vec3) Vec3 {
	return v.Sub(n.Mul(dotProduct(v, n) * 2))
}

// Refract implements vector refraction
func (v Vec3) Refract(n Vec3, ior float64) Vec3 {
	cosi := clamp(dotProduct(v, n), -1, 1)
	etai := 1.0
	if cosi < 0 {
		cosi = -cosi
	} else {
		etai, ior = ior, etai
		n.Mul(-1)
	}
	eta := etai / ior
	k := 1 - eta*eta*(1-cosi*cosi)
	if k < 0 {
		return zeroVec
	}
	return v.Mul(eta).Add(n.Mul(eta*cosi - math.Sqrt(k)))
}

// Inverse returns the inverse of a vector
func (v Vec3) Inverse() Vec3 {
	v.X = 1 / v.X
	v.Y = 1 / v.Y
	v.Z = 1 / v.Z
	return v
}

// Distance returns the distance between two vectors
func (v Vec3) Distance(v2 Vec3) float64 {
	sub := v.Sub(v2)
	dot := dotProduct(sub, sub)
	return math.Sqrt(dot)
}

// Equals checks for vector equality
func (v Vec3) Equals(v2 Vec3) bool {
	return v.X == v2.X && v.Y == v2.Y && v.Z == v2.Z
}

// Magnitude returns the magnitude of the vector
func (v Vec3) Magnitude() float64 {
	return math.Sqrt(dotProduct(v, v))
}

// Normalize returns a normalized vector
func (v Vec3) Normalize() Vec3 {
	return v.Mul(1 / math.Sqrt(dotProduct(v, v)))
}

// Min returns the min x, y, and z from two vectors
func (v Vec3) Min(b Vec3) Vec3 {
	return Vec3{math.Min(v.X, b.X), math.Min(v.Y, b.Y), math.Min(v.Z, b.Z)}
}

// Max returns the max x, y, and z from two vectors
func (v Vec3) Max(b Vec3) Vec3 {
	return Vec3{math.Max(v.X, b.X), math.Max(v.Y, b.Y), math.Max(v.Z, b.Z)}
}

func dotProduct(a, b Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func crossProduct(a, b Vec3) Vec3 {
	return Vec3{a.Y*b.Z - a.Z*b.Y, a.Z*b.X - a.X*b.Z, a.X*b.Y - a.Y*b.X}
}
