package main

type Mesh struct {
	triangles []*Triangle
	bv        BoundingVolume
}

func (m Mesh) GetColor() Vec3 {
	return Vec3{0.1, 0.7, 0.9}
}

func (m Mesh) IntersectHit(r Ray) (bool, Hit) {
	var pHit Hit
	var minDistance = infinity
	var closestTriangle Geometry
	if isHit := m.bv.Intersect(r); !isHit {
		return false, zeroHit
	}
	for _, triangle := range m.triangles {
		if isHit, hit := triangle.IntersectHit(r); isHit {
			if hit.distance < minDistance {
				pHit = hit
				closestTriangle = triangle
				minDistance = pHit.distance
			}
		}
	}
	if closestTriangle == nil {
		return false, zeroHit
	}
	return true, pHit
}
