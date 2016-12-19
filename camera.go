package main

import (
	"math"
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

// Camera is used to calculate rays into the scene
type Camera struct {
	eye         Vec3
	width       float64
	height      float64
	fov         float64
	scale       float64
	aspectRatio float64
	depth       int
}

func degToRad(d float64) float64 { return d * math.Pi / 180 }
func radToDeg(r float64) float64 { return r / math.Pi / 180 }

// Init sets up the camera struct
func (c *Camera) Init(eye Vec3, w, h int) {
	c.eye = eye
	c.fov = 90
	c.width = float64(w)
	c.height = float64(h)
	c.depth = 1
	c.scale = math.Tan(degToRad(c.fov * 0.5))
	c.aspectRatio = c.width / c.height
}

func (c *Camera) rayForPixel(x int, y int) Ray {
	dir := Vec3{(2*(float64(x)+0.5)/c.width - 1) * c.aspectRatio * c.scale, (1 - 2*(float64(y)+0.5)/c.height) * c.scale,
		float64(c.depth)}.Normalize()
	return Ray{c.eye, dir}
}
