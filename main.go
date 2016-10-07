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
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"runtime"
	"sync"

	pb "gopkg.in/cheggaaa/pb.v1"
)

var infinity = math.Inf(1)
var zeroVec = Vec3{0, 0, 0}
var delta = math.Sqrt(1.0E-16)

const MaxDepth = 2

var wg sync.WaitGroup

var backgroundColor = Vec3{0.1, 0.1, 0.1}

type Ray struct {
	Origin, Direction Vec3
}

type BoundingVolume interface {
	Intersect(r Ray) bool
}

type Geometry interface {
	IntersectHit(r Ray) Hit
	GetColor() Vec3
}

type Rect struct {
	left   int
	right  int
	top    int
	bottom int
}

func main() {
	//defer profile.Start().Stop()
	// Image size
	imageRes := 512
	w, h := imageRes, imageRes
	// define chunk size for rendering
	chunkSize := 16
	t := image.NewRGBA(image.Rect(0, 0, w, h))
	// Create geometry for the scene
	mesh, err := LoadOBJ("teapot.obj")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	sp1 := &Sphere{center: Vec3{0, 0, 5}, radius: 1.0, color: Vec3{0, 0.7, 0}}
	sp2 := &Sphere{center: Vec3{-2, -1.5, 3}, radius: 1.0, color: Vec3{0.1, 0.9, .7}}
	sp3 := &Sphere{center: Vec3{-2, 1.5, 5}, radius: 1.0, color: Vec3{0.9, 0.9, .1}}
	sp4 := &Sphere{center: Vec3{2, 1.5, 5}, radius: 1.0, color: Vec3{0.9, 0.1, .9}}
	sp5 := &Sphere{center: Vec3{2, -1.5, 5}, radius: 1.0, color: Vec3{0.2, 0.4, .6}}

	// Setup the renderer
	light := Light{Vec3{-1.0, -2.0, 2.0}.Normalize(), 20}
	scene := &Scene{light, []Geometry{mesh, sp1, sp2, sp3, sp4, sp5}}
	eye := Vec3{0, 0, -4.0}
	camera := Camera{eye, w, h, imageRes}
	jobChan := make(chan Rect, 10)
	renderer := Renderer{scene, t, &camera, jobChan}
	///////////////////
	fmt.Println("Rendering...")
	bar := pb.StartNew((imageRes / chunkSize) * (imageRes / chunkSize))
	// Create workers to render chunks
	for i := 0; i < runtime.NumCPU()*2; i++ {
		go renderer.worker(bar)
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
	bar.FinishPrint("")
	outputPath := "img.png"
	fmt.Printf("Writing output to: %s ... ", outputPath)
	defer fmt.Println("Done")
	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer outFile.Close()
	bufWriter := bufio.NewWriter(outFile)
	defer bufWriter.Flush()
	png.Encode(bufWriter, renderer.img)
}
