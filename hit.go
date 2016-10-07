package main

var NoHit = Hit{infinity, zeroVec, zeroVec}

type Hit struct {
	T      float64
	Point  Vec3
	Normal Vec3
}

type HitInfo struct {
	object Geometry
}

func (h Hit) IsHit() bool {
	if h.T < infinity {
		return true
	}
	return false
}
