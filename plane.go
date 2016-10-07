package main

type Plane struct {
	Point  Vec3
	Normal Vec3
	Color  Vec3
}

func (p *Plane) GetColor() Vec3 {
	return p.Color
}

func (p *Plane) IntersectHit(r Ray) Hit {
	denom := dotProduct(p.Normal, r.Direction)
	if denom > EPSILON {
		p0l0 := p.Point.Sub(r.Origin)
		t := dotProduct(p0l0, p.Normal) / denom
		return Hit{t, r.Origin.Add(r.Direction.Mul(t)), p.Normal}
	}

	return NoHit
}
