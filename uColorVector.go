package godel

import (
	"github.com/go-gl/mathgl/mgl64"
	"image/color"
	"math"
)

func ColorToVec(c color.Color) mgl64.Vec4 {
	r, g, b, a := c.RGBA()
	var res mgl64.Vec4
	res[0] = mgl64.Clamp(float64(r) / math.MaxUint16, 0, 1)
	res[1] = mgl64.Clamp(float64(g) / math.MaxUint16, 0, 1)
	res[2] = mgl64.Clamp(float64(b) / math.MaxUint16, 0, 1)
	res[3] = mgl64.Clamp(float64(a) / math.MaxUint16, 0, 1)
	return res
}
func VecToColor(v mgl64.Vec4) color.Color {
	return color.RGBA64{
		R: uint16(mgl64.Clamp(v[0] * math.MaxUint16, 0, math.MaxUint16)),
		G: uint16(mgl64.Clamp(v[1] * math.MaxUint16, 0, math.MaxUint16)),
		B: uint16(mgl64.Clamp(v[2] * math.MaxUint16, 0, math.MaxUint16)),
		A: uint16(mgl64.Clamp(v[3] * math.MaxUint16, 0, math.MaxUint16)),
	}
}
