package godel

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Skeleton struct {
	Name string
	bones []*Bone
}

type Bone struct {
	Name string
	Parent *Bone
	Children []*Bone
	Socket string
	//
	T mgl64.Vec3
	R mgl64.Quat
	S mgl64.Vec3
	W []float64
}