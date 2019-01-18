package godel1

import (
	"github.com/go-gl/mathgl/mgl64"
)

// Directional
// Point
// Spot
// Probe
// GlobalProbe
type Light interface {
}

type DirectionalLight struct {
	dir mgl64.Vec3
	clr mgl64.Vec3
}
func NewDirectionalLight(dir mgl64.Vec3, clr mgl64.Vec3) *DirectionalLight {
	return &DirectionalLight{dir: dir, clr: clr}
}

type PointLight struct {
	pos               mgl64.Vec3
	clr               mgl64.Vec3
	attenuationRadius float32
}
func NewPointLight(pos mgl64.Vec3, clr mgl64.Vec3, attenuationRadius float32) *PointLight {
	return &PointLight{pos: pos, clr: clr, attenuationRadius: attenuationRadius}
}

type SpotLight struct {
	pos               mgl64.Vec3
	dir               mgl64.Vec3
	clr               mgl64.Vec3
	cutoff            float32
	attenuationRadius float32
}

func NewSpotLight(pos mgl64.Vec3, dir mgl64.Vec3, clr mgl64.Vec3, cutoff float32, attenuationRadius float32) *SpotLight {
	return &SpotLight{pos: pos, dir: dir, clr: clr, cutoff: cutoff, attenuationRadius: attenuationRadius}
}
//type ProbeLight struct {
//	pos      mgl64.Vec3
//	brdf     f32i.Image32
//	diffuse  [6]f32i.Image32
//	specular [6]f32i.Image32
//}
//
//func NewProbeLight(pos mgl64.Vec3, brdf f32i.Image32, diffuse [6]f32i.Image32, specular [6]f32i.Image32) *ProbeLight {
//	return &ProbeLight{pos: pos, brdf: brdf, diffuse: diffuse, specular: specular}
//}
//
//type GlobalProbeLight struct {
//	brdf     f32i.Image32
//	diffuse  [6]f32i.Image32
//	specular [6]f32i.Image32
//}
