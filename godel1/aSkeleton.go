package godel1

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Skeleton struct {
	Name string
	Root *Bone
}

type Bone struct {
	Name string
	Parent *Bone
	Children []*Bone
	T mgl64.Vec3
	R mgl64.Quat
	S mgl64.Vec3
	W []float64
}