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
	"image"
	"image/color"
	"math"

	pb "gopkg.in/cheggaaa/pb.v1"
)

type Renderer struct {
	scene      *Scene
	maxX, maxY int
	pixelChan  chan Pixel
	cam        *Camera
	jobChan    chan rect
}

type Pixel struct {
	x, y  int
	color color.Color
}

func (renderer *Renderer) renderRect(r *rect) {
	for y := r.top; y < r.bottom; y++ {
		for x := r.left; x < r.right; x++ {
			// Compute primary ray direction
			ray := renderer.cam.rayForPixel(x, y)
			g := renderer.scene.rayTrace(ray, 0)
			g.linearToSRGB()
			colour := color.RGBA64{ratioToColor(g.X), ratioToColor(g.Y), ratioToColor(g.Z), 65535}
			renderer.pixelChan <- Pixel{x, y, colour}
		}
	}
}

// CreateImage reads the pixel channel to create the final image
func (renderer *Renderer) CreateImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, renderer.maxX, renderer.maxY))
	for pixel := range renderer.pixelChan {
		img.Set(pixel.x, pixel.y, pixel.color)
	}
	return img
}

func (renderer *Renderer) worker(bar *pb.ProgressBar) {
	for r := range renderer.jobChan {
		renderer.renderRect(&r)
		bar.Increment()
		wg.Done()
	}
}

type Scene struct {
	light    Light
	geometry []Geometry
}

type Light struct {
	direction Vec3
	intensity float64
}

func (s *Scene) rayTrace(ray Ray, depth int) Vec3 {
	var pHit Hit
	var minDistance = infinity
	var closestObject Geometry
	// Search for ray intersection in scene
	for _, object := range s.geometry {
		// Test if the ray hits any scene geometry
		if hit := object.IntersectHit(ray); hit.IsHit() {
			// Ensure we draw the closest item
			if hit.T < minDistance {
				pHit = hit
				closestObject = object
				minDistance = pHit.T
			}
		}
	}
	// If the ray misses
	if closestObject == nil {
		return backgroundColor
	}

	light := s.light.direction.Mul(-1)
	shadowRay := Ray{pHit.Point.Add(pHit.Normal.Mul(EPSILON)), light}
	for _, object := range s.geometry {
		if hit := object.IntersectHit(shadowRay); hit.IsHit() {
			return zeroVec
		}
	}
	// How much light is reflected
	albedo := 0.18
	normalLightProduct := dotProduct(pHit.Normal, light)
	diffColor := albedo / math.Pi * s.light.intensity * math.Max(0, normalLightProduct)
	return closestObject.Color().Mul(diffColor)
}
