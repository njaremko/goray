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
type Vec3 [3]float64

func (v Vec3) Add(vectors ...Vec3) Vec3 {
	for _, vector := range vectors {
		v[0] += vector[0]
		v[1] += vector[1]
		v[2] += vector[2]
	}
	return v
}

func (v Vec3) Sub(vectors ...Vec3) Vec3 {
	for _, vector := range vectors {
		v[0] -= vector[0]
		v[1] -= vector[1]
		v[2] -= vector[2]
	}
	return v
}

func (v Vec3) Mul(x float64) Vec3 {
	v[0] *= x
	v[1] *= x
	v[2] *= x
	return v
}

func (v Vec3) MulVec(v2 Vec3) Vec3 {
	v[0] *= v2[0]
	v[1] *= v2[1]
	v[2] *= v2[2]
	return v
}

func (v Vec3) Reflect(n Vec3) Vec3 {
	return v.Sub(n.Mul(dotProduct(v, n) * 2))
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}

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

func (v Vec3) Inverse() Vec3 {
	v[0] = 1 / v[0]
	v[1] = 1 / v[1]
	v[2] = 1 / v[2]
	return v
}

func (v Vec3) Distance(v2 Vec3) float64 {
	sub := v.Sub(v2)
	dot := dotProduct(sub, sub)
	return math.Sqrt(dot)
}

func (v Vec3) Magnitude() float64 {
	return math.Sqrt(dotProduct(v, v))
}

func (v Vec3) normalize() Vec3 {
	return v.Mul(1 / math.Sqrt(dotProduct(v, v)))
}

func dotProduct(a, b Vec3) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func crossProduct(a, b Vec3) Vec3 {
	return Vec3{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}
