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
	"math"
)

// Matrix is a 4x4 matrix struct used for transforming geometry
type Matrix struct {
	x00, x01, x02, x03 float64
	x10, x11, x12, x13 float64
	x20, x21, x22, x23 float64
	x30, x31, x32, x33 float64
}

// Identity returns the identity matrix
func Identity() Matrix {
	return Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

// Translate performs a matrix translation of input v
func Translate(v Vec3) Matrix {
	return Matrix{
		1, 0, 0, v.X,
		0, 1, 0, v.Y,
		0, 0, 1, v.Z,
		0, 0, 0, 1}
}

// Scale performs a matrix scale on input v
func Scale(v Vec3) Matrix {
	return Matrix{
		v.X, 0, 0, 0,
		0, v.Y, 0, 0,
		0, 0, v.Z, 0,
		0, 0, 0, 1}
}

// Rotate rotates a given matrix a degrees
func Rotate(v Vec3, a float64) Matrix {
	v = v.Normalize()
	s := math.Sin(a)
	c := math.Cos(a)
	m := 1 - c
	return Matrix{
		m*v.X*v.X + c, m*v.X*v.Y + v.Z*s, m*v.Z*v.X - v.Y*s, 0,
		m*v.X*v.Y - v.Z*s, m*v.Y*v.Y + c, m*v.Y*v.Z + v.X*s, 0,
		m*v.Z*v.X + v.Y*s, m*v.Y*v.Z - v.X*s, m*v.Z*v.Z + c, 0,
		0, 0, 0, 1}
}

// Translate performs a translation on m with vector v
func (m Matrix) Translate(v Vec3) Matrix {
	return Translate(v).Mul(m)
}

// Scale performs a scale on m with vector v
func (m Matrix) Scale(v Vec3) Matrix {
	return Scale(v).Mul(m)
}

// Rotate performs a rotation on m with vector v and angle a
func (m Matrix) Rotate(v Vec3, a float64) Matrix {
	return Rotate(v, a).Mul(m)
}

// Mul performs matrix multiplication on two matrices
func (m Matrix) Mul(b Matrix) Matrix {
	result := Matrix{}
	result.x00 = m.x00*b.x00 + m.x01*b.x10 + m.x02*b.x20 + m.x03*b.x30
	result.x10 = m.x10*b.x00 + m.x11*b.x10 + m.x12*b.x20 + m.x13*b.x30
	result.x20 = m.x20*b.x00 + m.x21*b.x10 + m.x22*b.x20 + m.x23*b.x30
	result.x30 = m.x30*b.x00 + m.x31*b.x10 + m.x32*b.x20 + m.x33*b.x30
	result.x01 = m.x00*b.x01 + m.x01*b.x11 + m.x02*b.x21 + m.x03*b.x31
	result.x11 = m.x10*b.x01 + m.x11*b.x11 + m.x12*b.x21 + m.x13*b.x31
	result.x21 = m.x20*b.x01 + m.x21*b.x11 + m.x22*b.x21 + m.x23*b.x31
	result.x31 = m.x30*b.x01 + m.x31*b.x11 + m.x32*b.x21 + m.x33*b.x31
	result.x02 = m.x00*b.x02 + m.x01*b.x12 + m.x02*b.x22 + m.x03*b.x32
	result.x12 = m.x10*b.x02 + m.x11*b.x12 + m.x12*b.x22 + m.x13*b.x32
	result.x22 = m.x20*b.x02 + m.x21*b.x12 + m.x22*b.x22 + m.x23*b.x32
	result.x32 = m.x30*b.x02 + m.x31*b.x12 + m.x32*b.x22 + m.x33*b.x32
	result.x03 = m.x00*b.x03 + m.x01*b.x13 + m.x02*b.x23 + m.x03*b.x33
	result.x13 = m.x10*b.x03 + m.x11*b.x13 + m.x12*b.x23 + m.x13*b.x33
	result.x23 = m.x20*b.x03 + m.x21*b.x13 + m.x22*b.x23 + m.x23*b.x33
	result.x33 = m.x30*b.x03 + m.x31*b.x13 + m.x32*b.x23 + m.x33*b.x33
	return result
}

// MulPoint applies a transofrmation to a Vec3 point
func (m Matrix) MulPoint(b Vec3) Vec3 {
	x := m.x00*b.X + m.x01*b.Y + m.x02*b.Z + m.x03
	y := m.x10*b.X + m.x11*b.Y + m.x12*b.Z + m.x13
	z := m.x20*b.X + m.x21*b.Y + m.x22*b.Z + m.x23
	return Vec3{x, y, z}
}

// MulDirection applies a transformation matrix to Vec3 direction
func (m Matrix) MulDirection(b Vec3) Vec3 {
	x := m.x00*b.X + m.x01*b.Y + m.x02*b.Z
	y := m.x10*b.X + m.x11*b.Y + m.x12*b.Z
	z := m.x20*b.X + m.x21*b.Y + m.x22*b.Z
	return Vec3{x, y, z}.Normalize()
}

// MulRay applies a transformation matrix to a Ray
func (m Matrix) MulRay(b Ray) Ray {
	return Ray{m.MulPoint(b.Origin), m.MulDirection(b.Direction)}
}

// MulBox applies a transformation matrix to a Box
func (m Matrix) MulBox(box Box) Box {
	// http://dev.theomader.com/transform-bounding-boxes/
	r := Vec3{m.x00, m.x10, m.x20}
	u := Vec3{m.x01, m.x11, m.x21}
	b := Vec3{m.x02, m.x12, m.x22}
	t := Vec3{m.x03, m.x13, m.x23}
	xa := r.Mul(box.min.X)
	xb := r.Mul(box.max.X)
	ya := u.Mul(box.min.Y)
	yb := u.Mul(box.max.Y)
	za := b.Mul(box.min.Z)
	zb := b.Mul(box.max.Z)
	xa, xb = xa.Min(xb), xa.Max(xb)
	ya, yb = ya.Min(yb), ya.Max(yb)
	za, zb = za.Min(zb), za.Max(zb)
	min := xa.Add(ya).Add(za).Add(t)
	max := xb.Add(yb).Add(zb).Add(t)
	return Box{min, max}
}

// Transpose returns a transpose of input matrix
func (m Matrix) Transpose() Matrix {
	return Matrix{
		m.x00, m.x10, m.x20, m.x30,
		m.x01, m.x11, m.x21, m.x31,
		m.x02, m.x12, m.x22, m.x32,
		m.x03, m.x13, m.x23, m.x33}
}

// Determinant returns the determinant of a given matrix
func (m Matrix) Determinant() float64 {
	return (m.x00*m.x11*m.x22*m.x33 - m.x00*m.x11*m.x23*m.x32 +
		m.x00*m.x12*m.x23*m.x31 - m.x00*m.x12*m.x21*m.x33 +
		m.x00*m.x13*m.x21*m.x32 - m.x00*m.x13*m.x22*m.x31 -
		m.x01*m.x12*m.x23*m.x30 + m.x01*m.x12*m.x20*m.x33 -
		m.x01*m.x13*m.x20*m.x32 + m.x01*m.x13*m.x22*m.x30 -
		m.x01*m.x10*m.x22*m.x33 + m.x01*m.x10*m.x23*m.x32 +
		m.x02*m.x13*m.x20*m.x31 - m.x02*m.x13*m.x21*m.x30 +
		m.x02*m.x10*m.x21*m.x33 - m.x02*m.x10*m.x23*m.x31 +
		m.x02*m.x11*m.x23*m.x30 - m.x02*m.x11*m.x20*m.x33 -
		m.x03*m.x10*m.x21*m.x32 + m.x03*m.x10*m.x22*m.x31 -
		m.x03*m.x11*m.x22*m.x30 + m.x03*m.x11*m.x20*m.x32 -
		m.x03*m.x12*m.x20*m.x31 + m.x03*m.x12*m.x21*m.x30)
}

// Inverse returns the inverse of a given matrix
func (m Matrix) Inverse() Matrix {
	result := Matrix{}
	d := m.Determinant()
	result.x00 = (m.x12*m.x23*m.x31 - m.x13*m.x22*m.x31 + m.x13*m.x21*m.x32 - m.x11*m.x23*m.x32 - m.x12*m.x21*m.x33 + m.x11*m.x22*m.x33) / d
	result.x01 = (m.x03*m.x22*m.x31 - m.x02*m.x23*m.x31 - m.x03*m.x21*m.x32 + m.x01*m.x23*m.x32 + m.x02*m.x21*m.x33 - m.x01*m.x22*m.x33) / d
	result.x02 = (m.x02*m.x13*m.x31 - m.x03*m.x12*m.x31 + m.x03*m.x11*m.x32 - m.x01*m.x13*m.x32 - m.x02*m.x11*m.x33 + m.x01*m.x12*m.x33) / d
	result.x03 = (m.x03*m.x12*m.x21 - m.x02*m.x13*m.x21 - m.x03*m.x11*m.x22 + m.x01*m.x13*m.x22 + m.x02*m.x11*m.x23 - m.x01*m.x12*m.x23) / d
	result.x10 = (m.x13*m.x22*m.x30 - m.x12*m.x23*m.x30 - m.x13*m.x20*m.x32 + m.x10*m.x23*m.x32 + m.x12*m.x20*m.x33 - m.x10*m.x22*m.x33) / d
	result.x11 = (m.x02*m.x23*m.x30 - m.x03*m.x22*m.x30 + m.x03*m.x20*m.x32 - m.x00*m.x23*m.x32 - m.x02*m.x20*m.x33 + m.x00*m.x22*m.x33) / d
	result.x12 = (m.x03*m.x12*m.x30 - m.x02*m.x13*m.x30 - m.x03*m.x10*m.x32 + m.x00*m.x13*m.x32 + m.x02*m.x10*m.x33 - m.x00*m.x12*m.x33) / d
	result.x13 = (m.x02*m.x13*m.x20 - m.x03*m.x12*m.x20 + m.x03*m.x10*m.x22 - m.x00*m.x13*m.x22 - m.x02*m.x10*m.x23 + m.x00*m.x12*m.x23) / d
	result.x20 = (m.x11*m.x23*m.x30 - m.x13*m.x21*m.x30 + m.x13*m.x20*m.x31 - m.x10*m.x23*m.x31 - m.x11*m.x20*m.x33 + m.x10*m.x21*m.x33) / d
	result.x21 = (m.x03*m.x21*m.x30 - m.x01*m.x23*m.x30 - m.x03*m.x20*m.x31 + m.x00*m.x23*m.x31 + m.x01*m.x20*m.x33 - m.x00*m.x21*m.x33) / d
	result.x22 = (m.x01*m.x13*m.x30 - m.x03*m.x11*m.x30 + m.x03*m.x10*m.x31 - m.x00*m.x13*m.x31 - m.x01*m.x10*m.x33 + m.x00*m.x11*m.x33) / d
	result.x23 = (m.x03*m.x11*m.x20 - m.x01*m.x13*m.x20 - m.x03*m.x10*m.x21 + m.x00*m.x13*m.x21 + m.x01*m.x10*m.x23 - m.x00*m.x11*m.x23) / d
	result.x30 = (m.x12*m.x21*m.x30 - m.x11*m.x22*m.x30 - m.x12*m.x20*m.x31 + m.x10*m.x22*m.x31 + m.x11*m.x20*m.x32 - m.x10*m.x21*m.x32) / d
	result.x31 = (m.x01*m.x22*m.x30 - m.x02*m.x21*m.x30 + m.x02*m.x20*m.x31 - m.x00*m.x22*m.x31 - m.x01*m.x20*m.x32 + m.x00*m.x21*m.x32) / d
	result.x32 = (m.x02*m.x11*m.x30 - m.x01*m.x12*m.x30 - m.x02*m.x10*m.x31 + m.x00*m.x12*m.x31 + m.x01*m.x10*m.x32 - m.x00*m.x11*m.x32) / d
	result.x33 = (m.x01*m.x12*m.x20 - m.x02*m.x11*m.x20 + m.x02*m.x10*m.x21 - m.x00*m.x12*m.x21 - m.x01*m.x10*m.x22 + m.x00*m.x11*m.x22) / d
	return result
}

// Frustum calculates a frustum culling matrix given...
func Frustum(left, right, bottom, top, near, far float64) Matrix {
	t1 := 2 * near
	t2 := right - left
	t3 := top - bottom
	t4 := far - near
	return Matrix{
		t1 / t2, 0, (right + left) / t2, 0,
		0, t1 / t3, (top + bottom) / t3, 0,
		0, 0, (-far - near) / t4, (-t1 * far) / t4,
		0, 0, -1, 0}
}

// Orthographic calculates an orthographic projection matrix given...
func Orthographic(left, right, bottom, top, near, far float64) Matrix {
	return Matrix{
		2 / (right - left), 0, 0, -(right + left) / (right - left),
		0, 2 / (top - bottom), 0, -(top + bottom) / (top - bottom),
		0, 0, -2 / (far - near), -(far + near) / (far - near),
		0, 0, 0, 1}
}

// Perspective calculates a perspective matrix given...
func Perspective(fovy, aspect, near, far float64) Matrix {
	ymax := near * math.Tan(fovy*math.Pi/360)
	xmax := ymax * aspect
	return Frustum(-xmax, xmax, -ymax, ymax, near, far)
}

// LookAtMatrix calculates what the matrix looks like to an eye
func LookAtMatrix(eye, center, up Vec3) Matrix {
	up = up.Normalize()
	f := center.Sub(eye).Normalize()
	s := crossProduct(f, up).Normalize()
	u := crossProduct(s, f)
	m := Matrix{
		s.X, u.X, f.X, 0,
		s.Y, u.Y, f.Y, 0,
		s.Z, u.Z, f.Z, 0,
		0, 0, 0, 1,
	}
	return m.Transpose().Inverse().Translate(eye)
}

// Frustum returns a frustum culled matrix m given...
func (m Matrix) Frustum(left, right, bottom, top, near, far float64) Matrix {
	return Frustum(left, right, bottom, top, near, far).Mul(m)
}

// Orthographic returns an orthographic projection matrix given...
func (m Matrix) Orthographic(left, right, bottom, top, near, far float64) Matrix {
	return Orthographic(left, right, bottom, top, near, far).Mul(m)
}

// Perspective returns a perspective matrix given...
func (m Matrix) Perspective(fovy, aspect, near, far float64) Matrix {
	return Perspective(fovy, aspect, near, far).Mul(m)
}
