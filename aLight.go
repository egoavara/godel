package godel

import (
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
	"math"
)

type Lighting struct {
	Global *DirectionalLight
	// TODO Locals []*PointLight
	// TODO IBLs []*ImageBaseLight
}

func NewLighting() *Lighting {
	return &Lighting{
		Global: &DirectionalLight{
			Direction: mgl32.Vec3{0, -1, 0},
			Color:     mgl32.Vec3{1, 1, 1},
		},
	}
}

type DirectionalLight struct {
	Direction mgl32.Vec3
	Color     mgl32.Vec3
}

func (s *DirectionalLight) SetColor(c color.Color) {
	r, g, b, _ := c.RGBA()
	s.Color[0] = float32(r) / math.MaxUint16
	s.Color[1] = float32(g) / math.MaxUint16
	s.Color[2] = float32(b) / math.MaxUint16
}

type PointLight struct {
	Position    mgl32.Vec3
	Color       mgl32.Vec3
	Attenuation mgl32.Vec3
}

func (s *PointLight) SetColor(c color.Color) {
	r, g, b, _ := c.RGBA()
	s.Color[0] = float32(r) / math.MaxUint16
	s.Color[1] = float32(g) / math.MaxUint16
	s.Color[2] = float32(b) / math.MaxUint16
}

type SpotLight struct {
	Position  mgl32.Vec3
	Direction mgl32.Vec3
	Color     mgl32.Vec3
	Cutoff    float32
}

func (s *SpotLight) SetColor(c color.Color) {
	r, g, b, _ := c.RGBA()
	s.Color[0] = float32(r) / math.MaxUint16
	s.Color[1] = float32(g) / math.MaxUint16
	s.Color[2] = float32(b) / math.MaxUint16
}

const (
	IBL_UP = iota
	IBL_DOWN
	IBL_FORWARD
	IBL_BACKWARD
	IBL_RIGHT
	IBL_LEFT
	IBL_SIZE
)

type ImageBaseLight struct {
	BRDF      *image.RGBA
	Diffuse   [IBL_SIZE]*image.RGBA
	Specular  [IBL_SIZE]*image.RGBA
	uBRDF     uint32
	uDiffuse  uint32
	uSpecular uint32
}
