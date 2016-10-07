package main

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
