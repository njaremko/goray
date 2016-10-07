package main

import "fmt"

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
	kd        *KdTree
}

func NewMesh(triangles []*Triangle) *Mesh {
	fmt.Printf("Building k-d tree... ")
	kdTree := BuildTree(triangles)
	fmt.Println("Done")
	return &Mesh{triangles, kdTree}
}

func (m Mesh) GetColor() Vec3 {
	return Vec3{0.1, 0.7, 0.9}
}

func (m Mesh) IntersectHit(r Ray) Hit {
	hit := m.kd.Intersect(r)
	if hit.IsHit() {
		return hit
	}
	return NoHit
}
