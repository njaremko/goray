package main

import (
	"math"
)

func ratioToColor(f float64) uint16 {
	scaled := f * 65535
	switch {
	case scaled < 0:
		scaled = 0
	case scaled > 65535:
		scaled = 65535
	}
	return uint16(scaled)
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
	v.X = sRGBToLinear(v.X)
	v.Y = sRGBToLinear(v.Y)
	v.Z = sRGBToLinear(v.Z)
}

func (v *Vec3) linearToSRGB() {
	v.X = linearToSRGB(v.X)
	v.Y = linearToSRGB(v.Y)
	v.Z = linearToSRGB(v.Z)
}

func (v *Vec3) linearToGamma() {
	v.X = linearToGamma(v.X)
	v.Y = linearToGamma(v.Y)
	v.Z = linearToGamma(v.Z)
}

func (v *Vec3) gammaToLinear() {
	v.X = gammaToLinear(v.X)
	v.Y = gammaToLinear(v.Y)
	v.Z = gammaToLinear(v.Z)
}
