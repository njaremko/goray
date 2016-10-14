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
	"os"
	"strconv"
	"strings"
)

func parseIndex(value string, length int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		//		fmt.Fprintln(os.Stderr, err)
	}
	if parsed < 0 {
		parsed += length
	}
	return parsed
}

// OpenOBJ takes the path to an obj file and returns a Mesh pointer
func OpenOBJ(path string) (*Mesh, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	// Allocate slices to store data
	vectors := make([]Vec3, 1, 1024)
	textures := make([]Vec3, 1, 1024)
	normals := make([]Vec3, 1, 1024)
	var triangles []*Triangle
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		keyword := fields[0]
		args := fields[1:]
		switch keyword {
		case "v":
			f := parseFloats(args)
			v := Vec3{f[0], f[1], f[2]}
			vectors = append(vectors, v)
		case "vt":
			f := parseFloats(args)
			v := Vec3{f[0], f[1], 0}
			textures = append(textures, v)
		case "vn":
			f := parseFloats(args)
			v := Vec3{f[0], f[1], f[2]}
			normals = append(normals, v)
		case "f":
			fVectors := make([]int, len(args))
			fTextures := make([]int, len(args))
			fNormals := make([]int, len(args))
			for i, arg := range args {
				vertex := strings.Split(arg+"//", "/")
				fVectors[i] = parseIndex(vertex[0], len(vectors))
				fTextures[i] = parseIndex(vertex[1], len(textures))
				fNormals[i] = parseIndex(vertex[2], len(normals))
			}
			for i := 1; i < len(fVectors)-1; i++ {
				i1, i2, i3 := 0, i, i+1
				t := Triangle{}
				t.V1 = vectors[fVectors[i1]]
				t.V2 = vectors[fVectors[i2]]
				t.V3 = vectors[fVectors[i3]]
				t.T1 = textures[fTextures[i1]]
				t.T2 = textures[fTextures[i2]]
				t.T3 = textures[fTextures[i3]]
				t.N1 = normals[fNormals[i1]]
				t.N2 = normals[fNormals[i2]]
				t.N3 = normals[fNormals[i3]]
				t.fixNormals()
				triangles = append(triangles, &t)
			}
		}
	}

	return newMesh(triangles), scanner.Err()
}
