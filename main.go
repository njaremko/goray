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
	"bufio"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
)

var infinity = math.Inf(1)
var zeroVec = Vec3{0, 0, 0}
var zeroHit = Hit{0, zeroVec, zeroVec}
var delta = math.Sqrt(1.0E-16)

const MaxDepth = 2

var wg sync.WaitGroup

var backgroundColor = Vec3{0.1, 0.1, 0.1}

type Hit struct {
	distance float64
	point    Vec3
	normal   Vec3
}

type Ray struct {
	origin, dir Vec3
}

type BoundingVolume interface {
	Intersect(r Ray) bool
}

type Geometry interface {
	IntersectHit(r Ray) (bool, Hit)
	GetColor() Vec3
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
		if isHit, hit := object.IntersectHit(ray); isHit {
			// Ensure we draw the closest item
			if hit.distance < minDistance {
				pHit = hit
				closestObject = object
				minDistance = pHit.distance
			}
		}
	}
	// If the ray misses
	if closestObject == nil {
		return backgroundColor
	}

	light := s.light.direction.Mul(-1)
	shadowRay := Ray{pHit.point.Add(pHit.normal.Mul(delta)), light}
	for _, object := range s.geometry {
		if isHit, _ := object.IntersectHit(shadowRay); isHit {
			return zeroVec
		}
	}
	// How much light is reflected
	albedo := 0.18
	normalLightProduct := math.Abs(dotProduct(pHit.normal, light))
	diffColor := albedo / math.Pi * s.light.intensity * math.Max(0, normalLightProduct)
	return closestObject.GetColor().Mul(diffColor)

}

type Rect struct {
	left   int
	right  int
	top    int
	bottom int
}

type Camera struct {
	eye    Vec3
	width  int
	height int
	depth  int
}

func (c *Camera) rayForPixel(x int, y int) Ray {
	dir := Vec3{float64(x) - float64(c.width)*0.5, float64(y) - float64(c.height)*0.5,
		float64(c.depth)}.normalize()
	return Ray{c.eye, dir}
}

type Renderer struct {
	scene   *Scene
	img     *image.RGBA
	cam     *Camera
	jobChan chan Rect
}

func (renderer *Renderer) renderRect(r *Rect) {
	for y := r.top; y < r.bottom; y++ {
		for x := r.left; x < r.right; x++ {
			// Compute primary ray direction
			ray := renderer.cam.rayForPixel(x, y)
			g := renderer.scene.rayTrace(ray, 0)
			colour := color.RGBA{float2byte(g.x), float2byte(g.y), float2byte(g.z), 255}
			renderer.img.Set(x, renderer.cam.height-(y+1), colour)
		}
	}
}

func (renderer *Renderer) worker() {
	for r := range renderer.jobChan {
		renderer.renderRect(&r)
		wg.Done()
	}
}

func float2byte(f float64) byte {
	scaled := 0.5 + f*255.0
	switch {
	case scaled < 0:
		scaled = 0
	case scaled > 255:
		scaled = 255
	}
	return byte(scaled)
}

func main() {
	// Image size
	imageRes := 256
	w, h := imageRes, imageRes
	// define chunk size for rendering
	chunkSize := 16
	t := image.NewRGBA(image.Rect(0, 0, w, h))
	// Create geometry for the scene
	mesh := readObjFile("teapot.obj")
	/*sp1 := &Sphere{center: Vec3{0, 0, 0}, radius: 1.0, color: Vec3{0, 0.7, 0}}
	sp2 := &Sphere{center: Vec3{-2, -1.5, 1}, radius: 1.0, color: Vec3{0.1, 0.9, .7}}
	sp3 := &Sphere{center: Vec3{-2, 1.5, 1}, radius: 1.0, color: Vec3{0.9, 0.9, .1}}
	sp4 := &Sphere{center: Vec3{2, 1.5, 1}, radius: 1.0, color: Vec3{0.9, 0.1, .9}}
	sp5 := &Sphere{center: Vec3{2, -1.5, 1}, radius: 1.0, color: Vec3{0.2, 0.4, .6}}
	geometry := []Geometry{sp1, sp2, sp3, sp4, sp5}
	////////////////////////////////////*/
	// Setup the renderer
	light := Light{Vec3{-2.0, -3.0, 2.0}.normalize(), 1500}
	scene := &Scene{light, []Geometry{mesh}}
	eye := Vec3{0, 0, -4.0}
	camera := Camera{eye, w, h, imageRes}
	jobChan := make(chan Rect)
	renderer := Renderer{scene, t, &camera, jobChan}
	///////////////////

	// Create workers to render chunks
	for i := 0; i < runtime.NumCPU()*2; i++ {
		go renderer.worker()
	}
	// Send chunks to workers
	for y := 0; y < h; y += chunkSize {
		for x := 0; x < w; x += chunkSize {
			wg.Add(1)
			renderer.jobChan <- Rect{x, y, x + chunkSize, y + chunkSize}
		}
	}
	// Wait for all jobs to finish
	wg.Wait()
	outFile, err := os.Create("img.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	bufWriter := bufio.NewWriter(outFile)
	defer bufWriter.Flush()
	png.Encode(bufWriter, renderer.img)
}
