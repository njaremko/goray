package main

import (
	"math"
)

func ratioToColor(f float64) uint8 {
	scaled := f * 255.0
	switch {
	case scaled < 0:
		scaled = 0
	case scaled > 255:
		scaled = 255
	}
	return uint8(scaled)
}

func colorToRatio(i uint8) float64 {
	return float64(i) / 255.0
}

func sRGBToLinear(f float64) float64 {
	a := 0.055
	if f <= 0.04045 {
		return f / 12.92
	}
	return math.Pow((f+a)/(1+a), 2.4)
}

func linearToSRGB(f float64) float64 {
	a := 0.055
	if f <= 0.0031308 {
		return f * 12.92
	}
	return (1+a)*math.Pow(f, 1/2.4) - a
}

func linearToGamma(f float64) float64 {
	return math.Pow(f, 1/2.2)
}

func gammaToLinear(f float64) float64 {
	return math.Pow(f, 2.2)
}

func (v *Vec3) sRGBToLinear() {
	v.x = sRGBToLinear(v.x)
	v.y = sRGBToLinear(v.y)
	v.z = sRGBToLinear(v.z)
}

func (v *Vec3) linearToSRGB() {
	v.x = linearToSRGB(v.x)
	v.y = linearToSRGB(v.y)
	v.z = linearToSRGB(v.z)
}

func (v *Vec3) linearToGamma() {
	v.x = linearToGamma(v.x)
	v.y = linearToGamma(v.y)
	v.z = linearToGamma(v.z)
}

func (v *Vec3) gammaToLinear() {
	v.x = gammaToLinear(v.x)
	v.y = gammaToLinear(v.y)
	v.z = gammaToLinear(v.z)
}
