package godel

import (
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/color"
	"math"
)

type Lighting struct {
	Global *GlobalLight
	// TODO Locals []*LocalLight
	// TODO IBLs []*ImageBaseLight
}

func NewLighting() *Lighting {
	return &Lighting{
		Global: &GlobalLight{
			Direction: mgl32.Vec3{0, -1, 0},
			Color:     mgl32.Vec3{1, 1, 1},
		},
	}
}

type GlobalLight struct {
	Direction mgl32.Vec3
	Color     mgl32.Vec3
	// TODO Ambient float32
	// TODO SkySphere *ImageBaseLight
}

func (s *GlobalLight) SetColor(c color.Color) {
	r, g, b, _ := c.RGBA()
	s.Color[0] = float32(r) / math.MaxUint16
	s.Color[1] = float32(g) / math.MaxUint16
	s.Color[2] = float32(b) / math.MaxUint16
}

type LocalLight struct {
	Position mgl32.Vec3
	Color    mgl32.Vec3
	Lumen    float32
}

func (s *LocalLight) SetColor(c color.Color) {
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
	Position mgl32.Vec3
	Light    [IBL_SIZE]image.Image
	uBRDF uint32
	uDiffuse uint32
	uSpecular uint32
}
