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

// Plane is used to store information regarding an infinity plane
type Plane struct {
	Point  Vec3
	Normal Vec3
	color  Vec3
}

// Color returns the color of the plane, used to fufill Geometry interface
func (p *Plane) Color() Vec3 {
	return p.color
}

// IntersectHit performs an intersection test on the plan and returns a Hit
func (p *Plane) IntersectHit(r Ray) Hit {
	denom := dotProduct(p.Normal, r.Direction)
	if denom > EPSILON {
		p0l0 := p.Point.Sub(r.Origin)
		t := dotProduct(p0l0, p.Normal) / denom
		return Hit{t, r.Origin.Add(r.Direction.Mul(t)), p.Normal}
	}

	return NoHit
}
