package main

import (
	"image"
	"image/color"
	"math"

	pb "gopkg.in/cheggaaa/pb.v1"
)

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
			g.linearToSRGB()
			colour := color.RGBA64{ratioToColor(g.X), ratioToColor(g.Y), ratioToColor(g.Z), 65535}
			renderer.img.Set(x, renderer.cam.height-(y+1), colour)
		}
	}
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
	return closestObject.GetColor().Mul(diffColor)
}
