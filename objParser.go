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
	"log"
	"os"
	"strconv"
	"strings"
)

func readObjFile(file string) Mesh {
	inputFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	var vertexSlice []Vec3
	var normalSlice []Vec3
	var triangleSlice []*Triangle
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		parsedLine := strings.Split(scanner.Text(), " ")
		switch parsedLine[0] {
		case "v":
			tx, err := strconv.ParseFloat(parsedLine[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			ty, err := strconv.ParseFloat(parsedLine[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			tz, err := strconv.ParseFloat(parsedLine[3], 64)
			if err != nil {
				log.Fatal(err)
			}
			vertexSlice = append(vertexSlice, Vec3{tx, ty, tz})
		case "vn":
			tx, err := strconv.ParseFloat(parsedLine[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			ty, err := strconv.ParseFloat(parsedLine[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			tz, err := strconv.ParseFloat(parsedLine[3], 64)
			if err != nil {
				log.Fatal(err)
			}
			normalSlice = append(normalSlice, Vec3{tx, ty, tz})
		case "f":
			var v1, n1, v2, n2, v3, n3 int
			joined := strings.Join(parsedLine[1:], " ")
			fmt.Sscanf(joined, "%d//%d %d//%d %d//%d", &v1, &n1, &v2, &n2, &v3, &n3)
			triangleSlice = append(triangleSlice, &Triangle{v1: vertexSlice[v1-1], v2: vertexSlice[v2-1], v3: vertexSlice[v3-1], color: Vec3{0.1, 0.7, 0.9}})
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(os.Stderr, "reading standard input:", err)
	}
	boundingBox, err := computeBoundingBox(vertexSlice)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	return Mesh{triangleSlice, boundingBox}
}
