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

// Camera is used to calculate rays into the scene
type Camera struct {
	eye    Vec3
	width  int
	height int
	depth  int
}

func (c *Camera) rayForPixel(x int, y int) Ray {
	dir := Vec3{float64(x) - float64(c.width)*0.5, float64(y) - float64(c.height)*0.5,
		float64(c.depth)}.Normalize()
	return Ray{c.eye, dir}
}
