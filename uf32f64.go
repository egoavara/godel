package godel

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

func dVec3ToVec3(a mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(a[0]),
		float32(a[1]),
		float32(a[2]),
	}
}
