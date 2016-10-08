package main

// NoHit is a const that is used when rays miss
var NoHit = Hit{infinity, zeroVec, zeroVec}

// Hit represents a hit if one occurs
type Hit struct {
	T      float64
	Point  Vec3
	Normal Vec3
}

// HitInfo holds information regarding a hit
type HitInfo struct {
	object Geometry
}

// IsHit returns true if a hit occured
func (h Hit) IsHit() bool {
	return h.T < infinity
}
